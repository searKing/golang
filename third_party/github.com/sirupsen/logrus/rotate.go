// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logrus

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	os_ "github.com/searKing/golang/go/os"
	filepath_ "github.com/searKing/golang/go/path/filepath"
)

// WithRotation enhances logrus log to be written to local filesystem, with file rotation
// path sets log's base path prefix
// duration sets the time between rotation.
// maxCount sets the number of files should be kept before it gets purged from the file system.
// maxAge sets the max age of a log file before it gets purged from the file system.
func WithRotation(log *logrus.Logger, path string, duration time.Duration, maxCount int, maxAge time.Duration) error {
	if log == nil {
		return nil
	}
	dir := filepath_.ToDir(filepath.Dir(path))
	if err := filepath_.TouchAll(dir, filepath_.PrivateDirMode); err != nil {
		return errors.Wrap(err, fmt.Sprintf("create dir[%s] for log failed", dir))
	}

	file := os_.NewRotateFileWithStrftime(".%Y%m%d%H%M.log")
	file.FilePathPrefix = path
	file.FileLinkPath = filepath.Base(path) + ".log"
	file.RotateInterval = duration
	file.MaxCount = maxCount
	file.MaxAge = maxAge
	log.AddHook(HookFunc(func(entry *logrus.Entry) error {
		var msg []byte
		var err error

		if log.Formatter == nil {
			msg_, err_ := entry.String()
			msg, err = []byte(msg_), err_
		} else {
			msg, err = log.Formatter.Format(entry)
		}

		if err != nil {
			return err
		}
		_, err = file.Write(msg)
		return err
	}))
	return nil
}
