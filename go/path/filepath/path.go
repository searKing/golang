// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filepath

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	// PrivateFileMode grants owner to read/write a file.
	PrivateFileMode = 0600
	// PrivateDirMode grants owner to make/remove files inside the directory.
	PrivateDirMode = 0700
)

// Pathify Expand, Abs and Clean the path
func Pathify(path string) string {
	p := os.ExpandEnv(path)

	if filepath.IsAbs(p) {
		return filepath.Clean(p)
	}

	p, err := filepath.Abs(p)
	if err == nil {
		return filepath.Clean(p)
	}
	return ""
}

// ToDir returns the dir format ends with OS-specific path separator.
func ToDir(path string) string {
	sep := string(filepath.Separator)
	if strings.HasSuffix(path, sep) {
		return path
	}
	return path + sep
}

// Exist returns a boolean indicating whether the file is known to
// report that a file or directory does exist.
func Exist(name string) bool {
	_, err := os.Stat(name)
	if err == nil {
		return true
	}
	return !os.IsNotExist(err)
}

// TouchAll creates a file or a directory only if it does not already exist.
func TouchAll(path string, perm os.FileMode) error {
	dir, file := filepath.Split(path)

	fileInfo, err := os.Stat(path)
	if err == nil {
		// file type mismatched
		if fileInfo.IsDir() && file != "" {
			return os.ErrExist
		}
	} else {
		// other errors
		if !os.IsNotExist(err) {
			return err
		}
		// Not Exist
	}

	// touch dir
	if dir != "" && !Exist(dir) {
		if err := os.MkdirAll(dir, perm); err != nil {
			return err
		}
	}

	// return if path is a dir
	if file == "" {
		return nil
	}
	// touch file
	f, err := os.OpenFile(path, os.O_CREATE, perm)
	if err != nil {
		return err
	}
	_ = f.Close()
	return nil
}
