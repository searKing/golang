// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package os

import "syscall"

// statfsBlockSize returns the unit of blocks in st.
// Blocks are counted in Frsize, while Bsize is the optimal transfer block size,
// which may be bigger, such as on NFS. Frsize is 0 before Linux 2.6.
func statfsBlockSize(st *syscall.Statfs_t) int64 {
	if st.Frsize > 0 {
		return int64(st.Frsize)
	}
	return int64(st.Bsize)
}
