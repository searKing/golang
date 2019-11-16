// +build cgo

package cgo

/*
   #cgo LDFLAGS: -ldl
   #include "signal.cgo.h"
   #include <stdio.h>
   #include <stdlib.h>  // Needed for C.free
*/
import "C"

// signalAction act as signal.Notify, which invokes the Go signal handler.
// https://godoc.org/os/signal#hdr-Go_programs_that_use_cgo_or_SWIG
func
SignalAction(sig int) {
	C.CGOSignalHandlerSignalAction(C.int(sig))
}

// signalDumpTo redirect log to fd, stdout if not set.
func SetFd(fd int) { C.CGOSignalHandlerSetFd(C.int(fd)) }

// dumpBacktrace enables|disables log of bt when signal is triggered, disable if
// not set.
func SetBacktraceDump(enable bool) {
	C.CGOSignalHandlerSetBacktraceDump(C._Bool(enable))
}
