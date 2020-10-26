package webp

/*
#cgo LDFLAGS: -lwebp
#include <webp/decode.h>
*/
import "C"
import (
    "github.com/mocukie/webp-go/webp/colorx"
    "image"
    "image/color"
    "io"
    "io/ioutil"
    "unsafe"
)

const webpHeaderSize = 30 //riff(12) + VP8X(8)? + (VP8(10) or VP8L(5))

func DecoderVersion() Version {
    return Version(C.WebPGetDecoderVersion())
}

func Decode(r io.Reader) (image.Image, error) {
    return DecodeEX(r, NewDecOptions())
}

func DecodeConfig(r io.Reader) (image.Config, error) {
    data := make([]byte, webpHeaderSize)
    _, err := io.ReadFull(r, data)
    if err != nil {
        return image.Config{}, err
    }

    f, err := GetBitstreamFeatures(data)
    if err != nil {
        return image.Config{}, err
    }

    var model color.Model
    if f.Format == FormatLossy {
        if f.HasAlpha {
            model = colorx.NYCbCrBT601Model
        } else {
            model = colorx.YCbCrBT601Model
        }
    } else {
        if f.HasAlpha {
            model = color.NRGBAModel
        } else {
            model = colorx.RGBModel
        }
    }

    return image.Config{
        Width:      f.Width,
        Height:     f.Height,
        ColorModel: model,
    }, nil
}

func DecodeEX(r io.Reader, opts *DecodeOptions) (image.Image, error) {
    data, err := ioutil.ReadAll(r)
    if err != nil {
        return nil, err
    }
    return decode(data, opts)
}

func DecodeSlice(data []byte, opts *DecodeOptions) (image.Image, error) {
    return decode(data, opts)
}

func GetBitstreamFeatures(data []byte) (*BitstreamFeatures, error) {
    cData, cSize := bytesGetCPtr(data)
    var f C.WebPBitstreamFeatures
    if code := VP8StatusCode(C.WebPGetFeatures(cData, cSize, &f)); code != VP8StatusOk {
        return nil, code.error("could not get bits stream features, ")
    }
    return &BitstreamFeatures{
        Width:        int(f.width),
        Height:       int(f.height),
        HasAlpha:     int(f.has_alpha) == 1,
        HasAnimation: int(f.has_animation) == 1,
        Format:       BitStreamFormat(f.format),
    }, nil
}

func decode(input []byte, opts *DecodeOptions) (image.Image, error) {
    cData, cSize := bytesGetCPtr(input)
    config := &C.WebPDecoderConfig{}
    if code := VP8StatusCode(C.WebPGetFeatures(cData, cSize, &config.input)); code != VP8StatusOk {
        return nil, code.error("could not get bits stream features, ")
    }
    opts.assign(&config.options)

    var img image.Image
    var width, height = calcOutputSize(config)
    img = opts.ImageType(config, width, height)
    if code := VP8StatusCode(C.WebPDecode(cData, cSize, config)); code != VP8StatusOk {
        return nil, code.error("WebPDecode")
    }

    return img, nil
}

func decPixAuto(config *C.WebPDecoderConfig, width, height int) image.Image {
    var img image.Image
    hasAlpha := int(config.input.has_alpha) == 1
    lossy := BitStreamFormat(config.input.format) == FormatLossy
    if lossy {
        if hasAlpha {
            img = decPixYUVA(config, width, height)
        } else {
            img = decPixYUV(config, width, height)
        }
    } else {
        if hasAlpha {
            img = decPixNRGBA(config, width, height)
        } else {
            img = decPixRGB(config, width, height)
        }
    }
    return img
}

func decPixStd(config *C.WebPDecoderConfig, width, height int) image.Image {
    hasAlpha := int(config.input.has_alpha) == 1
    if hasAlpha {
        return decPixNRGBA(config, width, height)
    } else {
        return decPixRGBA(config, width, height)
    }
}

func decPixRGB(config *C.WebPDecoderConfig, width, height int) image.Image {
    img := NewRGB(image.Rect(0, 0, width, height))
    setupRGBBuf(config, img.Pix, img.Stride, ModeRGB)
    return img
}

func decPixRGBA(config *C.WebPDecoderConfig, width, height int) image.Image {
    img := image.NewRGBA(image.Rect(0, 0, width, height))
    setupRGBBuf(config, img.Pix, img.Stride, ModeRGBA)
    return img
}

func decPixNRGBA(config *C.WebPDecoderConfig, width, height int) image.Image {
    img := image.NewNRGBA(image.Rect(0, 0, width, height))
    setupRGBBuf(config, img.Pix, img.Stride, ModeNRGBA)
    return img
}

func decPixYUV(config *C.WebPDecoderConfig, width, height int) image.Image {
    img := NewYCbCr(image.Rect(0, 0, width, height))
    setupYUVABuf(config, img.Y, img.Cb, img.Cr, nil, img.YStride, img.CStride, -1, ModeYUV)
    return img
}

func decPixYUVA(config *C.WebPDecoderConfig, width, height int) image.Image {
    img := NewNYCbCrA(image.Rect(0, 0, width, height))
    setupYUVABuf(config, img.Y, img.Cb, img.Cr, img.A, img.YStride, img.CStride, img.AStride, ModeYUVA)
    return img
}

func setupRGBBuf(config *C.WebPDecoderConfig, pix []uint8, stride int, mode DecCspMode) {
    config.output.colorspace = C.WEBP_CSP_MODE(mode)
    config.output.is_external_memory = C.int(1)

    // Go represents a union as a byte array
    buf := (*C.WebPRGBABuffer)(unsafe.Pointer(&config.output.u[0]))
    buf.rgba, buf.size = bytesGetCPtr(pix)
    buf.stride = C.int(stride)
}

func setupYUVABuf(config *C.WebPDecoderConfig, y, u, v, a []uint8, yStride, uvStride, aStride int, mode DecCspMode) {
    config.output.colorspace = C.WEBP_CSP_MODE(mode)
    config.output.is_external_memory = C.int(1)

    // Go represents a union as a byte array
    buf := (*C.WebPYUVABuffer)(unsafe.Pointer(&config.output.u[0]))
    buf.y, buf.y_size = bytesGetCPtr(y)
    buf.u, buf.u_size = bytesGetCPtr(u)
    buf.v, buf.v_size = bytesGetCPtr(v)
    buf.y_stride = C.int(yStride)
    buf.u_stride = C.int(uvStride)
    buf.v_stride = C.int(uvStride)
    if a != nil {
        buf.a, buf.a_size = bytesGetCPtr(a)
        buf.a_stride = C.int(aStride)
    }
}

func calcOutputSize(config *C.WebPDecoderConfig) (w int, h int) {
    opts := config.options
    if int(opts.use_scaling) == 1 {
        w, h = int(opts.scaled_width), int(opts.scaled_height)
    } else if int(opts.use_cropping) == 1 {
        w, h = int(opts.crop_width), int(opts.crop_height)
    } else {
        w, h = int(config.input.width), int(config.input.height)
    }
    return
}

func init() {
    image.RegisterFormat("webp", "RIFF????WEBPVP8", Decode, DecodeConfig)
}
