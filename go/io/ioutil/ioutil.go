// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ioutil

import (
	"io"
	"os"

	os_ "github.com/searKing/golang/go/os"
)

// WriteAll writes data to a file named by filename.
// If the file does not exist, WriteAll creates it with mode 0666 (before umask)
// If the dir does not exist, WriteAll creates it with 0755 (before umask)
// otherwise WriteAll truncates it before writing, without changing permissions.
//
// As of Go 1.16, this function simply calls os.WriteAll.
func WriteAll(filename string, data []byte) error {
	return os_.WriteAll(filename, data)
}

// WriteFileAll is the generalized open call; most users will use WriteAll instead.
// It writes data to a file named by filename.
// If the file does not exist, WriteFileAll creates it with permissions fileperm (before umask)
// If the dir does not exist, WriteFileAll creates it with permissions dirperm (before umask)
// otherwise WriteFileAll truncates it before writing, without changing permissions.
//
// As of Go 1.16, this function simply calls os.WriteFileAll.
func WriteFileAll(filename string, data []byte, dirperm, fileperm os.FileMode) error {
	return os_.WriteFileAll(filename, data, dirperm, fileperm)
}

// WriteAllFrom writes data to a file named by filename from r until EOF or error.
// If the file does not exist, WriteAll creates it with mode 0666 (before umask)
// If the dir does not exist, WriteAll creates it with 0755 (before umask)
// otherwise WriteAll truncates it before writing, without changing permissions.
//
// As of Go 1.16, this function simply calls os.WriteAllFrom.
func WriteAllFrom(filename string, r io.Reader) error {
	return os_.WriteAllFrom(filename, r)
}

// WriteFileAllFrom is the generalized open call; most users will use WriteAllFrom instead.
// It writes data to a file named by filename from r until EOF or error.
// If the file does not exist, WriteFileAllFrom creates it with permissions fileperm (before umask)
// If the dir does not exist, WriteFileAllFrom creates it with permissions dirperm (before umask)
// otherwise WriteFileAllFrom truncates it before writing, without changing permissions.
//
// As of Go 1.16, this function simply calls os.WriteFileAllFrom.
func WriteFileAllFrom(filename string, r io.Reader, dirperm, fileperm os.FileMode) error {
	return os_.WriteFileAllFrom(filename, r, dirperm, fileperm)
}

// AppendAll appends data to a file named by filename.
// If the file does not exist, AppendAll creates it with mode 0666 (before umask)
// If the dir does not exist, AppendAll creates it with 0755 (before umask)
// (before umask); otherwise AppendAll appends it before writing, without changing permissions.
//
// As of Go 1.16, this function simply calls os.AppendAll.
func AppendAll(filename string, data []byte) error {
	return os_.AppendAll(filename, data)
}

// AppendFileAll is the generalized open call; most users will use AppendAll instead.
// It appends data to a file named by filename.
// If the file does not exist, AppendFileAll creates it with permissions fileperm (before umask)
// If the dir does not exist, AppendFileAll creates it with permissions dirperm (before umask)
// otherwise AppendFileAll appends it before writing, without changing permissions.
//
// As of Go 1.16, this function simply calls os.AppendFileAll.
func AppendFileAll(filename string, data []byte, dirperm, fileperm os.FileMode) error {
	return os_.AppendFileAll(filename, data, dirperm, fileperm)
}

// AppendAllFrom appends data to a file named by filename from r until EOF or error.
// If the file does not exist, AppendAllFrom creates it with mode 0666 (before umask)
// If the dir does not exist, AppendAllFrom creates it with 0755 (before umask)
// (before umask); otherwise AppendAllFrom appends it before writing, without changing permissions.
//
// As of Go 1.16, this function simply calls os.AppendAllFrom.
func AppendAllFrom(filename string, r io.Reader) error {
	return os_.AppendAllFrom(filename, r)
}

// AppendFileAllFrom is the generalized open call; most users will use AppendFileFrom instead.
// It appends data to a file named by filename from r until EOF or error.
// If the file does not exist, AppendFileAllFrom creates it with permissions fileperm (before umask)
// If the dir does not exist, AppendFileAllFrom creates it with permissions dirperm (before umask)
// otherwise AppendFileAllFrom appends it before writing, without changing permissions.
//
// As of Go 1.16, this function simply calls os.AppendFileAllFrom.
func AppendFileAllFrom(filename string, r io.Reader, dirperm, fileperm os.FileMode) error {
	return os_.AppendFileAllFrom(filename, r, dirperm, fileperm)
}

// WriteRenameAll writes data to a temp file and rename to the new file named by filename.
// If the file does not exist, WriteRenameAll creates it with mode 0666 (before umask)
// If the dir does not exist, WriteRenameAll creates it with 0755 (before umask)
// otherwise WriteRenameAll truncates it before writing, without changing permissions.
//
// As of Go 1.16, this function simply calls os.WriteRenameAll.
func WriteRenameAll(filename string, data []byte) error {
	return os_.WriteRenameAll(filename, data)
}

// WriteRenameFileAll is the generalized open call; most users will use WriteRenameAll instead.
// WriteRenameFileAll is safer than WriteFileAll as before Write finished, nobody can find the unfinished file.
// It writes data to a temp file and rename to the new file named by filename.
// If the file does not exist, WriteRenameFileAll creates it with permissions fileperm
// If the dir does not exist, WriteRenameFileAll creates it with permissions dirperm
// (before umask); otherwise WriteRenameFileAll truncates it before writing, without changing permissions.
//
// As of Go 1.16, this function simply calls os.WriteRenameFileAllFrom.
func WriteRenameFileAll(filename string, data []byte, dirperm os.FileMode) error {
	return os_.WriteRenameFileAll(filename, data, dirperm)
}

// WriteRenameAllFrom writes data to a temp file from r until EOF or error, and rename to the new file named by filename.
// WriteRenameAllFrom is safer than WriteAllFrom as before Write finished, nobody can find the unfinished file.
// If the file does not exist, WriteRenameAllFrom creates it with mode 0666 (before umask)
// If the dir does not exist, WriteRenameAllFrom creates it with 0755 (before umask)
// otherwise WriteRenameAllFrom truncates it before writing, without changing permissions.
//
// As of Go 1.16, this function simply calls os.WriteRenameAllFrom.
func WriteRenameAllFrom(filename string, r io.Reader) error {
	return os_.WriteRenameAllFrom(filename, r)
}

// WriteRenameFileAllFrom is the generalized open call; most users will use WriteRenameAllFrom instead.
// WriteRenameFileAllFrom is safer than WriteRenameAllFrom as before Write finished, nobody can find the unfinished file.
// It writes data to a temp file and rename to the new file named by filename.
// If the file does not exist, WriteRenameFileAllFrom creates it with permissions fileperm
// If the dir does not exist, WriteRenameFileAllFrom creates it with permissions dirperm
// (before umask); otherwise WriteRenameFileAllFrom truncates it before writing, without changing permissions.
//
// As of Go 1.16, this function simply calls os.WriteRenameFileAllFrom.
func WriteRenameFileAllFrom(filename string, r io.Reader, dirperm os.FileMode) error {
	return os_.WriteRenameFileAllFrom(filename, r, dirperm)
}

// TempAll creates a new temporary file in the directory dir,
// opens the file for reading and writing, and returns the resulting *os.File.
// If the file does not exist, TempAll creates it with mode 0600 (before umask)
// If the dir does not exist, TempAll creates it with 0755 (before umask)
// otherwise TempAll truncates it before writing, without changing permissions.
//
// As of Go 1.16, this function simply calls os.TempAll.
func TempAll(dir, pattern string) (f *os.File, err error) {
	return os_.TempAll(dir, pattern)
}

// TempFileAll is the generalized open call; most users will use TempAll instead.
// If the directory does not exist, it is created with mode dirperm (before umask).
//
// As of Go 1.16, this function simply calls os.TempFileAll.
func TempFileAll(dir, pattern string, dirperm os.FileMode) (f *os.File, err error) {
	return os_.TempFileAll(dir, pattern, dirperm)
}
