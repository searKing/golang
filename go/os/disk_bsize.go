// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build aix || darwin || dragonfly || freebsd

package os

import "syscall"

// statfsBlockSize returns the unit of blocks in st, which is the fundamental block size Bsize.
func statfsBlockSize(st *syscall.Statfs_t) int64 {
	return int64(st.Bsize)
}
