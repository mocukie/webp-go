package webp

/*
#cgo LDFLAGS: -lwebp
#include <webp/encode.h>
#include <webp/decode.h>
*/
import "C"
import "fmt"

type VP8EncodeError int

const (
	VP8EncOk                        VP8EncodeError = C.VP8_ENC_OK
	VP8EncErrorOutOfMemory          VP8EncodeError = C.VP8_ENC_ERROR_OUT_OF_MEMORY
	VP8EncErrorBitstreamOutOfMemory VP8EncodeError = C.VP8_ENC_ERROR_BITSTREAM_OUT_OF_MEMORY
	VP8EncErrorNullParameter        VP8EncodeError = C.VP8_ENC_ERROR_NULL_PARAMETER
	VP8EncErrorInvalidConfiguration VP8EncodeError = C.VP8_ENC_ERROR_INVALID_CONFIGURATION
	VP8EncErrorBadDimension         VP8EncodeError = C.VP8_ENC_ERROR_BAD_DIMENSION
	VP8EncErrorPartition0Overflow   VP8EncodeError = C.VP8_ENC_ERROR_PARTITION0_OVERFLOW
	VP8EncErrorPartitionOverflow    VP8EncodeError = C.VP8_ENC_ERROR_PARTITION_OVERFLOW
	VP8EncErrorBadWrite             VP8EncodeError = C.VP8_ENC_ERROR_BAD_WRITE
	VP8EncErrorFileTooBig           VP8EncodeError = C.VP8_ENC_ERROR_FILE_TOO_BIG
	VP8EncErrorUserAbort            VP8EncodeError = C.VP8_ENC_ERROR_USER_ABORT
	VP8EncErrorLast                 VP8EncodeError = C.VP8_ENC_ERROR_LAST
)

func (e VP8EncodeError) Error() string {
	return "VP8 encode failed: " + e.String()
}

func (e VP8EncodeError) String() string {
	switch e {
	case VP8EncOk:
		return "ok"
	case VP8EncErrorOutOfMemory:
		return "memory error allocating objects"
	case VP8EncErrorBitstreamOutOfMemory:
		return "memory error while flushing bits"
	case VP8EncErrorNullParameter:
		return "a pointer parameter is NULL"
	case VP8EncErrorInvalidConfiguration:
		return "configuration is invalid"
	case VP8EncErrorBadDimension:
		return "picture has invalid width/height"
	case VP8EncErrorPartition0Overflow:
		return "partition is bigger than 512k"
	case VP8EncErrorPartitionOverflow:
		return "partition is bigger than 16M"
	case VP8EncErrorBadWrite:
		return "error while flushing bytes"
	case VP8EncErrorFileTooBig:
		return "file is bigger than 4G"
	case VP8EncErrorUserAbort:
		return "abort request by user"
	case VP8EncErrorLast:
		return "list terminator. always last."
	}
	return "unknown"
}

type VP8StatusCode int

const (
	VP8StatusOk                 VP8StatusCode = C.VP8_STATUS_OK
	VP8StatusOutOfMemory        VP8StatusCode = C.VP8_STATUS_OUT_OF_MEMORY
	VP8StatusInvalidParam       VP8StatusCode = C.VP8_STATUS_INVALID_PARAM
	VP8StatusBitstreamError     VP8StatusCode = C.VP8_STATUS_BITSTREAM_ERROR
	VP8StatusUnsupportedFeature VP8StatusCode = C.VP8_STATUS_UNSUPPORTED_FEATURE
	VP8StatusSuspended          VP8StatusCode = C.VP8_STATUS_SUSPENDED
	VP8StatusUserAbort          VP8StatusCode = C.VP8_STATUS_USER_ABORT
	VP8StatusNotEnoughData      VP8StatusCode = C.VP8_STATUS_NOT_ENOUGH_DATA
)

func (c VP8StatusCode) String() string {
	switch c {
	case VP8StatusOk:
		return "VP8_STATUS_OK"
	case VP8StatusOutOfMemory:
		return "VP8_STATUS_OUT_OF_MEMORY"
	case VP8StatusInvalidParam:
		return "VP8_STATUS_INVALID_PARAM"
	case VP8StatusBitstreamError:
		return "VP8_STATUS_BITSTREAM_ERROR"
	case VP8StatusUnsupportedFeature:
		return "VP8_STATUS_UNSUPPORTED_FEATURE"
	case VP8StatusSuspended:
		return "VP8_STATUS_SUSPENDED"
	case VP8StatusUserAbort:
		return "VP8_STATUS_USER_ABORT"
	case VP8StatusNotEnoughData:
		return "VP8_STATUS_NOT_ENOUGH_DATA"
	}
	return "UNKNOWN"
}

func (c VP8StatusCode) error(msg string) VP8DecodeError {
	return VP8DecodeError(fmt.Sprintf("%s return code <%s>", msg, c.String()))
}

type VP8DecodeError string

func (e VP8DecodeError) Error() string {
	return "VP8 decode failed :" + string(e)
}
