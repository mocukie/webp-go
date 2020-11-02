package webp

import (
	"github.com/mocukie/webp-go/webp/colorx"
	"image"
	"image/color"
)

type RGBImg struct {
	// Pix holds the image's pixels, in R, G, B order.
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect image.Rectangle
}

func (p *RGBImg) ColorModel() color.Model { return colorx.RGBModel }

func (p *RGBImg) Bounds() image.Rectangle { return p.Rect }

func (p *RGBImg) At(x, y int) color.Color {
	return p.RGBAAt(x, y)
}

func (p *RGBImg) RGBAAt(x, y int) colorx.RGB {
	if !(image.Point{X: x, Y: y}.In(p.Rect)) {
		return colorx.RGB{}
	}
	off := (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*3
	return colorx.RGB{
		R: p.Pix[off+0],
		G: p.Pix[off+1],
		B: p.Pix[off+2],
	}
}

// NewRGB returns a new RGBImg with the given bounds.
func NewRGB(r image.Rectangle) *RGBImg {
	w, h := r.Dx(), r.Dy()
	return &RGBImg{
		Pix:    make([]uint8, w*h*3),
		Stride: 3 * w,
		Rect:   r,
	}
}

//encode helper for someone want to import custom image type self
type ARGBImg struct {
	// Pix holds the image's pixels, store A,R,G,B from hi to lo in uint32
	Pix []uint32
	// Stride is the Pix stride in pixels units, not bytes between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect image.Rectangle
}

func (p *ARGBImg) ColorModel() color.Model { return color.NRGBAModel }

func (p *ARGBImg) Bounds() image.Rectangle { return p.Rect }

func (p *ARGBImg) At(x, y int) color.Color {
	if !(image.Point{X: x, Y: y}.In(p.Rect)) {
		return color.NRGBA{}
	}
	argb := p.Pix[(y-p.Rect.Min.Y)*p.Stride+(x-p.Rect.Min.X)*4]
	return color.NRGBA{
		R: uint8(argb >> 16 & 0xff),
		G: uint8(argb >> 8 & 0xff),
		B: uint8(argb & 0xff),
		A: uint8(argb >> 24 & 0xff),
	}
}

func NewARGB(r image.Rectangle) *ARGBImg {
	return &ARGBImg{
		Pix:    make([]uint32, r.Dx()*r.Dy()),
		Stride: r.Dx(),
		Rect:   r,
	}
}

var _ image.Image = (*RGBImg)(nil)
var _ image.Image = (*ARGBImg)(nil)
