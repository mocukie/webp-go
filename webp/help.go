package webp

/*
#cgo LDFLAGS: -lwebp
#include <webp/types.h>
*/
import "C"
import (
	"fmt"
	"reflect"
	"unsafe"
)

func bool2CInt(b bool) C.int {
	if b {
		return C.int(1)
	} else {
		return C.int(0)
	}
}

func bytesGetCPtr(data []byte) (*C.uint8_t, C.size_t) {
	return (*C.uint8_t)(unsafe.Pointer(&data[0])), C.size_t(len(data))
}

type Version int

func (ver Version) V() (major, minor, revision int) {
	v := int(ver)
	major, minor, revision = (v>>16)&0xff, (v>>8)&0xff, v&0xff
	return
}

func (ver Version) Major() int {
	return int(ver) >> 16 & 0xff
}

func (ver Version) Minor() int {
	return int(ver) >> 8 & 0xff
}

func (ver Version) Revision() int {
	return int(ver) & 0xff
}

func (ver Version) String() string {
	v := int(ver)
	return fmt.Sprintf("%d.%d.%d", (v>>16)&0xff, (v>>8)&0xff, v&0xff)
}

type unsafeBytes []byte

func allocUnsafeBytes(size int) unsafeBytes {
	ptr := unsafe.Pointer(C.WebPMalloc(C.size_t(size)))
	if ptr == nil {
		return nil
	}
	s := &reflect.SliceHeader{
		Data: uintptr(ptr),
		Len:  size,
		Cap:  size,
	}
	return *(*unsafeBytes)(unsafe.Pointer(s))
}

func wrapUnsafeBytes(data *C.uint8_t, size C.size_t) unsafeBytes {
	s := &reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(data)),
		Len:  int(size),
		Cap:  int(size),
	}
	return *(*unsafeBytes)(unsafe.Pointer(s))
}

func (data unsafeBytes) pointer() unsafe.Pointer {
	if data == nil {
		return nil
	}
	s := (*reflect.SliceHeader)(unsafe.Pointer(&data))
	return unsafe.Pointer(s.Data)
}

func (data unsafeBytes) asSafe() []byte {
	if data == nil {
		return nil
	}
	defer data.release()
	safe := make([]byte, len(data))
	copy(safe, data)
	return safe
}

func (data unsafeBytes) release() {
	if data == nil {
		return
	}
	s := (*reflect.SliceHeader)(unsafe.Pointer(&data))
	C.WebPFree(unsafe.Pointer(s.Data))
}
