// Copyright (c) 2019 The searKing authors. All Rights Reserved.
//
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file in the root of the source
// tree. An additional intellectual property rights grant can be found
// in the file PATENTS.  All contributing project authors may
// be found in the AUTHORS file in the root of the source tree.
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
