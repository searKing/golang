// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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

// dumpPreviousStacktrace is fake for cgo
func dumpPreviousStacktrace() {
}

// previousStacktrace is fake for cgo
func previousStacktrace() string {
}

// setSigInvokeChain is fake for cgo
func setSigInvokeChain(from, to, wait, sleepInSeconds int) {
}
