// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !linux && !windows

package sync

// osThreadID returns the ID of the calling OS thread, and whether it is supported.
func osThreadID() (uint64, bool) {
	return 0, false
}
