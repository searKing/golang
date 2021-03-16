// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logrus

import (
	"time"

	time_ "github.com/searKing/golang/go/time"
)

// WithRotateLayout sets time layout to format rotate file.
func WithRotateLayout(layout string) RotateOption {
	return RotateOptionFunc(func(r *rotate) {
		r.FilePathRotateLayout = layout
	})
}

// WithRotateStrftime sets time layout in strftime format to format rotate file.
func WithRotateStrftime(layout string) RotateOption {
	return RotateOptionFunc(func(r *rotate) {
		r.FilePathRotateLayout = time_.LayoutStrftimeToSimilarTime(layout)
	})
}

// WithFileLinkPath sets the symbolic link name that gets linked to the current file name being used.
func WithFileLinkPath(link string) RotateOption {
	return RotateOptionFunc(func(r *rotate) {
		r.FileLinkPath = link
	})
}

// WithRotateInterval rotates files are rotated until RotateInterval expired before being removed
// take effects if only RotateInterval is bigger than 0.
func WithRotateInterval(interval time.Duration) RotateOption {
	return RotateOptionFunc(func(r *rotate) {
		r.RotateInterval = interval
	})
}

// WithRotateSize rotates files are rotated if they grow bigger then size bytes.
// take effects if only RotateSize is bigger than 0.
func WithRotateSize(size int64) RotateOption {
	return RotateOptionFunc(func(r *rotate) {
		r.RotateSize = size
	})
}

// WithMaxAge sets max age of a log file before it gets purged from the file system.
// Remove rotated logs older than duration. The age is only checked if the file is
// to be rotated.
// take effects if only MaxAge is bigger than 0.
func WithMaxAge(maxAge time.Duration) RotateOption {
	return RotateOptionFunc(func(r *rotate) {
		r.MaxAge = maxAge
	})
}

// WithMaxCount rotates files are rotated MaxCount times before being removed
// take effects if only MaxCount is bigger than 0.
func WithMaxCount(maxCount int) RotateOption {
	return RotateOptionFunc(func(r *rotate) {
		r.MaxCount = maxCount
	})
}

// WithMaxCount force rotates files directly when start up
func WithForceNewFileOnStartup(force bool) RotateOption {
	return RotateOptionFunc(func(r *rotate) {
		r.ForceNewFileOnStartup = force
	})
}

// WithForceNewFileOnStartup mutes writer of logrus.Output if level is InfoLevel、DebugLevel、TraceLevel...
func WithMuteDirectlyOutput(mute bool) RotateOption {
	return RotateOptionFunc(func(r *rotate) {
		r.MuteDirectlyOutput = mute
	})
}
