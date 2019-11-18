// +build !cgo

package signal

// signalAction is fake for cgo
func signalAction(sigs ...os.Signal) {
}

// dumpSignalTo is fake for cgo
func dumpSignalTo(fd int) {
}

// dumpStacktraceTo is fake for cgo
func dumpStacktraceTo(name string) {
}

// dumpPreviousStacktrace is fake for cgo
func dumpPreviousStacktrace() {
}

// previousStacktrace is fake for cgo
func previousStacktrace() string {
}
