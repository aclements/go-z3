// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package z3

import "unsafe"

/*
#cgo LDFLAGS: -lz3
#include <z3.h>
#include <stdlib.h>
*/
import "C"

// LogOpen creates a Z3 interaction log in a file called filename.
//
// The interaction log is a low-level trace of all Z3 API calls.
//
// It returns false if it fails to open the log.
func LogOpen(filename string) bool {
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))
	return C.Z3_open_log(cfilename) != 0
}

// LogAppend emits text to the Z3 interaction log.
func LogAppend(text string) {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.Z3_append_log(ctext)
}

// LogClose closes the Z3 interaction log file.
func LogClose() {
	C.Z3_close_log()
}
