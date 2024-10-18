//go:build uritawebview

package main

/*
// NOTE: -L../../build is required for compilation, -L. is for deployment
#cgo LDFLAGS: -L../../build -L. -lurita

#include <stdlib.h>
extern int uritaOpenWv(const char* url);
*/
import "C"

import (
	"unsafe"
)

func uritaOpenWv(url string) {
	cUrl := C.CString(url)
	defer C.free(unsafe.Pointer(cUrl))

	C.uritaOpenWv(cUrl)
}
