// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package signal

import (
	"os"
	"sync"
	"syscall"
)

const numSig = 256

var (
	sigtabMu sync.Mutex
	sigtab   = make(map[os.Signal]int)
)

func Signum(sig os.Signal) int {
	switch sig := sig.(type) {
	case syscall.Note:
		sigtabMu.Lock()
		defer sigtabMu.Unlock()
		n, ok := sigtab[sig]
		if !ok {
			n = len(sigtab) + 1
			if n > numSig {
				return -1
			}
			sigtab[sig] = n
		}
		return n
	default:
		return -1
	}
}

// allSignals returns nil, as notes can not be enumerated,
// and signal.Notify relays all notes if no signals are provided.
func allSignals() []os.Signal {
	return nil
}
