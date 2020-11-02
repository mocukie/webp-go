package webp

/*
#cgo LDFLAGS: -lwebp
#include "webp.h"

typedef enum OneSetpShopCall {
    OneSetpShopCallLossless = 1 << 0,
    OneSetpShopCallRGB  = 1 << 1,
    OneSetpShopCallBGR  = 1 << 2,
    OneSetpShopCallRGBA = 1 << 3,
    OneSetpShopCallBGRA = 1 << 4,
} OneSetpShopCall;

static size_t OneStepShopEncode(OneSetpShopCall call, const uint8_t* pix, int w, int h, int s, float q, uint8_t** out) {
    int lossless = call & OneSetpShopCallLossless;
    switch (call & ~OneSetpShopCallLossless) {
        case OneSetpShopCallRGB:
            return lossless ? WebPEncodeLosslessRGB(pix, w, h, s, out) : WebPEncodeRGB(pix, w, h, s, q, out);
        case OneSetpShopCallBGR:
            return lossless ? WebPEncodeLosslessBGR(pix, w, h, s, out) : WebPEncodeBGR(pix, w, h, s, q, out);
        case OneSetpShopCallRGBA:
            return lossless? WebPEncodeLosslessRGBA(pix, w, h, s, out) : WebPEncodeRGBA(pix, w, h, s, q, out);
        case OneSetpShopCallBGRA:
            return lossless ? WebPEncodeLosslessBGRA(pix, w, h, s, out) : WebPEncodeBGRA(pix, w, h, s, q, out);
    }
    return 0;
}
*/
import "C"
import (
	"errors"
	"image"
	"image/color"
	"io"
	"unsafe"
)

func EncoderVersion() Version {
	return Version(C.WebPGetEncoderVersion())
}

func Encode(w io.Writer, img image.Image, opts *EncodeOptions) error {
	data, err := encode(img, opts)
	if err != nil {
		return err
	}
	defer data.release()
	if _, err = w.Write(data); err != nil {
		return err
	}
	return nil
}

func EncodeSlice(img image.Image, opts *EncodeOptions) ([]byte, error) {
	data, err := encode(img, opts)
	if err != nil {
		return nil, err
	}
	return data.asSafe(), nil
}

func encode(img image.Image, opts *EncodeOptions) (unsafeBytes, error) {

	var config C.WebPConfig
	var pic C.WebPPicture

	opts.assign(&config)
	if !validateEncodeConfig(&config) {
		return nil, VP8EncErrorInvalidConfiguration
	}

	defer C.WebPPictureFree(&pic)

	holder, err := webpPictureImport(&pic, img, opts)
	if err != nil {
		return nil, err
	}

	var out *C.uint8_t
	var outSize C.size_t
	if holder == nil {
		outSize = C.GoWebPEncode(&pic, &config, &out)
	} else {
		outSize = C.GoWebPEncodeUseGoMem(&pic, &config, &out, *holder)
	}

	if int(outSize) == 0 {
		return nil, VP8EncodeError(pic.error_code)
	}

	return wrapUnsafeBytes(out, outSize), nil
}

func webpPictureImport(pic *C.WebPPicture, img image.Image, opts *EncodeOptions) (*C.PixMemHolder, error) {
	var holder *C.PixMemHolder
	var err error
	var width, height = img.Bounds().Dx(), img.Bounds().Dy()

	if width <= 0 || height <= 0 {
		return nil, VP8EncErrorBadDimension
	}

	pic.width = C.int(width)
	pic.height = C.int(height)
	pic.argb_stride = C.int(width)

	pic.use_argb = bool2CInt(opts.UseSharpYUV || opts.Lossless)
	switch m := img.(type) {
	case *RGBImg:
		err = webpPictureImportCGO(pic, m.Pix[:height*m.Stride], m.Stride, C.WebPPictureImporter(C.WebPPictureImportRGB))
	case *image.RGBA:
		if !opts.Lossless && !opts.UseSharpYUV && m.Opaque() { //import to yuv directly
			err = webpPictureImportCGO(pic, m.Pix[:height*m.Stride], m.Stride, C.WebPPictureImporter(C.WebPPictureImportRGBX))
		} else {
			holder, err = webpPictureImportRGBA(pic, m.Pix, width, height)
		}
	case *image.NRGBA:
		err = webpPictureImportCGO(pic, m.Pix[:height*m.Stride], m.Stride, C.WebPPictureImporter(C.WebPPictureImportRGBA))
	case *image.RGBA64:
		holder, err = webpPictureImportRGBA64(pic, m.Pix, width, height)
	case *image.NRGBA64:
		holder, err = webpPictureImportPix(pic, m.Pix, width, height, 8, 0, 2, 4, 6)
	case *image.Gray:
		holder, err = webpPictureImportPix(pic, m.Pix, width, height, 1, 0, 0, 0, -1)
	case *image.Gray16:
		holder, err = webpPictureImportPix(pic, m.Pix, width, height, 2, 0, 0, 0, -1)
	case *YCbCr:
		holder, err = webpPictureImportYUVA(pic, &m.YCbCr, nil, -1)
	case *NYCbCrA:
		holder, err = webpPictureImportYUVA(pic, &m.YCbCr, m.A, m.AStride)
	case *image.YCbCr, *image.NYCbCrA:
		//image.YCbCr is JFIF standard, but webp yuv is ITU-R BT.601 standard, so we need convert to rgb
		holder, err = webpPictureImportCommon(pic, img)
	case *ARGBImg:
		pic.use_argb = C.int(1)
		holder = &C.PixMemHolder{argb: unsafe.Pointer(&m.Pix[0])}
	default:
		holder, err = webpPictureImportCommon(pic, img)
	}

	return holder, err
}

func webpPictureAllocARGB(pic *C.WebPPicture) ([]uint32, error) {
	pic.use_argb = C.int(1)
	if ok := C.WebPPictureAlloc(pic); int(ok) == 0 {
		return nil, VP8EncErrorOutOfMemory
	}
	size := C.size_t(pic.width * pic.height)

	return (*(*[1 << 28]uint32)(unsafe.Pointer(pic.argb)))[:size:size], nil
}

func webpPictureImportCGO(pic *C.WebPPicture, pix []uint8, stride int, importer C.WebPPictureImporter) error {
	ret := C.GoWebPDoImport(importer, pic, (*C.uint8_t)(&pix[0]), C.int(stride))
	if int(ret) == 0 {
		if err := VP8EncodeError(pic.error_code); err != VP8EncOk {
			return err
		}
		return VP8EncodeError(-1)
	}
	return nil
}

func webpPictureImportPix(pic *C.WebPPicture, pix []uint8, width, height, step, ir, ig, ib, ia int) (*C.PixMemHolder, error) {
	pic.use_argb = C.int(1)
	argb := make([]uint32, width*height)
	holder := &C.PixMemHolder{argb: unsafe.Pointer(&argb[0])}
	if ia > -1 {
		for i := range argb {
			argb[i] = uint32(pix[ia])<<24 | uint32(pix[ir])<<16 | uint32(pix[ig])<<8 | uint32(pix[ib])
			ir += step
			ig += step
			ib += step
			ia += step
		}
	} else {
		for i := range argb {
			argb[i] = 0xff000000 | uint32(pix[ir])<<16 | uint32(pix[ig])<<8 | uint32(pix[ib])
			ir += step
			ig += step
			ib += step
		}
	}
	return holder, nil
}

func webpPictureImportRGBA(pic *C.WebPPicture, pix []uint8, width, height int) (*C.PixMemHolder, error) {
	pic.use_argb = C.int(1)
	argb := make([]uint32, width*height)
	holder := &C.PixMemHolder{argb: unsafe.Pointer(&argb[0])}
	for i := range argb {
		off := i * 4
		if a := uint32(pix[off+3]); a == 0xff {
			argb[i] = 0xff000000 | uint32(pix[off])<<16 | uint32(pix[off+1])<<8 | uint32(pix[off+2])
		} else if a != 0 {
			argb[i] = a<<24 | (uint32(pix[off])*0xff)/a<<16 | (uint32(pix[off+1])*0xff)/a<<8 | (uint32(pix[off+2])*0xff)/a
		}
	}
	return holder, nil
}

func webpPictureImportRGBA64(pic *C.WebPPicture, pix []uint8, width, height int) (*C.PixMemHolder, error) {
	pic.use_argb = C.int(1)
	argb := make([]uint32, width*height)
	holder := &C.PixMemHolder{argb: unsafe.Pointer(&argb[0])}
	for i := range argb {
		off := i * 8
		if a := uint32(pix[off+6])<<8 | uint32(pix[off+7]); a == 0xffff {
			argb[i] = 0xff000000 | uint32(pix[off])<<16 | uint32(pix[off+2])<<8 | uint32(pix[off+4])
		} else if a != 0 {
			argb[i] = a>>8<<24 |
				(((uint32(pix[off])<<8|uint32(pix[off+1]))*0xffff)/a)>>8<<16 |
				(((uint32(pix[off+2])<<8|uint32(pix[off+3]))*0xffff)/a)>>8<<8 |
				(((uint32(pix[off+4])<<8|uint32(pix[off+5]))*0xffff)/a)>>8
		}
	}
	return holder, nil
}

func webpPictureImportYUVA(pic *C.WebPPicture, m *image.YCbCr, A []uint8, AStride int) (*C.PixMemHolder, error) {
	holder := &C.PixMemHolder{}
	pic.use_argb = C.int(0)
	pic.colorspace = C.WEBP_YUV420
	holder.y = (*C.uint8_t)(unsafe.Pointer(&m.Y[0]))
	holder.u = (*C.uint8_t)(unsafe.Pointer(&m.Cb[0]))
	holder.v = (*C.uint8_t)(unsafe.Pointer(&m.Cr[0]))
	pic.y_stride = C.int(m.YStride)
	pic.uv_stride = C.int(m.CStride)
	if A != nil {
		pic.colorspace = C.WEBP_YUV420A
		holder.a = (*C.uint8_t)(unsafe.Pointer(&A[0]))
		pic.a_stride = C.int(AStride)
	}
	return holder, nil
}

func webpPictureImportCommon(pic *C.WebPPicture, img image.Image) (*C.PixMemHolder, error) {
	model := color.NRGBAModel
	rect := img.Bounds()
	w, h, mx, my := rect.Dx(), rect.Dy(), rect.Min.X, rect.Min.Y
	pic.use_argb = C.int(1)
	argb := make([]uint32, w*h)
	holder := &C.PixMemHolder{argb: unsafe.Pointer(&argb[0])}
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := model.Convert(img.At(mx+x, my+y)).(color.NRGBA)
			argb[y*w+x] = uint32(c.A)<<24 | uint32(c.R)<<16 | uint32(c.G)<<8 | uint32(c.B)
		}
	}
	return holder, nil
}

//-----------------
// One-stop-shop call! No questions asked:

// These functions compress using the lossy format, and the quality_factor
// can go from 0 (smaller output, lower quality) to 100 (best quality,
// larger output).
func EncodeRGB(rgb []uint8, width, height, stride int, quality float32) ([]byte, error) {
	return oneStepShopEncode(C.OneSetpShopCallRGB, rgb, width, height, stride, quality)
}

func EncodeBGR(bgr []uint8, width, height, stride int, quality float32) ([]byte, error) {
	return oneStepShopEncode(C.OneSetpShopCallBGR, bgr, width, height, stride, quality)
}

func EncodeRGBA(rgba []uint8, width, height, stride int, quality float32) ([]byte, error) {
	return oneStepShopEncode(C.OneSetpShopCallRGBA, rgba, width, height, stride, quality)
}

func EncodeBGRA(bgra []uint8, width, height, stride int, quality float32) ([]byte, error) {
	return oneStepShopEncode(C.OneSetpShopCallBGRA, bgra, width, height, stride, quality)
}

// These functions are the equivalent of the above, but compressing in a
// lossless manner. Files are usually larger than lossy format, but will
// not suffer any compression loss.
// Note these functions, like the lossy versions, use the library's default
// settings. For lossless this means 'exact' is disabled. RGB values in
// transparent areas will be modified to improve compression. To avoid this,
// use Encode(w io.Writer, img image.Image, opts *EncodeOptions) and set EncodeOptions exact to 1.
func EncodeLosslessRGB(rgb []uint8, width, height, stride int) ([]byte, error) {
	return oneStepShopEncode(C.OneSetpShopCallRGB|C.OneSetpShopCallLossless, rgb, width, height, stride, -1)
}

func EncodeLosslessBGR(bgr []uint8, width, height, stride int) ([]byte, error) {
	return oneStepShopEncode(C.OneSetpShopCallBGR|C.OneSetpShopCallLossless, bgr, width, height, stride, -1)
}

func EncodeLosslessRGBA(rgba []uint8, width, height, stride int) ([]byte, error) {
	return oneStepShopEncode(C.OneSetpShopCallRGBA|C.OneSetpShopCallLossless, rgba, width, height, stride, -1)
}

func EncodeLosslessBGRA(bgra []uint8, width, height, stride int) ([]byte, error) {
	return oneStepShopEncode(C.OneSetpShopCallBGRA|C.OneSetpShopCallLossless, bgra, width, height, stride, -1)
}

func oneStepShopEncode(call C.OneSetpShopCall, pix []uint8, width, height, stride int, quality float32) ([]byte, error) {
	var out *C.uint8_t
	size := C.OneStepShopEncode(call, (*C.uint8_t)(unsafe.Pointer(&pix[0])), C.int(width), C.int(height), C.int(stride), C.float(quality), &out)
	defer C.WebPFree(unsafe.Pointer(out))
	if size == 0 {
		return nil, errors.New("WebP one step shop encode failed")
	}
	output := C.GoBytes(unsafe.Pointer(out), C.int(size))
	return output, nil
}
