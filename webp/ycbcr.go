package webp

import "C"

import (
    "github.com/mocukie/webp-go/webp/colorx"
    "image"
    "image/color"
)

//webp yuv420
type YCbCr struct {
    image.YCbCr
}

func (p *YCbCr) ColorModel() color.Model {
    return colorx.YCbCrBT601Model
}

func (p *YCbCr) At(x, y int) color.Color {
    return p.YCbCrAt(x, y)
}

func (p *YCbCr) YCbCrAt(x, y int) colorx.YCbCrBT601 {
    if !(image.Point{x, y}.In(p.Rect)) {
        return colorx.YCbCrBT601{Y: 16, Cb: 16, Cr: 16}
    }
    yi := p.YOffset(x, y)
    ci := p.COffset(x, y)
    return colorx.YCbCrBT601{
        Y:  p.Y[yi],
        Cb: p.Cb[ci],
        Cr: p.Cr[ci],
    }
}

func (p *YCbCr) Set(x, y int, c color.Color) {
    if !(image.Point{x, y}.In(p.Rect)) {
        return
    }
    ycbcr := colorx.YCbCrBT601Model.Convert(c).(colorx.YCbCrBT601)
    yi := p.YOffset(x, y)
    ci := p.COffset(x, y)
    p.Y[yi] = ycbcr.Y
    p.Cb[ci] = ycbcr.Cb
    p.Cr[ci] = ycbcr.Cr
}

func NewYCbCr(r image.Rectangle) *YCbCr {
    return &YCbCr{
        YCbCr: *image.NewYCbCr(r, image.YCbCrSubsampleRatio420),
    }
}

//webp yuva420
type NYCbCrA struct {
    image.NYCbCrA
}

func (p *NYCbCrA) ColorModel() color.Model {
    return colorx.NYCbCrBT601Model
}

func (p *NYCbCrA) At(x, y int) color.Color {
    return p.NYCbCrAAt(x, y)
}

func (p *NYCbCrA) NYCbCrAAt(x, y int) colorx.NYCbCrBT601 {
    if !(image.Point{X: x, Y: y}.In(p.Rect)) {
        return colorx.NYCbCrBT601{YCbCrBT601: colorx.YCbCrBT601{Y: 16, Cb: 16, Cr: 16}}
    }
    yi := p.YOffset(x, y)
    ci := p.COffset(x, y)
    ai := p.AOffset(x, y)
    return colorx.NYCbCrBT601{
        YCbCrBT601: colorx.YCbCrBT601{
            Y:  p.Y[yi],
            Cb: p.Cb[ci],
            Cr: p.Cr[ci],
        },
        A: p.A[ai],
    }
}

func NewNYCbCrA(r image.Rectangle, ) *NYCbCrA {
    return &NYCbCrA{
        NYCbCrA: *image.NewNYCbCrA(r, image.YCbCrSubsampleRatio420),
    }
}
