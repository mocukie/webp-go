package colorx

import "image/color"

func RGBToBT601(r, g, b uint8) (uint8, uint8, uint8) {
    /* ITU-R BT.601
       Y' =  0.2569*R + 0.5044*G + 0.0979*B + 16
       Cb = -0.1483*R - 0.2911*G + 0.4394*B + 128
       Cr =  0.4394*R - 0.3679*G - 0.0715*B + 128
    */

    R := int32(r)
    G := int32(g)
    B := int32(b)
    y := (16836*R+33056*G+6416*B+32768)>>16 + 16
    cb := (-9718*R-19078*G+28797*B+32768)>>16 + 128
    cr := (28797*R-24111*G-4686*B+32768)>>16 + 128

    return uint8(y), uint8(cb), uint8(cr)
}

func BT601ToRGB(y, cb, cr uint8) (uint32, uint32, uint32) {
    /* ITU-R BT.601
       R = 1.164 * (Y-16)                    + 1.596 * (Cr-128)
       G = 1.164 * (Y-16) - 0.392 * (Cb-128) - 0.813 * (Cr-128)
       B = 1.164 * (Y-16) + 2.017 * (Cb-128)
    */

    Y := int32(y-16)*(76284) + 32768
    U := int32(cb) - 128
    V := int32(cr) - 128

    r := Y + 104595*V
    if uint32(r)&0xff000000 == 0 {
        r >>= 16
    } else {
        r = ^(r >> 31) & 0xff
    }

    g := Y - 25690*U - 53283*V
    if uint32(g)&0xff000000 == 0 {
        g >>= 16
    } else {
        g = ^(g >> 31) & 0xff
    }

    b := Y + 132186*U
    if uint32(b)&0xff000000 == 0 {
        b >>= 16
    } else {
        b = ^(b >> 31) & 0xff
    }

    return uint32(r<<8 | r), uint32(g<<8 | g), uint32(b<<8 | b)
}

type YCbCrBT601 struct {
    Y, Cb, Cr uint8
}

func (c YCbCrBT601) RGBA() (uint32, uint32, uint32, uint32) {
    r, g, b := BT601ToRGB(c.Y, c.Cb, c.Cr)
    return r, g, b, 0xffff
}

var YCbCrBT601Model color.Model = color.ModelFunc(yCbCrBT601Model)

func yCbCrBT601Model(c color.Color) color.Color {
    if _, ok := c.(YCbCrBT601); ok {
        return c
    }
    r, g, b, _ := c.RGBA()
    y, u, v := RGBToBT601(uint8(r>>8), uint8(g>>8), uint8(b>>8))
    return YCbCrBT601{y, u, v}
}

type NYCbCrBT601 struct {
    YCbCrBT601
    A uint8
}

func (c NYCbCrBT601) RGBA() (uint32, uint32, uint32, uint32) {
    r, g, b, _ := c.YCbCrBT601.RGBA()
    a := uint32(c.A)<<8 | uint32(c.A)
    if a != 0xffff {
        r = r * a / 0xffff
        g = g * a / 0xffff
        b = b * a / 0xffff
    }

    return r, g, b, a
}

var NYCbCrBT601Model = color.ModelFunc(nYCbCrBT601Model)

func nYCbCrBT601Model(c color.Color) color.Color {
    switch c := c.(type) {
    case NYCbCrBT601:
        return c
    case YCbCrBT601:
        return NYCbCrBT601{c, 0xff}
    }

    r, g, b, a := c.RGBA()
    if a != 0 {
        r = (r * 0xffff) / a
        g = (g * 0xffff) / a
        b = (b * 0xffff) / a
    }

    y, u, v := RGBToBT601(uint8(r>>8), uint8(g>>8), uint8(b>>8))
    return NYCbCrBT601{YCbCrBT601{Y: y, Cb: u, Cr: v}, uint8(a >> 8)}
}

var RGBModel = color.ModelFunc(rgbModel)

func rgbModel(c color.Color) color.Color {
    if _, ok := c.(RGB); ok {
        return c
    }
    r, g, b, _ := c.RGBA()
    return RGB{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)}
}

type RGB struct {
    R, G, B uint8
}

func (c RGB) RGBA() (r, g, b, a uint32) {
    r = uint32(c.R)
    r |= r << 8
    g = uint32(c.G)
    g |= g << 8
    b = uint32(c.B)
    b |= b << 8
    a = uint32(0xffff)
    return
}
