// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package signal

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

// enhance signal.Notify with stacktrace of cgo.
// redirects signal log to stdout
func init() {
	DumpSignalTo(int(syscall.Stdout))
	// FIXME https://github.com/golang/go/issues/35814
	//RegisterOnSignal(OnSignalHandlerFunc(func(signum os.Signal) {}))

	var dumpfile string
	if f, err := os.CreateTemp("", "*.stacktrace.dump"); err == nil {
		dumpfile = f.Name()
	} else {
		dumpfile = filepath.Join(os.TempDir(), fmt.Sprintf("stacktrace.%d.dump", time.Now().UnixNano()))
	}

	DumpStacktraceTo(dumpfile)
	defer os.Remove(dumpfile)
}

type OnSignalHandler interface {
	OnSignal(signum os.Signal)
}

type OnSignalHandlerFunc func(signum os.Signal)

func (f OnSignalHandlerFunc) OnSignal(signum os.Signal) {
	f(signum)
}

// Notify act as signal.Notify, which invokes the Go signal handler.
// https://godoc.org/os/signal#hdr-Go_programs_that_use_cgo_or_SWIG
// Notify must be called again when one sig is called on windows system
// as windows is based on signal(), which will reset sig's handler to SIG_DFL before sig's handler is called
// While unix-like os will remain sig's handler always.
func Notify(c chan<- os.Signal, sigs ...os.Signal) {
	if len(sigs) == 0 {
		for n := 0; n < numSig; n++ {
			sigs = append(sigs, syscall.Signal(n))
		}
	}
	signal.Notify(c, sigs...)
	setSig(sigs...)
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

// DumpPreviousStacktrace dumps the previous human readable stacktrace to fd, which is set by SetSignalDumpToFd.
func DumpPreviousStacktrace() {
	dumpPreviousStacktrace()
}

// PreviousStacktrace returns a human readable stacktrace
func PreviousStacktrace() string {
	return previousStacktrace()
}

// SetSigInvokeChain sets a rule to raise signal to {to} and wait until {wait}, done with sleep {sleepInSeconds}s
func SetSigInvokeChain(to os.Signal, wait os.Signal, sleepInSeconds int, froms ...os.Signal) {
	for _, from := range froms {
		setSigInvokeChain(Signum(from), Signum(to), Signum(wait), sleepInSeconds)
	}
}
