// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build unix || (js && wasm) || wasip1 || windows

package signal

import (
	"os"
	"syscall"
)

const (
	numSig = 65 // max across all systems
)

func Signum(sig os.Signal) int {
	switch sig := sig.(type) {
	case syscall.Signal:
		i := int(sig)
		if i < 0 || i >= numSig {
			return -1
		}
		return i
	default:
		return -1
	}
}

// allSignals returns signals numbered in [0, numSig).
func allSignals() []os.Signal {
	sigs := make([]os.Signal, 0, numSig)
	for n := 0; n < numSig; n++ {
		sigs = append(sigs, syscall.Signal(n))
	}
	return sigs
}
