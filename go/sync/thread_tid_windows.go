// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync

import "golang.org/x/sys/windows"

// osThreadID returns the ID of the calling OS thread, and whether it is supported.
func osThreadID() (uint64, bool) {
	return uint64(windows.GetCurrentThreadId()), true
}
