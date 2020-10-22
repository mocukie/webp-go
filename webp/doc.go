/*
Package webp provides a cgo binding of libwebp to process WebP image.

Require

gcc

libwebp 1.1.0

Usage

Encode lossy
    var opts *webp.EncodeOptions
    //default options
    opts, _ = webp.NewEncOptions()
    //or init with preset
    opts, _ = webp.NewEncOptionsByPreset(webp.PresetPicture, webp.LossyDefaultQuality)

    var buf = bytes.NewBuffer(nil)
    err = webp.Encode(buf, img, opts)
    if err != nil {
        panic(err)
    }
    ioutil.WriteFile("foo_lossy.webp", buf.Bytes(), os.ModePerm)


Encode lossless
    opts, _ := webp.NewEncOptions()
    opts.SetupLosslessPreset(webp.LosslessDefaultLevel)
    //or just
    opts.Lossless = true
    opts.Quality = webp.LosslessDefaultQuality
    //set true if you want preserve RGB values under transparent area
    opts.Exact = true

    var buf = bytes.NewBuffer(nil)
    err = webp.Encode(buf, img, opts)
    if err != nil {
        panic(err)
    }
    ioutil.WriteFile("foo_lossless.webp", buf.Bytes(), os.ModePerm)


Decode
    fin, _ := os.Open("foo.webp")
    webpImg, err := webp.Decode(fin)
    if err != nil {
        panic(err)
    }
    fin.Seek(0, io.SeekStart)
    //decode with options
    decOpts := webp.NewDecOptions()
    decOpts.ImageType = webp.TypeNRGBA //decode as image.NRGBA
    webpImg, err = webp.DecodeEX(fin, options)
    if err != nil {
        panic(err)
    }


Get and Set metadata chunk
    iccp, err := webp.GetMetadata(webpData, webp.ICCP)
    if err != nil {
        newWebpData, err := webp.SetMetadata(webpdata2, webp.ICCP, iccp)
    }
*/
package webp
