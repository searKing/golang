// Copyright 2025 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.25

package slog

import (
	"log/slog"
)

// source returns a Source for the log event.
// If the Record was created without the necessary information,
// or if the location is unavailable, it returns a non-nil *Source
// with zero fields.
func source(r slog.Record) *slog.Source {
	src := r.Source()
	if src == nil {
		src = &slog.Source{}
	}
	return src
}

// ShortSource returns a Source for the log event.
// If the Record was created without the necessary information,
// or if the location is unavailable, it returns a non-nil *Source
// with zero fields.
func ShortSource(r slog.Record) *slog.Source {
	src := r.Source()
	if src == nil {
		src = &slog.Source{}
	}
	return &slog.Source{
		Function: shortFunction(src.Function),
		File:     shortFile(src.File),
		Line:     src.Line,
	}
}
