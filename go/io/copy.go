// Copyright (c) 2019 The searKing authors. All Rights Reserved.
//
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file in the root of the source
// tree. An additional intellectual property rights grant can be found
// in the file PATENTS.  All contributing project authors may
// be found in the AUTHORS file in the root of the source tree.

package io

import "C"
import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

var ErrNotImplemented = errors.New("not implemented")

// Mode indicates whether to use hardlink or copy content
type Mode int

const (
	// Content creates a new file, and copies the content of the file
	Content Mode = iota
	// Hardlink creates a new hardlink to the existing file
	Hardlink
)

func CopyRegular(srcPath, dstPath string, fileInfo os.FileInfo) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// If the destination file already exists, we shouldn't blow it away
	dstFile, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, fileInfo.Mode())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if err = doCopyFileClone(srcFile, dstFile); err == nil {
		return nil
	}

	if err = doCopyWithFileRange(srcFile, dstFile, fileInfo); err == nil {
		return nil
	}

	return legacyCopy(srcFile, dstFile)
}

func doCopyFileClone(srcFile, dstFile *os.File) error {
	return copyFileClone(dstFile, srcFile)
}

func doCopyWithFileRange(srcFile, dstFile *os.File, fileInfo os.FileInfo) error {
	amountLeftToCopy := fileInfo.Size()

	for amountLeftToCopy > 0 {
		n, err := copyFileRange(int(srcFile.Fd()), nil, int(dstFile.Fd()), nil, int(amountLeftToCopy), 0)
		if err != nil {
			return err
		}

		amountLeftToCopy = amountLeftToCopy - int64(n)
	}

	return nil
}

func legacyCopy(srcFile io.Reader, dstFile io.Writer) error {
	_, err := io.Copy(dstFile, srcFile)
	return err
}

// CopyDir copies or hardlinks the contents of one directory to another,
// properly handling mods, and soft links
func CopyDir(srcDir, dstDir string, copyMode Mode) error {
	return filepath.Walk(srcDir, func(srcPath string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Rebase path
		relPath, err := filepath.Rel(srcDir, srcPath)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dstDir, relPath)
		return copy(srcPath, dstPath, f, copyMode)
	})
}
