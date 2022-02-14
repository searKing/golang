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

// UnlinkOldestFiles unlink old files if need
func UnlinkOldestFiles(pattern string, maxAge time.Duration, maxCount int, used, iUsed float32) error {
	if maxAge <= 0 && maxCount <= 0 && used <= 0 && iUsed <= 0 {
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
		if maxAge <= 0 {
			filesNotExpired = append(filesNotExpired, name)
			return false
		}

		if now.Sub(fi.ModTime()) < maxAge {
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
	if maxCount > 0 && len(filesNotExpired) > 0 {
		removeCount := len(filesNotExpired) - maxCount
		if removeCount < 0 {
			removeCount = 0
		}
		sort.Sort(rotateFileSlice(filesNotExpired))
		filesExceedMaxCount = filesNotExpired[:removeCount]
		filesLeft = filesNotExpired[removeCount:]
	}
	var errs []error
	for _, path := range filesExpired {
		err = os.Remove(path)
		if err != nil {
			errs = append(errs, err)
		}
	}
	for _, path := range filesExceedMaxCount {
		err = os.Remove(path)
		if err != nil {
			errs = append(errs, err)
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

		if used > 0 && float32(avail)/float32(total) > used {
			return true
		}
		if iUsed > 0 && float32(inodes-inodesFree)/float32(inodes) > iUsed {
			return true
		}
		return false
	}

	for _, path := range filesLeft {
		if !needGC(path) {
			return nil
		}
		err = os.Remove(path)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errors_.Multi(errs...)
}
