// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ioutil

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	os_ "github.com/searKing/golang/go/os"
)

// WriteAll writes data to a file named by filename.
// If the file does not exist, WriteAll creates it with mode 0666 (before umask)
// If the dir does not exist, WriteAll creates it with 0755 (before umask)
// otherwise WriteAll truncates it before writing, without changing permissions.
func WriteAll(filename string, data []byte) error {
	return WriteFileAll(filename, data, 0755, 0666)
}

// WriteFileAll is the generalized open call; most users will use WriteAll instead.
// It writes data to a file named by filename.
// If the file does not exist, WriteFileAll creates it with permissions fileperm (before umask)
// If the dir does not exist, WriteFileAll creates it with permissions dirperm (before umask)
// otherwise WriteFileAll truncates it before writing, without changing permissions.
func WriteFileAll(filename string, data []byte, dirperm, fileperm os.FileMode) error {
	return WriteFileAllFrom(filename, bytes.NewReader(data), dirperm, fileperm)
}

// WriteAllFrom writes data to a file named by filename from r until EOF or error.
// If the file does not exist, WriteAll creates it with mode 0666 (before umask)
// If the dir does not exist, WriteAll creates it with 0755 (before umask)
// otherwise WriteAll truncates it before writing, without changing permissions.
func WriteAllFrom(filename string, r io.Reader) error {
	return WriteFileAllFrom(filename, r, 0755, 0666)
}

// WriteFileAllFrom is the generalized open call; most users will use WriteAllFrom instead.
// It writes data to a file named by filename from r until EOF or error.
// If the file does not exist, WriteFileAllFrom creates it with permissions fileperm (before umask)
// If the dir does not exist, WriteFileAllFrom creates it with permissions dirperm (before umask)
// otherwise WriteFileAllFrom truncates it before writing, without changing permissions.
func WriteFileAllFrom(filename string, r io.Reader, dirperm, fileperm os.FileMode) error {
	f, err := os_.OpenFileAll(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, dirperm, fileperm)
	if err != nil {
		return err
	}
	_, err = f.ReadFrom(r)
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}

// AppendAll appends data to a file named by filename.
// If the file does not exist, AppendAll creates it with mode 0666 (before umask)
// If the dir does not exist, AppendAll creates it with 0755 (before umask)
// (before umask); otherwise AppendAll appends it before writing, without changing permissions.
func AppendAll(filename string, data []byte) error {
	return AppendFileAll(filename, data, 0755, 0666)
}

// AppendFileAll is the generalized open call; most users will use AppendAll instead.
// It appends data to a file named by filename.
// If the file does not exist, WriteFileAll creates it with permissions fileperm (before umask)
// If the dir does not exist, WriteFileAll creates it with permissions dirperm (before umask)
// otherwise WriteFileAll appends it before writing, without changing permissions.
func AppendFileAll(filename string, data []byte, dirperm, fileperm os.FileMode) error {
	return AppendFileAllFrom(filename, bytes.NewReader(data), dirperm, fileperm)
}

// AppendAllFrom appends data to a file named by filename from r until EOF or error.
// If the file does not exist, AppendAllFrom creates it with mode 0666 (before umask)
// If the dir does not exist, AppendAllFrom creates it with 0755 (before umask)
// (before umask); otherwise AppendAllFrom appends it before writing, without changing permissions.
func AppendAllFrom(filename string, r io.Reader) error {
	return AppendFileAllFrom(filename, r, 0755, 0666)
}

// AppendFileAllFrom is the generalized open call; most users will use AppendFileFrom instead.
// It appends data to a file named by filename from r until EOF or error.
// If the file does not exist, AppendFileAllFrom creates it with permissions fileperm (before umask)
// If the dir does not exist, AppendFileAllFrom creates it with permissions dirperm (before umask)
// otherwise AppendFileAllFrom appends it before writing, without changing permissions.
func AppendFileAllFrom(filename string, r io.Reader, dirperm, fileperm os.FileMode) error {
	f, err := os_.OpenFileAll(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, dirperm, fileperm)
	if err != nil {
		return err
	}
	_, err = f.ReadFrom(r)
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}

// WriteRenameAll writes data to a temp file and rename to the new file named by filename.
// If the file does not exist, WriteRenameAll creates it with mode 0666 (before umask)
// If the dir does not exist, WriteRenameAll creates it with 0755 (before umask)
// otherwise WriteRenameAll truncates it before writing, without changing permissions.
func WriteRenameAll(filename string, data []byte) error {
	return WriteRenameFileAll(filename, data, 0755)
}

// WriteRenameFileAll is the generalized open call; most users will use WriteRenameAll instead.
// WriteRenameFileAll is safer than WriteFileAll as before Write finished, nobody can find the unfinished file.
// It writes data to a temp file and rename to the new file named by filename.
// If the file does not exist, WriteRenameFileAll creates it with permissions fileperm
// If the dir does not exist, WriteRenameFileAll creates it with permissions dirperm
// (before umask); otherwise WriteRenameFileAll truncates it before writing, without changing permissions.
func WriteRenameFileAll(filename string, data []byte, dirperm os.FileMode) error {
	return WriteRenameFileAllFrom(filename, bytes.NewReader(data), 0755)
}

// WriteRenameAllFrom writes data to a temp file from r until EOF or error, and rename to the new file named by filename.
// WriteRenameAllFrom is safer than WriteAllFrom as before Write finished, nobody can find the unfinished file.
// If the file does not exist, WriteRenameAllFrom creates it with mode 0666 (before umask)
// If the dir does not exist, WriteRenameAllFrom creates it with 0755 (before umask)
// otherwise WriteRenameAllFrom truncates it before writing, without changing permissions.
func WriteRenameAllFrom(filename string, r io.Reader) error {
	return WriteRenameFileAllFrom(filename, r, 0755)
}

// WriteRenameFileAllFrom is the generalized open call; most users will use WriteRenameAllFrom instead.
// WriteRenameFileAllFrom is safer than WriteRenameAllFrom as before Write finished, nobody can find the unfinished file.
// It writes data to a temp file and rename to the new file named by filename.
// If the file does not exist, WriteRenameFileAllFrom creates it with permissions fileperm
// If the dir does not exist, WriteRenameFileAllFrom creates it with permissions dirperm
// (before umask); otherwise WriteRenameFileAllFrom truncates it before writing, without changing permissions.
func WriteRenameFileAllFrom(filename string, r io.Reader, dirperm os.FileMode) error {
	tempDir := filepath.Dir(filename)
	if tempDir != "" {
		// mkdir -p dir
		if err := os.MkdirAll(tempDir, dirperm); err != nil {
			return err
		}
	}

	tempFile, err := ioutil.TempFile(tempDir, "")
	if err != nil {
		return err
	}
	defer tempFile.Close()

	tempFilePath := tempFile.Name()
	defer os.Remove(tempFilePath)
	_, err = tempFile.ReadFrom(r)
	if err != nil {
		return err
	}
	return os_.RenameFileAll(tempFilePath, filename, dirperm)
}

// TempAll creates a new temporary file in the directory dir,
// opens the file for reading and writing, and returns the resulting *os.File.
// If the file does not exist, TempAll creates it with mode 0600 (before umask)
// If the dir does not exist, TempAll creates it with 0755 (before umask)
// otherwise TempAll truncates it before writing, without changing permissions.
func TempAll(dir, pattern string) (f *os.File, err error) {
	return TempFileAll(dir, pattern, 0755)
}

// TempFileAll is the generalized open call; most users will use TempAll instead.
// If the directory does not exist, it is created with mode dirperm (before umask).
func TempFileAll(dir, pattern string, dirperm os.FileMode) (f *os.File, err error) {
	if dir != "" {
		// mkdir -p dir
		if err := os.MkdirAll(dir, dirperm); err != nil {
			return nil, err
		}
	}
	return ioutil.TempFile(dir, pattern)
}
