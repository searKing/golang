// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package os

import (
	"io"
	"os"
	"path/filepath"
	"time"
)

const (
	DefaultPermissionFile      os.FileMode = 0644
	DefaultPermissionDirectory os.FileMode = 0755
)

func GetAbsBinDir() (dir string, err error) {
	return filepath.Abs(filepath.Dir(os.Args[0]))
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Chtimes changes the access and modification times of the named
// file with Now, similar to the Unix utime() or utimes() functions.
//
// The underlying filesystem may truncate or round the values to a
// less precise time unit.
// If there is an error, it will be of type *PathError.
func ChtimesNow(name string) error {
	now := time.Now()
	return os.Chtimes(name, now, now)
}

// MakeAll creates a directory named path,
// along with any necessary parents, and returns nil,
// or else returns an error.
// If the dir does not exist, it is created with mode 0755 (before umask).
func MakeAll(name string) error {
	return os.MkdirAll(name, 0755)
}

// MakeAll creates a directory named path and returns nil,
// or else returns an error.
// If the dir does not exist, it is created with mode 0755 (before umask).
func Make(name string) error {
	return os.Mkdir(name, 0755)
}

// TouchAll creates the named file or dir. If the file already exists, it is touched to now.
// If the file does not exist, it is created with mode 0666 (before umask).
// If the dir does not exist, it is created with mode 0755 (before umask).
func TouchAll(path string) (*os.File, error) {
	f, err := OpenFileAll(path, os.O_WRONLY|os.O_CREATE, 0755, 0666)
	if err != nil {
		return nil, err
	}
	if err := ChtimesNow(path); err != nil {
		defer f.Close()
		return nil, err
	}
	return f, nil
}

// CreateAll creates or truncates the named file or dir. If the file already exists,
// it is truncated. If the file does not exist, it is created with mode 0666 (before umask).
// If the dir does not exist, it is created with mode 0755 (before umask).
func CreateAll(path string) (*os.File, error) {
	return OpenFileAll(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755, 0666)
}

// CreateAllIfNotExist creates the named file or dir. If the file does not exist, it is created
// with mode 0666 (before umask).
// If the dir does not exist, it is created with mode 0755 (before umask).
// If path is already a directory, CreateAllIfNotExist does nothing and returns nil.
func CreateAllIfNotExist(path string) (*os.File, error) {
	return OpenFileAll(path, os.O_RDWR|os.O_CREATE, 0755, 0666)
}

// OpenAll opens the named file or dir for reading. If successful, methods on
// the returned file or dir can be used for reading; the associated file
// descriptor has mode O_RDONLY.
// If there is an error, it will be of type *PathError.
func OpenAll(path string) (*os.File, error) {
	return OpenFileAll(path, os.O_RDONLY, 0, 0)
}

// OpenFileAll is the generalized open call; most users will use OpenAll
// or CreateAll instead. It opens the named file or directory with specified flag
// (O_RDONLY etc.).
// If the file does not exist, and the O_CREATE flag is passed, it is created with mode fileperm (before umask).
// If the directory does not exist,, it is created with mode dirperm (before umask).
// If successful, methods on the returned File can be used for I/O.
// If there is an error, it will be of type *PathError.
func OpenFileAll(path string, flag int, dirperm, fileperm os.FileMode) (*os.File, error) {
	dir, file := filepath.Split(path)
	// file or dir exists
	if _, err := os.Stat(path); err == nil {
		return os.OpenFile(path, flag, 0)
	}

	// mkdir -p dir
	if err := os.MkdirAll(dir, dirperm); err != nil {
		return nil, err
	}

	// create file if needed
	if file == "" {
		return nil, nil
	}

	return os.OpenFile(path, flag, fileperm)
}

// CopyAll creates or truncates the dst file or dir, filled with content from src file.
// If the dst file already exists, it is truncated.
// If the dst file does not exist, it is created with mode 0666 (before umask).
// If the dst dir does not exist, it is created with mode 0755 (before umask).
func CopyAll(dst string, src string) error {
	return CopyFileAll(dst, src, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755, 0666)
}

// AppendAll creates or appends the dst file or dir, filled with content from src file.
// If the dst file already exists, it is truncated.
// If the dst file does not exist, it is created with mode 0666 (before umask).
// If the dst dir does not exist, it is created with mode 0755 (before umask).
func AppendAll(dst string, src string) error {
	return CopyFileAll(dst, src, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755, 0666)
}

// CopyFileAll is the generalized open call; most users will use CopyAll
// or AppendAll instead. It opens the named file or directory with specified flag
// (O_RDONLY etc.).
// If the dst file does not exist, and the O_CREATE flag is passed, it is created with mode fileperm (before umask).
// If the dst directory does not exist,, it is created with mode dirperm (before umask).
// If successful, methods on the returned File can be used for I/O.
// If there is an error, it will be of type *PathError.
func CopyFileAll(dst string, src string, flag int, dirperm, fileperm os.FileMode) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := OpenFileAll(dst, flag, dirperm, fileperm)
	if err != nil {
		return err
	}

	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// Copy creates or truncates the dst file or dir, filled with content from src file.
// If the dst file already exists, it is truncated.
// If the dst file does not exist, it is created with mode 0666 (before umask).
// If the dst dir does not exist, it is created with mode 0755 (before umask).
// parent dirs will not be created, otherwise, use CopyAll instead.
func Copy(dst string, src string) error {
	return CopyFile(dst, src, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
}

// Append creates or appends the dst file or dir, filled with content from src file.
// If the dst file already exists, it is truncated.
// If the dst file does not exist, it is created with mode 0666 (before umask).
// If the dst dir does not exist, it is created with mode 0755 (before umask).
// parent dirs will not be created, otherwise, use AppendAll instead.
func Append(dst string, src string) error {
	return CopyFile(dst, src, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
}

// CopyFile is the generalized open call; most users will use Copy
// or Append instead. It opens the named file or directory with specified flag
// (O_RDONLY etc.).
// CopyFile copies from src to dst.
// parent dirs will not be created, otherwise, use CopyFileAll instead.
func CopyFile(dst string, src string, flag int, perm os.FileMode) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dst, flag, perm)
	if err != nil {
		return err
	}

	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// SameFile reports whether fi1 and fi2 describe the same file.
// Overload os.SameFile by file path
func SameFile(fi1, fi2 string) bool {
	stat1, err := os.Stat(fi1)
	if err != nil {
		return false
	}

	stat2, err := os.Stat(fi2)
	if err != nil {
		return false
	}
	return os.SameFile(stat1, stat2)
}
