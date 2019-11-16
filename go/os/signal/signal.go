// +build aix darwin dragonfly freebsd js,wasm linux nacl netbsd openbsd solaris windows

package signal

/*
#cgo linux,amd64 pkg-config: ${SRCDIR}/cgo/pkgconfig/libsignal_cgo.linux.amd64.pc
#cgo darwin,amd64 pkg-config: ${SRCDIR}/cgo/pkgconfig/libsignal_cgo.darwin.amd64.pc
#include <stdlib.h>  // Needed for C.free
#include <stdio.h>
#include "signal.cgo.h"
*/
import "C"
import (
	"os"
)

func SignalAction(sigs ...os.Signal) {
	for _, sig := range sigs {
		C.CGOSignalHandlerSignalAction(C.int(Signum(sig)))
	}
}

func SignalDumpTo(fd int) {
	C.CGOSignalHandlerSetFd(C.int(fd))
}

func DumpBacktrace(enable bool) {
	C.CGOSignalHandlerSetBacktraceDump(C._Bool(enable))
}
