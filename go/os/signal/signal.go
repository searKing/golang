package signal

import "C"
import (
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"

	"github.com/google/uuid"
)

// enhance signal.Notify with stacktrace of cgo.
// redirects signal log to stdout
func init() {
	DumpSignalTo(syscall.Stdout)
	// FIXME https://github.com/golang/go/issues/35814
	//RegisterOnSignal(OnSignalHandlerFunc(func(signum os.Signal) {}))

	var dumpfile string
	if f, err := ioutil.TempFile("", "*.stacktrace.dump"); err == nil {
		dumpfile = f.Name()
	} else {
		dumpfile = filepath.Join(os.TempDir(), uuid.New().String()+".stacktrace.dump")
	}

	DumpStacktraceTo(dumpfile)
	defer os.Remove(dumpfile)

	var sigsToDo []os.Signal
	for n := 0; n < numSig; n++ {
		sigsToDo = append(sigsToDo, syscall.Signal(n))
	}
	if len(sigsToDo) == 0 {
		return
	}
	setSig(sigsToDo...)
}

type OnSignalHandler interface {
	OnSignal(signum os.Signal)
}

type OnSignalHandlerFunc func(signum os.Signal)

func (f OnSignalHandlerFunc) OnSignal(signum os.Signal) {
	f(signum)
}

// DumpSignalTo redirects log to fd, -1 if not set; muted if < 0.
func DumpSignalTo(fd int) {
	dumpSignalTo(fd)
}

// DumpStacktraceTo set dump file path of stacktrace when signal is triggered
// "*.stacktrace.dump" under a temp dir if not set.
func DumpStacktraceTo(name string) {
	dumpStacktraceTo(name)
}

func RegisterOnSignal(onSignal OnSignalHandler) {
	registerOnSignal(onSignal)
}

// DumpPreviousStacktrace dumps the previous human readable stacktrace to fd, which is set by SetSignalDumpToFd.
func DumpPreviousStacktrace() {
	dumpPreviousStacktrace()
}

// PreviousStacktrace returns a human readable stacktrace
func PreviousStacktrace() string {
	return previousStacktrace()
}
