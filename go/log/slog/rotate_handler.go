// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slog

import (
	"io"
	"log/slog"
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
// If opts is nil, the default options are used.
func NewRotateHandler(h NewHandler, path string, opts *slog.HandlerOptions, options ...RotateOption) (slog.Handler, error) {
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
