// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync

import "golang.org/x/sys/windows"

// callerID returns the ID of the calling OS thread, which identifies the worker goroutine
// of a Thread, as no other goroutine runs on the OS thread it locked.
func callerID() uint64 {
	return uint64(windows.GetCurrentThreadId())
}
