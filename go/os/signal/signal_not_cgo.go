// +build !cgo

package signal

// setSig is fake for cgo
func setSig(sigs ...os.Signal) {
}

// dumpSignalTo is fake for cgo
func dumpSignalTo(fd int) {
}

// dumpStacktraceTo is fake for cgo
func dumpStacktraceTo(name string) {
}

func registerOnSignal(onSignal OnSignalHandler) {
}

// dumpPreviousStacktrace is fake for cgo
func dumpPreviousStacktrace() {
}

// previousStacktrace is fake for cgo
func previousStacktrace() string {
}
