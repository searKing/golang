// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slog

import (
	"time"
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
}
