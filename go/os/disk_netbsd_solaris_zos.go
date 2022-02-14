// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build netbsd || solaris || zos
// +build netbsd solaris zos

package os

import (
	"golang.org/x/sys/unix"
)

// DiskUsage returns total and free bytes available in a directory, e.g. `/`.
func DiskUsage(path string) (total int64, free int64, avail int64, inodes int64, inodesFree int64, err error) {
	st := unix.Statvfs_t{}
	if err := unix.Statvfs(path, &st); err != nil {
		return 0, 0, 0, 0, 0, err
	}
	reservedBlocks := int64(st.Bfree) - int64(st.Bavail)
	// Bsize   uint64 /* file system block size */
	// Frsize  uint64 /* fundamental fs block size */
	// Blocks  uint64 /* number of blocks (unit f_frsize) */
	// Bfree   uint64 /* free blocks in file system */
	// Bavail  uint64 /* free blocks for non-root */
	// Files   uint64 /* total file inodes */
	// Ffree   uint64 /* free file inodes */
	// Favail  uint64 /* free file inodes for to non-root */
	// Fsid    uint64 /* file system id */
	// Flag    uint64 /* bit mask of f_flag values */
	// Namemax uint64 /* maximum filename length */
	return int64(st.Bsize) * (int64(st.Blocks) - reservedBlocks), int64(st.Bsize) * int64(st.Bfree), int64(st.Bsize) * int64(st.Bavail), int64(st.Files), int64(st.Ffree), nil
}
