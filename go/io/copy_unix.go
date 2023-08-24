// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !windows && !linux
// +build !windows,!linux

package io

import (
	"fmt"
	"os"
	"syscall"

	"golang.org/x/sys/unix"
)

func copyPath(srcPath, dstPath string, f os.FileInfo, copyMode Mode) error {
	stat, ok := f.Sys().(*syscall.Stat_t)
	if !ok {
		return fmt.Errorf("unable to get raw syscall.Stat_t data for %s", srcPath)
	}

	isHardlink := false

	switch mode := f.Mode(); {
	case mode.IsRegular():
		// the type is 32bit on mips
		if copyMode == Hardlink {
			isHardlink = true
			if err := os.Link(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := CopyRegular(srcPath, dstPath, f); err != nil {
				return err
			}
		}

	case mode.IsDir():
		if err := os.Mkdir(dstPath, f.Mode()); err != nil && !os.IsExist(err) {
			return err
		}

	case mode&os.ModeSymlink != 0:
		link, err := os.Readlink(srcPath)
		if err != nil {
			return err
		}

		if err := os.Symlink(link, dstPath); err != nil {
			return err
		}

	case mode&os.ModeNamedPipe != 0:
		fallthrough
	case mode&os.ModeSocket != 0:
		if err := unix.Mkfifo(dstPath, uint32(stat.Mode)); err != nil {
			return err
		}

	case mode&os.ModeDevice != 0:
		if err := unix.Mknod(dstPath, uint32(stat.Mode), int(stat.Rdev)); err != nil {
			return err
		}

	default:
		return fmt.Errorf("unknown file type (%d / %s) for %s", f.Mode(), f.Mode().String(), srcPath)
	}

	// Everything below is copying metadata from src to dst. All this metadata
	// already shares an inode for hardlinks.
	if isHardlink {
		return nil
	}

	if err := os.Lchown(dstPath, int(stat.Uid), int(stat.Gid)); err != nil {
		return err
	}

	isSymlink := f.Mode()&os.ModeSymlink != 0

	// There is no LChmod, so ignore mode for symlink. Also, this
	// must happen after chown, as that can modify the file mode
	if !isSymlink {
		if err := os.Chmod(dstPath, f.Mode()); err != nil {
			return err
		}
	}

	return nil
}
