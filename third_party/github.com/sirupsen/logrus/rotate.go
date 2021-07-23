// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logrus

import (
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"

	os_ "github.com/searKing/golang/go/os"
	time_ "github.com/searKing/golang/go/time"
)

//go:generate go-option -type "rotate"
type rotate struct {
	// Time layout to format rotate file
	FilePathRotateLayout string

	// sets the symbolic link name that gets linked to the current file name being used.
	FileLinkPath string

	// Rotate files are rotated until RotateInterval expired before being removed
	// take effects if only RotateInterval is bigger than 0.
	RotateInterval time.Duration

	// Rotate files are rotated if they grow bigger then size bytes.
	// take effects if only RotateSize is bigger than 0.
	RotateSize int64

	// max age of a log file before it gets purged from the file system.
	// Remove rotated logs older than duration. The age is only checked if the file is
	// to be rotated.
	// take effects if only MaxAge is bigger than 0.
	MaxAge time.Duration

	// Rotate files are rotated MaxCount times before being removed
	// take effects if only MaxCount is bigger than 0.
	MaxCount int

	// Force File Rotate when start up
	ForceNewFileOnStartup bool

	// mute writer of logrus.Output if level is less or equal than MuteDirectlyOutputLogLevel
	// default is true
	MuteDirectlyOutput bool

	// sets the logger level of directly output.
	// default is warn
	MuteDirectlyOutputLogLevel logrus.Level
}

// WithRotate enhances logrus log to be written to local filesystem, with file rotation
// path sets log's base path prefix
func WithRotate(log *logrus.Logger, path string, options ...RotateOption) error {
	if log == nil {
		return nil
	}
	if err := os_.MakeAll(filepath.Dir(path)); err != nil {
		return err
	}

	var opt rotate
	opt.FilePathRotateLayout = time_.LayoutStrftimeToSimilarTime(".%Y%m%d%H%M%S.log")
	opt.FileLinkPath = filepath.Base(path) + ".log"
	opt.MuteDirectlyOutput = true
	opt.MuteDirectlyOutputLogLevel = logrus.WarnLevel
	opt.ApplyOptions(options...)

	file := os_.NewRotateFile(opt.FilePathRotateLayout)
	file.FilePathPrefix = path
	file.FileLinkPath = opt.FileLinkPath
	file.RotateInterval = opt.RotateInterval
	file.RotateSize = opt.RotateSize
	file.MaxAge = opt.MaxAge
	file.MaxCount = opt.MaxCount
	file.ForceNewFileOnStartup = opt.ForceNewFileOnStartup

	var out = ioutil.Discard
	if opt.MuteDirectlyOutput {
		out = log.Out
		log.SetOutput(ioutil.Discard)
	}
	log.AddHook(HookFunc(func(entry *logrus.Entry) error {
		var msg []byte
		var err error

		if log.Formatter == nil {
			msg_, err_ := entry.String()
			msg, err = []byte(msg_), err_
		} else {
			switch f := log.Formatter.(type) {
			case *logrus.TextFormatter:
				var disableColors = f.DisableColors
				// disable colors in log file
				f.DisableColors = true
				msg, err = log.Formatter.Format(entry)
				f.DisableColors = disableColors
			default:
				msg, err = log.Formatter.Format(entry)
			}
		}

		if err != nil {
			return err
		}

		if opt.MuteDirectlyOutput && entry.Level <= opt.MuteDirectlyOutputLogLevel {
			if out != nil {
				_, _ = out.Write(msg)
			}
		}
		_, err = file.Write(msg)
		return err
	}))
	return nil
}
