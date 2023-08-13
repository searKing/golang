// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slog

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"

	os_ "github.com/searKing/golang/go/os"
	time_ "github.com/searKing/golang/go/time"
)

// NewHandler creates a slog.Handler that writes to w,
// using the given options.
// If opts is nil, the default options are used.
type NewHandler func(w io.Writer, opts *slog.HandlerOptions) slog.Handler

// NewRotateHandler creates a slog.Handler that writes to rotate file,
// using the given options.
// If path is empty, the default os.Stdout are used.
// If opts is nil, the default options are used.
func NewRotateHandler(h NewHandler, path string, opts *slog.HandlerOptions, options ...RotateOption) (slog.Handler, error) {
	if path == "" {
		return h(os.Stdout, opts), nil
	}

	if err := os_.MakeAll(filepath.Dir(path)); err != nil {
		return nil, err
	}

	var opt rotate
	opt.FilePathRotateLayout = time_.LayoutStrftimeToSimilarTime(".%Y%m%d%H%M%S.log")
	opt.FileLinkPath = filepath.Base(path) + ".log"
	opt.ApplyOptions(options...)

	file := os_.NewRotateFile(opt.FilePathRotateLayout)
	file.FilePathPrefix = path
	file.FileLinkPath = opt.FileLinkPath
	file.RotateInterval = opt.RotateInterval
	file.RotateSize = opt.RotateSize
	file.MaxAge = opt.MaxAge
	file.MaxCount = opt.MaxCount
	file.ForceNewFileOnStartup = opt.ForceNewFileOnStartup
	file.PostRotateHandler = GlogRotateHeader
	return h(file, opts), nil
}

// NewRotateJSONHandler creates a JSONHandler that writes to rotate file,
// using the given options.
// If opts is nil, the default options are used.
func NewRotateJSONHandler(path string, opts *slog.HandlerOptions, options ...RotateOption) (slog.Handler, error) {
	return NewRotateHandler(func(w io.Writer, opts *slog.HandlerOptions) slog.Handler {
		return slog.NewJSONHandler(w, opts)
	}, path, opts, options...)
}

// NewRotateTextHandler creates a TextHandler that writes to rotate file,
// using the given options.
// If opts is nil, the default options are used.
func NewRotateTextHandler(path string, opts *slog.HandlerOptions, options ...RotateOption) (slog.Handler, error) {
	return NewRotateHandler(func(w io.Writer, opts *slog.HandlerOptions) slog.Handler {
		return slog.NewTextHandler(w, opts)
	}, path, opts, options...)
}

// NewRotateGlogHandler creates a GlogHandler that writes to rotate file,
// using the given options.
// If opts is nil, the default options are used.
// # LOG LINE PREFIX FORMAT
//
// Log lines have this form:
//
//	Lyyyymmdd hh:mm:ss.uuuuuu threadid file:line] msg...
//
// where the fields are defined as follows:
//
//	L                A single character, representing the log level
//	                 (eg 'I' for INFO)
//	yyyy             The year
//	mm               The month (zero padded; ie May is '05')
//	dd               The day (zero padded)
//	hh:mm:ss.uuuuuu  Time in hours, minutes and fractional seconds
//	threadid         The space-padded thread ID as returned by GetTID()
//	                 (this matches the PID on Linux)
//	file             The file name
//	line             The line number
//	msg              The user-supplied message
//
// Example:
//
//	I1103 11:57:31.739339 24395 google.cc:2341] Command line: ./some_prog
//	I1103 11:57:31.739403 24395 google.cc:2342] Process id 24395
func NewRotateGlogHandler(path string, opts *slog.HandlerOptions, options ...RotateOption) (slog.Handler, error) {
	return NewRotateHandler(func(w io.Writer, opts *slog.HandlerOptions) slog.Handler {
		return NewGlogHandler(w, opts)
	}, path, opts, options...)
}

// NewRotateGlogHumanHandler creates a human-readable GlogHandler that writes to rotate file,
// using the given options.
// If opts is nil, the default options are used.
// # LOG LINE PREFIX FORMAT
//
// Log lines have this form:
//
//	[LLLLL] [yyyymmdd hh:mm:ss.uuuuuu] [threadid] [file:line(func)] msg...
//
// where the fields are defined as follows:
//
//	LLLLL            Five characters, representing the log level
//	                 (eg 'INFO ' for INFO)
//	yyyy             The year
//	mm               The month (zero padded; ie May is '05')
//	dd               The day (zero padded)
//	hh:mm:ss.uuuuuu  Time in hours, minutes and fractional seconds
//	threadid         The space-padded thread ID as returned by GetTID()
//	                 (this matches the PID on Linux)
//	file             The file name
//	line             The line number
//	func             The func name
//	msg              The user-supplied message
//
// Example:
//
//	[INFO] [1103 11:57:31.739339] [24395] [google.cc:2341] Command line: ./some_prog
//	[INFO] [1103 11:57:31.739403 24395] [google.cc:2342] Process id 24395
func NewRotateGlogHumanHandler(path string, opts *slog.HandlerOptions, options ...RotateOption) (slog.Handler, error) {
	return NewRotateHandler(func(w io.Writer, opts *slog.HandlerOptions) slog.Handler {
		return NewGlogHumanHandler(w, opts)
	}, path, opts, options...)
}
