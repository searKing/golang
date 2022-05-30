// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package os

import (
	"os"
	"sort"
	"time"

	errors_ "github.com/searKing/golang/go/errors"
	filepath_ "github.com/searKing/golang/go/path/filepath"
)

type DiskQuota struct {
	MaxAge             time.Duration // max age of files
	MaxCount           int           // max count of files
	MaxUsedProportion  float32       // max used proportion of files
	MaxIUsedProportion float32       // max used proportion of inodes
}

func (q DiskQuota) NoLimit() bool {
	return q.MaxAge <= 0 && q.MaxCount <= 0 && q.MaxUsedProportion <= 0 && q.MaxIUsedProportion <= 0
}

func (q DiskQuota) ExceedBytes(avail, total int64) bool {
	return q.MaxUsedProportion > 0 && float32(total-avail) > q.MaxUsedProportion*float32(total)
}

func (q DiskQuota) ExceedInodes(inodes, inodesFree int64) bool {
	return q.MaxIUsedProportion > 0 && float32(inodes-inodesFree) > q.MaxIUsedProportion*float32(inodes)
}

// UnlinkOldestFiles unlink old files if need
func UnlinkOldestFiles(pattern string, quora DiskQuota) error {
	return UnlinkOldestFilesFunc(pattern, quora, func(name string) bool { return true })
}

// UnlinkOldestFilesFunc unlink old files satisfying f(c) if need
func UnlinkOldestFilesFunc(pattern string, quora DiskQuota, f func(name string) bool) error {
	if quora.NoLimit() {
		return nil
	}

	now := time.Now()

	// find old files
	var filesNotExpired []string
	filesExpired, err := filepath_.GlobFunc(pattern, func(name string) bool {
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
		return true
	})
	if err != nil {
		return err
	}

	var filesExceedMaxCount []string
	var filesLeft = filesNotExpired
	if quora.MaxCount > 0 && len(filesNotExpired) > 0 {
		removeCount := len(filesNotExpired) - quora.MaxCount
		if removeCount < 0 {
			removeCount = 0
		}
		sort.Sort(rotateFileSlice(filesNotExpired))
		filesExceedMaxCount = filesNotExpired[:removeCount]
		filesLeft = filesNotExpired[removeCount:]
	}
	if f == nil {
		f = func(name string) bool {
			return true
		}
	}

	var errs []error
	for _, path := range filesExpired {
		if f(path) {
			err = os.Remove(path)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	for _, path := range filesExceedMaxCount {
		if f(path) {
			err = os.Remove(path)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

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
		if quora.ExceedInodes(inodes-inodesFree, inodes) {
			return true
		}
		return false
	}

	for _, path := range filesLeft {
		if !needGC(path) {
			return nil
		}

		if f(path) {
			err = os.Remove(path)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errors_.Multi(errs...)
}
