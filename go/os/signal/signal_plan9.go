// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package signal

import (
	"os"
	"syscall"
)

const numSig = 256

func Signum(sig os.Signal) int {
	switch sig := sig.(type) {
	case syscall.Note:
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
