// +build !cgo

package signal

// signalAction is fake for cgo
func signalAction(sigs ...os.Signal) {
}

// signalAction is fake for cgo
func signalDumpTo(fd int) {
}

// signalAction is fake for cgo
func dumpBacktrace(enable bool) {
}
