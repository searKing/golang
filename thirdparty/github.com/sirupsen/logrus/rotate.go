// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logrus

import (
	"fmt"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	filepath_ "github.com/searKing/golang/go/path/filepath"
	"github.com/sirupsen/logrus"
)

// WithRotation enhances logrus log to be written to local filesystem, with file rotation
// path sets log's base path prefix
// duration sets the time between rotation.
// maxCount sets the number of files should be kept before it gets purged from the file system.
// maxAge sets the max age of a log file before it gets purged from the file system.
func WithRotation(log *logrus.Logger, path string, duration time.Duration, maxCount uint, maxAge time.Duration) error {
	if log == nil {
		return nil
	}
	dir := filepath_.ToDir(filepath.Dir(path))
	if err := filepath_.TouchAll(dir, filepath_.PrivateDirMode); err != nil {
		return errors.Wrap(err, fmt.Sprintf("create dir[%s] for log failed", dir))
	}

	writer, err := rotatelogs.New(
		path+".%Y%m%d%H%M.log",
		rotatelogs.WithLinkName(filepath.Base(path+".log")), // 生成软链，指向最新日志文件
		rotatelogs.WithRotationTime(duration),               // 日志切割时间间隔
		rotatelogs.WithRotationCount(maxCount),              // 文件片段最大保存个数
		rotatelogs.WithMaxAge(maxAge),                       // 文件最大保存时间
	)
	if err != nil {
		return errors.Wrap(err, "create rotate logs failed")
	}
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer, // 为不同级别设置不同的输出目的
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, log.Formatter)
	log.AddHook(lfHook)
	return nil
}
