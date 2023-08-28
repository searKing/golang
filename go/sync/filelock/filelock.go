// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package filelock provides a platform-independent API for advisory file
// locking. Calls to functions in this package on platforms that do not support
// advisory locks will return errors for which IsNotSupported returns true.
//
// Copied from go/gc/src/cmd/go/internal/lockedfile/internal/filelock/filelock.go

package filelock

import (
	"errors"
	"io/fs"
	"os"
)

// A File provides the minimal set of methods required to lock an open file.
// File implementations must be usable as map keys.
// The usual implementation is *os.File.
type File interface {
	// Name returns the name of the file.
	Name() string

	// Fd returns a valid file descriptor.
	// (If the File is an *os.File, it must not be closed.)
	Fd() uintptr

	// Stat returns the FileInfo structure describing file.
	Stat() (fs.FileInfo, error)
}

// Lock places an advisory write lock on the file, blocking until it can be
// locked.
//
// If Lock returns nil, no other process will be able to place a read or write
// lock on the file until this process exits, closes f, or calls Unlock on it.
//
// If f's descriptor is already read- or write-locked, the behavior of Lock is
// unspecified.
//
// Closing the file may or may not release the lock promptly. Callers should
// ensure that Unlock is always called when Lock succeeds.
func Lock(f File) error {
	return lock(f, writeLock, false)
}

// TryLock tries to lock rw for writing and reports whether it succeeded.
//
// Note that while correct uses of TryLock do exist, they are rare,
// and use of TryLock is often a sign of a deeper problem
// in a particular use of mutexes.
func TryLock(f File) (bool, error) {
	err := lock(f, writeLock, true)
	if err != nil {
		if os.IsTimeout(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// RLock places an advisory read lock on the file, blocking until it can be locked.
//
// If RLock returns nil, no other process will be able to place a write lock on
// the file until this process exits, closes f, or calls Unlock on it.
//
// If f is already read- or write-locked, the behavior of RLock is unspecified.
//
// Closing the file may or may not release the lock promptly. Callers should
// ensure that Unlock is always called if RLock succeeds.
func RLock(f File) error {
	return lock(f, readLock, false)
}

// TryRLock tries to lock rw for reading and reports whether it succeeded.
//
// Note that while correct uses of TryRLock do exist, they are rare,
// and use of TryRLock is often a sign of a deeper problem
// in a particular use of mutexes.
func TryRLock(f File) (bool, error) {
	err := lock(f, readLock, true)
	if err != nil {
		if os.IsTimeout(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Unlock removes an advisory lock placed on f by this process.
//
// The caller must not attempt to unlock a file that is not locked.
func Unlock(f File) error {
	return unlock(f)
}

// String returns the name of the function corresponding to lt
// (Lock, RLock, or Unlock).
func (lt lockType) String() string {
	switch lt {
	case readLock:
		return "RLock"
	case writeLock:
		return "Lock"
	default:
		return "Unlock"
	}
}

// IsNotSupported returns a boolean indicating whether the error is known to
// report that a function is not supported (possibly for a specific input).
// It is satisfied by errors.ErrUnsupported as well as some syscall errors.
func IsNotSupported(err error) bool {
	return errors.Is(err, errors.ErrUnsupported)
}
