// +build cgo

package signal


/*
#cgo pkg-config: ${SRCDIR}/cgo/pkgconfig/libsignal_cgo.pc
#include <stdlib.h>  // Needed for C.free
#include <stdio.h>
#include "signal.cgo.h"
*/
import "C"
import (
	"os"
)

// signalAction act as signal.Notify, which invokes the Go signal handler.
// https://godoc.org/os/signal#hdr-Go_programs_that_use_cgo_or_SWIG
func signalAction(sigs ...os.Signal) {
	for _, sig := range sigs {
		C.CGOSignalHandlerSignalAction(C.int(Signum(sig)))
	}
}

// signalDumpTo redirect log to fd, stdout if not set.
func signalDumpTo(fd int) {
	C.CGOSignalHandlerSetFd(C.int(fd))
}

// dumpBacktrace enables|disables log of bt when signal is triggered, disable if not set.
func dumpBacktrace(enable bool) {
	C.CGOSignalHandlerSetBacktraceDump(C._Bool(enable))
}
