// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filelock

import (
	"io"
	"io/fs"
	"os"

	os_ "github.com/searKing/golang/go/os"
)

// OpenFile is like os.OpenFile, but returns a locked file.
// If flag includes os.O_WRONLY or os.O_RDWR, the file is write-locked;
// otherwise, it is read-locked.
func OpenFile(name string, flag int, perm fs.FileMode) (*LockedFile[*os.File], error) {
	file, err := openFile(name, flag, perm)
	if err != nil {
		return nil, err
	}
	return NewLockedFile(file)
}

// Open is like os.Open, but returns a read-locked file.
func Open(name string) (*LockedFile[*os.File], error) {
	return OpenFile(name, os.O_RDONLY, 0)
}

// Create is like os.Create, but returns a write-locked file.
func Create(name string) (*LockedFile[*os.File], error) {
	return OpenFile(name, os_.DefaultFlagCreateTruncate, os_.DefaultPermissionFile)
}

// Edit creates the named file with mode 0666 (before umask),
// but does not truncate existing contents.
//
// If Edit succeeds, methods on the returned File can be used for I/O.
// The associated file descriptor has mode O_RDWR and the file is write-locked.
func Edit(name string) (*LockedFile[*os.File], error) {
	return OpenFile(name, os_.DefaultFlagCreateIfNotExist, os_.DefaultPermissionFile)
}

// Read opens the named file with a read-lock and returns its contents.
func Read(name string) ([]byte, error) {
	f, err := Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return io.ReadAll(f.File)
}

// Write opens the named file (creating it with the given permissions if needed),
// then write-locks it and overwrites it with the given content.
func Write(name string, content io.Reader, perm fs.FileMode) (err error) {
	f, err := OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}

	_, err = io.Copy(f.File, content)
	if closeErr := f.Close(); err == nil {
		err = closeErr
	}
	return err
}

// Transform invokes t with the result of reading the named file, with its lock
// still held.
//
// If t returns a nil error, Transform then writes the returned contents back to
// the file, making a best effort to preserve existing contents on error.
//
// t must not modify the slice passed to it.
func Transform(name string, t func([]byte) ([]byte, error)) (err error) {
	f, err := Edit(name)
	if err != nil {
		return err
	}
	defer f.Close()

	old, err := io.ReadAll(f.File)
	if err != nil {
		return err
	}

	new, err := t(old)
	if err != nil {
		return err
	}

	if len(new) > len(old) {
		// The overall file size is increasing, so write the tail first: if we're
		// about to run out of space on the disk, we would rather detect that
		// failure before we have overwritten the original contents.
		if _, err := f.File.WriteAt(new[len(old):], int64(len(old))); err != nil {
			// Make a best effort to remove the incomplete tail.
			f.File.Truncate(int64(len(old)))
			return err
		}
	}

	// We're about to overwrite the old contents. In case of failure, make a best
	// effort to roll back before we close the file.
	defer func() {
		if err != nil {
			if _, err := f.File.WriteAt(old, 0); err == nil {
				f.File.Truncate(int64(len(old)))
			}
		}
	}()

	if len(new) >= len(old) {
		if _, err := f.File.WriteAt(new[:len(old)], 0); err != nil {
			return err
		}
	} else {
		if _, err := f.File.WriteAt(new, 0); err != nil {
			return err
		}
		// The overall file size is decreasing, so shrink the file to its final size
		// after writing. We do this after writing (instead of before) so that if
		// the write fails, enough filesystem space will likely still be reserved
		// to contain the previous contents.
		if err := f.File.Truncate(int64(len(new))); err != nil {
			return err
		}
	}

	return nil
}
