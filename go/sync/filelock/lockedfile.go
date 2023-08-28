// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Copied from go/gc/src/cmd/go/internal/lockedfile/lockedfile.go

package filelock

import (
	"fmt"
	"io/fs"
	"os"
	"runtime"
)

type LockedFile[T File] struct {
	File   T
	closed bool
}

// NewLockedFile returns a locked file.
// If flag includes os.O_WRONLY or os.O_RDWR, the file is write-locked;
// otherwise, it is read-locked.
func NewLockedFile(file *os.File) (*LockedFile[*os.File], error) {
	var f LockedFile[*os.File]
	f.File = file

	// Although the operating system will drop locks for open files when the go
	// command exits, we want to hold locks for as little time as possible, and we
	// especially don't want to leave a file locked after we're done with it. Our
	// Close method is what releases the locks, so use a finalizer to report
	// missing Close calls on a best-effort basis.
	runtime.SetFinalizer(&f, func(f *LockedFile[*os.File]) {
		panic(fmt.Sprintf("lockedfile.File %s became unreachable without a call to Close", f.File.Name()))
	})
	return &f, nil
}

// Close unlocks and closes the underlying file.
//
// Close may be called multiple times; all calls after the first will return a
// non-nil error.
func (f *LockedFile[T]) Close() error {
	if f.closed {
		return &fs.PathError{
			Op:   "close",
			Path: f.File.Name(),
			Err:  fs.ErrClosed,
		}
	}
	f.closed = true

	err := closeFile(f.File)
	runtime.SetFinalizer(f, nil)
	return err
}
func (f *LockedFile[T]) Lock() error {
	return Lock(f.File)
}

func (f *LockedFile[T]) TryLock() (bool, error) {
	return TryLock(f.File)
}

func (f *LockedFile[T]) Unlock() error {
	return Unlock(f.File)
}

func (f *LockedFile[T]) RLock() error {
	return RLock(f.File)
}

func (f *LockedFile[T]) RUnlock() error {
	return Unlock(f.File)
}
