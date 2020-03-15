// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build cgo

package signal

import "C"
import (
	"os"

	"github.com/searKing/golang/go/os/signal/cgo"
)

// SetSig act as signal.Notify, which invokes the Go signal handler.
// https://godoc.org/os/signal#hdr-Go_programs_that_use_cgo_or_SWIG
func setSig(sigs ...os.Signal) {
	for _, sig := range sigs {
		cgo.SetSig(Signum(sig))
	}
}

// dumpSignalTo redirects log to fd, -1 if not set; muted if < 0.
func dumpSignalTo(fd int) {
	cgo.SetSignalDumpToFd(fd)
}

// dumpStacktraceTo set dump file path of stacktrace when signal is triggered, nop if not set.
func dumpStacktraceTo(name string) {
	cgo.SetBacktraceDumpToFile(name)
}

// dumpPreviousStacktrace dumps human readable stacktrace to fd, which is set by SetSignalDumpToFd.
func dumpPreviousStacktrace() {
	cgo.DumpPreviousStacktrace()
}

// previousStacktrace returns a human readable stacktrace
func previousStacktrace() string {
	return cgo.PreviousStacktrace()
}

// setSigInvokeChain sets a rule to raise signal to {to} and wait until {wait}, done with sleep {sleepInSeconds}s
func setSigInvokeChain(from, to, wait, sleepInSeconds int) {
	cgo.SetSigInvokeChain(from, to, wait, sleepInSeconds)
}
