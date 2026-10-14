// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build unix

package signal

import (
	"os"
	"syscall"
)

// isRuntimeSignal reports whether sig is raised by the Go runtime frequently, such as
// SIGURG for preemption of goroutines, and SIGPROF for profiling.
func isRuntimeSignal(sig os.Signal) bool {
	return sig == syscall.SIGURG || sig == syscall.SIGPROF
}
