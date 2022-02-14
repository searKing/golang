// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build openbsd || (aix && ppc64)
// +build openbsd aix,ppc64

package os

import (
	"syscall"
)

// DiskUsage returns total and free bytes available in a directory, e.g. `/`.
func DiskUsage(path string) (total int64, free int64, avail int64, inodes int64, inodesFree int64, err error) {
	var st syscall.Statfs_t
	if err := syscall.Statfs(path, &st); err != nil {
		return 0, 0, 0, 0, 0, err
	}
	reservedBlocks := int64(st.F_bfree) - int64(st.F_bavail)
	return int64(st.F_bsize) * (int64(st.F_blocks) - reservedBlocks), int64(st.F_bsize) * int64(st.F_bfree), int64(st.F_bsize) * int64(st.F_bavail), int64(st.F_files), int64(st.F_ffree), nil
}
