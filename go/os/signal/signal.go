// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package signal

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

// enhance signal.Notify with stacktrace of cgo.
// redirects signal log to stderr
func init() {
	DumpSignalTo(int(syscall.Stderr))
	// FIXME https://github.com/golang/go/issues/35814
	//RegisterOnSignal(OnSignalHandlerFunc(func(signum os.Signal) {}))

	var dumpfile string
	if f, err := os.CreateTemp("", "*.stacktrace.dump"); err == nil {
		dumpfile = f.Name()
		_ = f.Close()
		// only the name is reserved, the file is created when stacktrace is dumped,
		// so that PreviousStacktrace returns empty before that.
		_ = os.Remove(dumpfile)
	} else {
		dumpfile = filepath.Join(os.TempDir(), fmt.Sprintf("stacktrace.%d.dump", time.Now().UnixNano()))
	}

	DumpStacktraceTo(dumpfile)
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
// Signals raised by the Go runtime frequently, such as SIGURG and SIGPROF, are relayed to c only,
// without stacktraces dumped.
func Notify(c chan<- os.Signal, sigs ...os.Signal) {
	if len(sigs) == 0 {
		sigs = allSignals()
	}
	signal.Notify(c, sigs...)
	for _, sig := range sigs {
		if !isRuntimeSignal(sig) {
			setSig(sig)
		}
	}
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
// It waits forever if {wait} is never received, see SetSigInvokeChainContext.
func SetSigInvokeChain(to os.Signal, wait os.Signal, sleepInSeconds int, froms ...os.Signal) {
	for _, from := range froms {
		setSigInvokeChain(Signum(from), Signum(to), Signum(wait), sleepInSeconds)
	}
}

// SetSigInvokeChainContext acts as SetSigInvokeChain, with sleep rounded up to seconds,
// but the chain is removed when ctx is done, and its wait in progress returns.
// The wait returns within about 10ms after ctx is done, if goroutines can still run.
func SetSigInvokeChainContext(ctx context.Context, to os.Signal, wait os.Signal, sleep time.Duration, froms ...os.Signal) {
	if ctx.Err() != nil {
		return
	}
	for _, from := range froms {
		cancel := setSigInvokeChainCancelable(Signum(from), Signum(to), Signum(wait), ceilSeconds(sleep))
		context.AfterFunc(ctx, cancel)
	}
}

// ceilSeconds returns d in seconds rounded up, 0 if d <= 0.
func ceilSeconds(d time.Duration) int {
	if d <= 0 {
		return 0
	}
	return int((d + time.Second - 1) / time.Second)
}
