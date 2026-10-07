// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package os

import (
	"errors"
	"os"
	"sort"
	"time"

	filepath_ "github.com/searKing/golang/go/path/filepath"
)

// DiskQuota limits files matched by a pattern. A limit takes effect only if it is bigger than 0.
type DiskQuota struct {
	MaxAge             time.Duration // max age of files, by ModTime
	MaxCount           int           // max count of files
	MaxUsedProportion  float32       // max used proportion of disk bytes, in (0, 1]
	MaxIUsedProportion float32       // max used proportion of disk inodes, in (0, 1]
}

// NoLimit reports whether no limit is set.
func (q DiskQuota) NoLimit() bool {
	return q.MaxAge <= 0 && q.MaxCount <= 0 && q.NoUsageLimit()
}

// NoUsageLimit reports whether neither MaxUsedProportion nor MaxIUsedProportion is set.
func (q DiskQuota) NoUsageLimit() bool {
	return q.MaxUsedProportion <= 0 && q.MaxIUsedProportion <= 0
}

// ExceedCount reports whether n files exceed MaxCount.
func (q DiskQuota) ExceedCount(n int) bool {
	return q.MaxCount > 0 && n > q.MaxCount
}

// ExceedBytes reports whether used disk bytes exceed MaxUsedProportion of total.
func (q DiskQuota) ExceedBytes(avail, total int64) bool {
	return q.MaxUsedProportion > 0 && float32(total-avail) > q.MaxUsedProportion*float32(total)
}

// ExceedInodes reports whether used disk inodes exceed MaxIUsedProportion of inodes.
func (q DiskQuota) ExceedInodes(inodes, inodesFree int64) bool {
	return q.MaxIUsedProportion > 0 && float32(inodes-inodesFree) > q.MaxIUsedProportion*float32(inodes)
}

// UnlinkOldestFiles unlinks files matching pattern which exceed quora.
// See [UnlinkOldestFilesFunc].
func UnlinkOldestFiles(pattern string, quora DiskQuota) error {
	return UnlinkOldestFilesFunc(pattern, quora, func(name string) bool { return true })
}

// UnlinkOldestFilesFunc unlinks files matching pattern which exceed quora,
// skipping files for which f(name) returns false.
//
// Files are unlinked in the following order:
//  1. files older than MaxAge, symbolic links excluded;
//  2. the oldest files by ModTime, until no more than MaxCount files left;
//  3. the oldest files by ModTime, until disk usage is within MaxUsedProportion and MaxIUsedProportion.
func UnlinkOldestFilesFunc(pattern string, quora DiskQuota, f func(name string) bool) error {
	if quora.NoLimit() {
		return nil
	}

	now := time.Now()

	// unlink expired files in place, collect the others
	var filesNotExpired []string

	var errs []error
	_, err := filepath_.GlobFunc(pattern, func(name string) bool {
		fi, err := os.Stat(name)
		if err != nil {
			return false
		}

		fl, err := os.Lstat(name)
		if err != nil {
			return false
		}
		if quora.MaxAge <= 0 {
			filesNotExpired = append(filesNotExpired, name)
			return false
		}

		if now.Sub(fi.ModTime()) < quora.MaxAge {
			filesNotExpired = append(filesNotExpired, name)
			return false
		}

		if fl.Mode()&os.ModeSymlink == os.ModeSymlink {
			return false
		}

		// this file is expired, delete it inplace.
		if f(name) {
			err = os.Remove(name)
			if err != nil {
				errs = append(errs, err)
			}
		}
		return false
	})
	if err != nil {
		return errors.Join(append(errs, err)...)
	}

	if len(filesNotExpired) == 0 {
		return errors.Join(errs...)
	}

	// special case: no file left need to be unlinked, no need to order files.
	if !quora.ExceedCount(len(filesNotExpired)) && quora.NoUsageLimit() {
		return errors.Join(errs...)
	}

	var filesExceedMaxCount []string
	var filesLeftOrdered = filesNotExpired

	// prefer to delete files ordered by ModTime from oldest to newest.
	sort.Sort(rotateFileSlice(filesLeftOrdered))
	if quora.ExceedCount(len(filesLeftOrdered)) {
		removeCount := len(filesLeftOrdered) - quora.MaxCount
		filesExceedMaxCount = filesLeftOrdered[:removeCount]
		filesLeftOrdered = filesLeftOrdered[removeCount:]
	}

	for _, path := range filesExceedMaxCount {
		if f(path) {
			err = os.Remove(path)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	// needGC reports whether disk usage of the file system containing name exceeds quora.
	var needGC = func(name string) bool {
		total, _, avail, inodes, inodesFree, err := DiskUsage(name)
		if err != nil {
			return false
		}
		if total <= 0 {
			return false
		}
		if quora.ExceedBytes(avail, total) {
			return true
		}
		if quora.ExceedInodes(inodes, inodesFree) {
			return true
		}
		return false
	}

	for _, path := range filesLeftOrdered {
		if !needGC(path) {
			break
		}

		if f(path) {
			err = os.Remove(path)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errors.Join(errs...)
}
