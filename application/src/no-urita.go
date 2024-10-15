// +build nowebview

package main

/*
// NOTE: -L../../build is required for compilation, -L. is for deployment
#cgo LDFLAGS: -L../../build -L. -lurita

#include <stdlib.h>
extern int uritaOpenWv(const char* url);
*/


func uritaOpenWv(url string) {
	
}
