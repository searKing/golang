// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package os

import (
	"os"
	"path/filepath"
	"time"
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

// TouchAll creates the named file or dir. If the file already exists,
// it is touched to now. If the file does not exist, it is created with mode 0666 (before umask).
// If the dir does not exist, it is created with mode 0666 (before umask).
func TouchAll(path string, perm os.FileMode) error {
	return createAll(path, perm, true, false)
}

// CreateAll creates or truncates the named file or dir. If the file already exists,
// it is truncated. If the file does not exist, it is created with mode 0666 (before umask).
// If the dir does not exist, it is created with mode 0666 (before umask).
func CreateAll(path string, perm os.FileMode) error {
	return createAll(path, perm, false, true)
}

// CreateAllIfNotExist creates the named file or dir. If the file does not exist, it is created
// with mode 0666 (before umask).
// If the dir does not exist, it is created with mode 0666 (before umask).
// If path is already a directory, CreateAllIfNotExist does nothing
// and returns nil.
func CreateAllIfNotExist(path string, perm os.FileMode) error {
	return createAll(path, perm, false, false)
}

func createAll(path string, perm os.FileMode, touch bool, truncate bool) error {
	dir, file := filepath.Split(path)
	// file or dir exists
	if fi, err := os.Stat(path); err == nil {
		if touch {
			return ChtimesNow(path)
		}
		if fi.IsDir() || !truncate {
			return nil
		}
		// truncates file
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		return f.Close()
	}

	// mkdir -p dir
	if err := os.MkdirAll(dir, perm); err != nil {
		return err
	}

	// create file if needed
	if file == "" {
		return nil
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	_ = f.Close()
	return os.Chmod(path, perm)
}
