package signal

import "C"
import "os"

// SignalAction act as signal.Notify, which invokes the Go signal handler.
// https://godoc.org/os/signal#hdr-Go_programs_that_use_cgo_or_SWIG
func SignalAction(sigs ...os.Signal) {
	signalAction(sigs...)
}

// SignalDumpTo redirect log to fd, stdout if not set.
func SignalDumpTo(fd int) {
	signalDumpTo(fd)
}

// DumpBacktrace enables|disables log of bt when signal is triggered, disable if not set.
func DumpBacktrace(enable bool) {
	dumpBacktrace(enable)
}
