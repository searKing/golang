// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// +build linux,cgo

package io

/*
   #include <linux/fs.h>

   #ifndef FICLONE
   #define FICLONE		_IOW(0x94, 9, int)
   #endif
*/
import "C"
import (
	"os"

	"golang.org/x/sys/unix"
)

func copyFileClone(srcFile, dstFile *os.File) error {
	_, _, err := unix.Syscall(unix.SYS_IOCTL, dstFile.Fd(), C.FICLONE, srcFile.Fd())
	return err
}
