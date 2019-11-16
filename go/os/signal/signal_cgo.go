// +build cgo

package signal

import "C"
import (
	"os"

	"github.com/searKing/golang/go/os/signal/cgo"
)

// signalAction act as signal.Notify, which invokes the Go signal handler.
// https://godoc.org/os/signal#hdr-Go_programs_that_use_cgo_or_SWIG
func signalAction(sigs ...os.Signal) {
	for _, sig := range sigs {
		cgo.SignalAction(Signum(sig))
	}
}

// signalDumpTo redirect log to fd, stdout if not set.
func signalDumpTo(fd int) {
	cgo.SetFd(fd)
}

// dumpBacktrace enables|disables log of bt when signal is triggered, disable if not set.
func dumpBacktrace(enable bool) {
	cgo.SetBacktraceDump(enable)
}
