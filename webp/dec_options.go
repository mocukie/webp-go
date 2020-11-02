package webp

/*
#cgo LDFLAGS: -lwebp
#include <webp/decode.h>
*/
import "C"
import "image"

type DecCspMode int

const (
	ModeRGB   DecCspMode = C.MODE_RGB
	ModeRGBA  DecCspMode = C.MODE_rgbA
	ModeNRGBA DecCspMode = C.MODE_RGBA
	ModeYUV   DecCspMode = C.MODE_YUV
	ModeYUVA  DecCspMode = C.MODE_YUVA
)

type DecPixelFormat func(config *C.WebPDecoderConfig, w, h int) image.Image

var (
	// auto detect decoded image type, return webp.YCbCr/webp.NYCbCrA for lossy webp, otherwise webp.RGBImg/image.NRGBA
	TypeAuto DecPixelFormat = decPixAuto
	// auto detect decoded image type, return image.NRGBA if has alpha, otherwise image.RGBA
	TypeStd   DecPixelFormat = decPixStd
	TypeRGB   DecPixelFormat = decPixRGB
	TypeRGBA  DecPixelFormat = decPixRGBA
	TypeNRGBA DecPixelFormat = decPixNRGBA
	TypeYUV   DecPixelFormat = decPixYUV
	TypeYUVA  DecPixelFormat = decPixYUVA
)

type DecodeOptions struct {
	BypassFiltering        bool            // if true, skip the in-loop filtering
	NoFancyUpsampling      bool            // if true, use faster pointwise upsampler
	Crop                   image.Rectangle // do cropping if not empty, this is applied _first_
	Scale                  image.Rectangle // do scaling if not empty,  this is applied _afterward_
	UseThreads             bool            // if true, use multi-threaded decoding
	DitheringStrength      int             // dithering strength (0=Off, 100=full)
	Flip                   bool            // flip output vertically
	AlphaDitheringStrength int             // alpha dithering strength in [0..100]
	ImageType              DecPixelFormat  // decoded image type
}

func NewDecOptions() *DecodeOptions {
	return &DecodeOptions{ImageType: TypeAuto, UseThreads: true}
}

func (opts *DecodeOptions) assign(c *C.WebPDecoderOptions) {
	c.bypass_filtering = bool2CInt(opts.BypassFiltering)
	c.no_fancy_upsampling = bool2CInt(opts.NoFancyUpsampling)
	c.use_threads = bool2CInt(opts.UseThreads)
	c.dithering_strength = C.int(opts.DitheringStrength)
	c.flip = bool2CInt(opts.Flip)
	c.alpha_dithering_strength = C.int(opts.AlphaDitheringStrength)

	if !opts.Crop.Empty() {
		c.use_cropping = 1
		c.crop_left = C.int(opts.Crop.Min.X)
		c.crop_top = C.int(opts.Crop.Min.Y)
		c.crop_width = C.int(opts.Crop.Dx())
		c.crop_height = C.int(opts.Crop.Dy())
	}

	if !opts.Scale.Empty() {
		c.use_scaling = 1
		c.scaled_width = C.int(opts.Scale.Dx())
		c.scaled_height = C.int(opts.Scale.Dy())
	}
}

type BitStreamFormat int

const (
	FormatMixed BitStreamFormat = iota
	FormatLossy
	FormatLossless
)

type BitstreamFeatures struct {
	Width, Height int
	Format        BitStreamFormat
	HasAlpha      bool
	HasAnimation  bool
}
