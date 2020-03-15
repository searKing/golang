// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build cgo

package cgo

/*
	#cgo CXXFLAGS: -I${SRCDIR}/include/
	#cgo windows CXXFLAGS: -g
	#cgo !windows CXXFLAGS: -g -D_GNU_SOURCE
	#cgo linux LDFLAGS: -ldl

	#include "signal.cgo.h"
	#include <stdio.h>
   	#include <stdlib.h>  // Needed for C.free

*/
import "C"
import (
	"unsafe"

	_ "github.com/searKing/golang/go/os/signal/cgo/include"
)

// signalAction act as signal.Notify, which invokes the Go signal handler.
// https://godoc.org/os/signal#hdr-Go_programs_that_use_cgo_or_SWIG
func SetSig(sig int) {
	C.CGO_SignalHandlerSetSig(C.int(sig))
}

// SetSignalDumpToFd redirect log to fd, -1 if not set; muted if < 0.
func SetSignalDumpToFd(fd int) { C.CGO_SignalHandlerSetSignalDumpToFd(C.int(fd)) }

// SetBacktraceDumpToFile set dump file path of stacktrace when signal is triggered, nop if not set.
func SetBacktraceDumpToFile(name string) {
	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))
	C.CGO_SignalHandlerSetStacktraceDumpToFile(cs)
}

// DumpPreviousStacktrace dumps human readable stacktrace to fd, which is set by SetSignalDumpToFd.
func DumpPreviousStacktrace() {
	C.CGO_SignalHandlerDumpPreviousStacktrace()
}

// PreviousStacktrace returns a human readable stacktrace
func PreviousStacktrace() string {
	stacktraceChars := C.CGO_PreviousStacktrace()
	defer C.free(unsafe.Pointer(stacktraceChars))
	return C.GoString(stacktraceChars)
}

// PreviousStacktrace sets a rule to raise signal to {to} and wait until {wait}, done with sleep {sleepInSeconds}s
func SetSigInvokeChain(from, to, wait, sleepInSeconds int) {
	C.CGO_SetSigInvokeChain(C.int(from), C.int(to), C.int(wait), C.int(sleepInSeconds))
}
