// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cgosymbolizer provides a cgo symbolizer based on libbacktrace.
// This will be used to provide a symbolic backtrace of cgo functions.
// This package does not export any symbols.
// To use it on all platforms, add a line like
//   import _ "github.com/searKing/golang/go/runtime/cgosymbolizer"
// 	 go build main.go
// Advance Usage can be set by go build tags:
//   BOOST_STACKTRACE_USE_WINDBG
//   BOOST_STACKTRACE_USE_WINDBG_CACHED
//   BOOST_STACKTRACE_USE_BACKTRACE
//   BOOST_STACKTRACE_USE_ADDR2LINE
//   BOOST_STACKTRACE_USE_NOOP
// all tags defined in https://www.boost.org/doc/libs/develop/doc/html/stacktrace/configuration_and_build.html
//   go build -tags BOOST_STACKTRACE_USE_BACKTRACE main.go
// somewhere in your program.
// for linux only, you can use `cgosymbolizer` by ianlancetaylor instead.
//   import _ "github.com/ianlancetaylor/cgosymbolizer"
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
