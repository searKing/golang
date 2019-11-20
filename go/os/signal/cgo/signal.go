// +build cgo

package cgo

/*
	#cgo CXXFLAGS: -I${SRCDIR}/include/
	#cgo darwin CXXFLAGS: -g -D_GNU_SOURCE
	#cgo !darwin CXXFLAGS: -g
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
func SignalAction(sig int) {
	C.CGOSignalHandlerSignalAction(C.int(sig))
}

// SetSignalDumpToFd redirect log to fd, -1 if not set; muted if < 0.
func SetSignalDumpToFd(fd int) { C.CGOSignalHandlerSetSignalDumpToFd(C.int(fd)) }

// SetBacktraceDumpToFile set dump file path of stacktrace when signal is triggered, nop if not set.
func SetBacktraceDumpToFile(name string) {
	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))
	C.CGOSignalHandlerSetStacktraceDumpToFile(cs)
}

// DumpPreviousStacktrace dumps human readable stacktrace to fd, which is set by SetSignalDumpToFd.
func DumpPreviousStacktrace() {
	C.CGOSignalHandlerDumpPreviousHumanReadableStacktrace()
}

// PreviousStacktrace returns a human readable stacktrace
func PreviousStacktrace() string {
	stacktraceChars := C.CGOPreviousHumanReadableStacktrace()
	defer C.free(unsafe.Pointer(stacktraceChars))
	return C.GoString(stacktraceChars)
}
