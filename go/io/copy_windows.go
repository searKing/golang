// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package io

import (
	"fmt"
	"os"
)

func copyFileClone(srcFile, dstFile *os.File) error {
	return ErrNotImplemented
}

func copyFileRange(rfd int, roff *int64, wfd int, woff *int64, len int, flags int) (n int, err error) {
	return 0, ErrNotImplemented
}

func copyPath(srcPath, dstPath string, f os.FileInfo, copyMode Mode) error {
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

	default:
		return fmt.Errorf("unknown file type (%d / %s) for %s", f.Mode(), f.Mode().String(), srcPath)
	}

	// Everything below is copying metadata from src to dst. All this metadata
	// already shares an inode for hardlinks.
	if isHardlink {
		return nil
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
