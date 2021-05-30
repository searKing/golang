// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cgosymbolizer provides a cgo symbolizer based on libbacktrace.
// This will be used to provide a symbolic backtrace of cgo functions.
// This package does not export any symbols.
// To use it, add a line like
//   import _ "github.com/ianlancetaylor/cgosymbolizer"
// somewhere in your program.
package cgosymbolizer

// extern void cgoTraceback(void*);
// extern void cgoSymbolizer(void*);
import "C"

import (
	"runtime"
	"unsafe"
)

func init() {
	runtime.SetCgoTraceback(0, unsafe.Pointer(C.cgoTraceback), nil, unsafe.Pointer(C.cgoSymbolizer))
}

//// PreviousStacktrace returns a human readable stacktrace
//func PreviousStacktrace() string {
//	stacktraceChars := C.CGO_PreviousStacktrace()
//	defer C.free(unsafe.Pointer(stacktraceChars))
//	return C.GoString(stacktraceChars)
//}
