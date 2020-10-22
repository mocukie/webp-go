package webp

/*
#cgo LDFLAGS: -lwebp -lwebpmux -lwebpdemux
#include "webp.h"
*/
import "C"
import (
    "unsafe"
)

type MuxError int

const (
    MuxOk              MuxError = C.WEBP_MUX_OK
    MuxNotFound        MuxError = C.WEBP_MUX_NOT_FOUND
    MuxInvalidArgument MuxError = C.WEBP_MUX_INVALID_ARGUMENT
    MuxBadData         MuxError = C.WEBP_MUX_BAD_DATA
    MuxMemoryError     MuxError = C.WEBP_MUX_MEMORY_ERROR
    MuxNotEnoughData   MuxError = C.WEBP_MUX_NOT_ENOUGH_DATA
)

func (c MuxError) Error() string {
    switch c {
    case MuxOk:
        return "WebPMuxError: ok"
    case MuxNotFound:
        return "WebPMuxError: not found"
    case MuxInvalidArgument:
        return "WebPMuxError: invalid argument"
    case MuxBadData:
        return "WebPMuxError: bad data"
    case MuxMemoryError:
        return "WebPMuxError: memory error"
    case MuxNotEnoughData:
        return "WebPMuxError: not enough data "
    }
    return "WebPMuxError: unknown"
}

type FourCC [4]byte

var (
    ICCP = FourCC{'I', 'C', 'C', 'P'}
    XMP  = FourCC{'X', 'M', 'P', ' '}
    EXIF = FourCC{'E', 'X', 'I', 'F'}
)

func SetMetadata(img []byte, fourcc FourCC, chunk []byte) ([]byte, error) {
    imgPtr, imgSize := bytesGetCPtr(img)
    chunkPtr, chunkSize := bytesGetCPtr(chunk)
    var out C.WebPData
    if err := MuxError(C.GoSetWebPChunk(imgPtr, imgSize, (*C.char)(unsafe.Pointer(&fourcc)), chunkPtr, chunkSize, &out)); err != MuxOk {
        return nil, err
    }
    return wrapUnsafeBytes(out.bytes, out.size).asSafe(), nil
}

// return MuxNotFound if meta not exists
func GetMetadata(img []byte, fourcc FourCC) ([]byte, error) {
    imgPtr, imgSize := bytesGetCPtr(img)
    var cOff C.int
    var cSize C.size_t
    if err := MuxError(C.GoGetWebPChunk(imgPtr, imgSize, (*C.char)(unsafe.Pointer(&fourcc)), &cOff, &cSize)); err != MuxOk {
        return nil, err
    }
    off, size := int(cOff), int(cSize)
    return img[off : off+size : off+size], nil
}

// return MuxNotFound if meta not exists
func DeleteMetadata(img []byte, fourcc FourCC) error {
    imgPtr, imgSize := bytesGetCPtr(img)
    if err := MuxError(C.GoDeleteWebPChunk(imgPtr, imgSize, (*C.char)(unsafe.Pointer(&fourcc)))); err != MuxOk {
        return err
    }
    return nil
}
