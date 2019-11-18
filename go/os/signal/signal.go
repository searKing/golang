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
func init() {
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
	signalAction(sigsToDo...)
}

// DumpSignalTo redirects log to fd, stdout if not set; muted if < 0.
func DumpSignalTo(fd int) {
	dumpSignalTo(fd)
}

// DumpStacktraceTo set dump file path of stacktrace when signal is triggered
// "*.stacktrace.dump" under a temp dir if not set.
func DumpStacktraceTo(name string) {
	dumpStacktraceTo(name)
}

// DumpPreviousStacktrace dumps the previous human readable stacktrace to fd, which is set by SetSignalDumpToFd.
func DumpPreviousStacktrace() {
	dumpPreviousStacktrace()
}

// PreviousStacktrace returns a human readable stacktrace
func PreviousStacktrace() string {
	return previousStacktrace()
}
