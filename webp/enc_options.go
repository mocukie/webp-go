package webp

/*
#cgo LDFLAGS: -lwebp
#include "webp.h"
*/
import "C"

const LosslessDefaultLevel int = 6
const LosslessDefaultQuality float32 = 70.0
const LossyDefaultQuality float32 = 75.0

type EncodePreset int

const (
	PresetDefault EncodePreset = C.WEBP_PRESET_DEFAULT // default preset.
	PresetPicture EncodePreset = C.WEBP_PRESET_PICTURE // digital picture, like portrait, inner shot
	PresetPhoto   EncodePreset = C.WEBP_PRESET_PHOTO   // outdoor photograph, with natural lighting
	PresetDrawing EncodePreset = C.WEBP_PRESET_DRAWING // hand or line drawing, with high-contrast details
	PresetIcon    EncodePreset = C.WEBP_PRESET_ICON    // small-sized colorful images
	PresetText    EncodePreset = C.WEBP_PRESET_TEXT    // text-like
)

type ImageHint int

const (
	HintDefault ImageHint = C.WEBP_HINT_DEFAULT // default preset.
	HintPicture ImageHint = C.WEBP_HINT_PICTURE // digital picture, like portrait, inner shot
	HintPhoto   ImageHint = C.WEBP_HINT_PHOTO   // outdoor photograph, with natural lighting
	HintGraph   ImageHint = C.WEBP_HINT_GRAPH   // Discrete tone image (graph, map-tile etc).
	HintLast    ImageHint = C.WEBP_HINT_LAST
)

type FilterType int

const (
	SimpleFilter FilterType = 0
	StrongFilter FilterType = 1
)

type AlphaFilter int

const (
	NoneAlphaFilter AlphaFilter = iota
	FastAlphaFilter
	BestAlphaFilter
)

type EncodeOptions struct {
	Lossless bool
	// between 0 and 100. For lossy, 0 gives the smallest
	// size and 100 the largest. For lossless, this
	// parameter is the amount of effort put into the
	// compression: 0 is the fastest but gives larger
	// files compared to the slowest, but best, 100.
	Quality float32
	// quality/speed trade-off (0=fast, 6=slower-better)
	Method int
	// Hint for image type (lossless only for now)
	ImageHint ImageHint
	// if non-zero, set the desired target size in bytes. Takes precedence over the 'compression' parameter.
	TargetSize int
	// if non-zero, specifies the minimal distortion to try to achieve. Takes precedence over target_size
	TargetPSNR float32
	// maximum number of segments to use, in [1..4]
	Segments int
	// Spatial Noise Shaping. 0=off, 100=maximum
	SnsStrength int
	// range: [0 = off .. 100 = strongest]
	FilterStrength int
	// range: [0 = off .. 7 = least sharp]
	FilterSharpness int
	// filtering type: 0 = simple, 1 = strong (only used if filter_strength > 0 or autofilter > 0)
	FilterType FilterType
	// Auto adjust filter's strength
	AutoFilter bool
	// Algorithm for encoding the alpha plane (0 = none, 1 = compressed with WebP lossless). Default is 1
	AlphaCompression int
	// Predictive filtering method for alpha plane. 0: none, 1: fast, 2: best. Default if 1.
	AlphaFiltering AlphaFilter
	// Between 0 (smallest size) and 100 (lossless). Default is 100.
	AlphaQuality int
	// number of entropy-analysis passes (in [1..10]).
	Pass int
	// if true, export the compressed picture back. In-loop filtering is not applied.
	ShowCompressed bool
	// preprocessing filter: 0=none, 1=segment-smooth, 2=pseudo-random dithering
	Preprocessing int
	// log2(number of token partitions) in [0..3]. Default is set to 0 for easier progressive decoding.
	Partitions int
	// quality degradation allowed to fit the 512k limit on prediction modes coding (0: no degradation, 100: maximum possible degradation).
	PartitionLimit int
	// If true, compression parameters will be remapped
	// to better match the expected output size from
	// JPEG compression. Generally, the output size will
	// be similar but the degradation will be lower.
	EmulateJpegSize bool
	// If true, try and use multi-threaded encoding.
	ThreadLevel bool
	// If set, reduce memory usage (but increase CPU use).
	LowMemory bool
	// Near lossless encoding [0 = max loss .. 100 = off
	// (default)].
	NearLossless int
	// if true, preserve the exact RGB values under
	// transparent area. Otherwise, discard this invisible
	// RGB information for better compression. The default
	// value is 0.
	Exact bool
	// reserved for future lossless feature
	UseDeltaPalette int
	// if needed, use sharp (and slow) RGB->YUV conversion
	UseSharpYUV bool
}

func (opts *EncodeOptions) from(c *C.WebPConfig) {
	opts.Lossless = int(c.lossless) == 1
	opts.Quality = float32(c.quality)
	opts.Method = int(c.method)
	opts.ImageHint = ImageHint(c.image_hint)
	opts.TargetSize = int(c.target_size)
	opts.TargetPSNR = float32(c.target_PSNR)
	opts.Segments = int(c.segments)
	opts.SnsStrength = int(c.sns_strength)
	opts.FilterStrength = int(c.filter_strength)
	opts.FilterSharpness = int(c.filter_sharpness)
	opts.FilterType = FilterType(c.filter_type)
	opts.AutoFilter = int(c.autofilter) == 1
	opts.AlphaCompression = int(c.alpha_compression)
	opts.AlphaFiltering = AlphaFilter(c.alpha_filtering)
	opts.AlphaQuality = int(c.alpha_quality)
	opts.Pass = int(c.pass)
	opts.ShowCompressed = int(c.show_compressed) == 1
	opts.Preprocessing = int(c.preprocessing)
	opts.Partitions = int(c.partitions)
	opts.PartitionLimit = int(c.partition_limit)
	opts.EmulateJpegSize = int(c.emulate_jpeg_size) == 1
	opts.ThreadLevel = int(c.thread_level) == 1
	opts.LowMemory = int(c.low_memory) == 1
	opts.NearLossless = int(c.near_lossless)
	opts.Exact = int(c.exact) == 1
	opts.UseDeltaPalette = int(c.use_delta_palette) // reserved for future lossless feature
	opts.UseSharpYUV = int(c.use_sharp_yuv) == 1
}

func (opts *EncodeOptions) assign(c *C.WebPConfig) {
	c.lossless = bool2CInt(opts.Lossless)
	c.quality = C.float(opts.Quality)
	c.method = C.int(opts.Method)
	c.image_hint = C.WebPImageHint(opts.ImageHint)
	c.target_size = C.int(opts.TargetSize)
	c.target_PSNR = C.float(opts.TargetPSNR)
	c.segments = C.int(opts.Segments)
	c.sns_strength = C.int(opts.SnsStrength)
	c.filter_strength = C.int(opts.FilterStrength)
	c.filter_sharpness = C.int(opts.FilterSharpness)
	c.filter_type = C.int(opts.FilterType)
	c.autofilter = bool2CInt(opts.AutoFilter)
	c.alpha_compression = C.int(opts.AlphaCompression)
	c.alpha_filtering = C.int(opts.AlphaFiltering)
	c.alpha_quality = C.int(opts.AlphaQuality)
	c.pass = C.int(opts.Pass)
	c.show_compressed = bool2CInt(opts.ShowCompressed)
	c.preprocessing = C.int(opts.Preprocessing)
	c.partitions = C.int(opts.Partitions)
	c.partition_limit = C.int(opts.PartitionLimit)
	c.emulate_jpeg_size = bool2CInt(opts.EmulateJpegSize)
	c.thread_level = bool2CInt(opts.ThreadLevel)
	c.low_memory = bool2CInt(opts.LowMemory)
	c.near_lossless = C.int(opts.NearLossless)
	c.exact = bool2CInt(opts.Exact)
	c.use_delta_palette = C.int(opts.UseDeltaPalette)
	c.use_sharp_yuv = bool2CInt(opts.UseSharpYUV)
}

func NewEncOptions() (*EncodeOptions, error) {
	var c C.WebPConfig
	ret := C.WebPConfigInit(&c)
	if int(ret) == 0 {
		return nil, VP8EncErrorInvalidConfiguration
	}
	opts := new(EncodeOptions)
	opts.from(&c)
	return opts, nil
}

func NewEncOptionsByPreset(preset EncodePreset, quality float32) (*EncodeOptions, error) {
	var c C.WebPConfig
	ret := C.WebPConfigPreset(&c, C.WebPPreset(preset), C.float(quality))
	if int(ret) == 0 {
		return nil, VP8EncErrorInvalidConfiguration
	}
	opts := new(EncodeOptions)
	opts.from(&c)
	return opts, nil
}

// Activate the lossless compression mode with the desired efficiency level
// between 0 (fastest, lowest compression) and 9 (slower, best compression).
// A good default level is '6', providing a fair tradeoff between compression
// speed and final compressed size.
// This function will overwrite several fields from config: 'method', 'quality'
// and 'lossless'.
func (opts *EncodeOptions) SetupLosslessPreset(level int) error {
	var c C.WebPConfig
	opts.assign(&c)
	ret := C.WebPConfigLosslessPreset(&c, C.int(level))
	if int(ret) == 0 {
		return VP8EncErrorInvalidConfiguration
	}
	opts.from(&c)
	return nil
}

func (opts *EncodeOptions) Validate() error {
	var c C.WebPConfig
	opts.assign(&c)
	if !validateEncodeConfig(&c) {
		return VP8EncErrorInvalidConfiguration
	}
	return nil
}

func validateEncodeConfig(config *C.WebPConfig) bool {
	ret := C.WebPValidateConfig(config)
	return int(ret) != 0
}
