// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slog

import (
	"fmt"
	"log/slog"
	"path"
	"runtime"
	"strings"
)

// source returns a Source for the log event.
// If the Record was created without the necessary information,
// or if the location is unavailable, it returns a non-nil *Source
// with zero fields.
func source(r slog.Record) *slog.Source {
	fs := runtime.CallersFrames([]uintptr{r.PC})
	f, _ := fs.Next()
	return &slog.Source{
		Function: f.Function,
		File:     f.File,
		Line:     f.Line,
	}
}

// ShortSource returns a Source for the log event.
// If the Record was created without the necessary information,
// or if the location is unavailable, it returns a non-nil *Source
// with zero fields.
func ShortSource(r slog.Record) *slog.Source {
	fs := runtime.CallersFrames([]uintptr{r.PC})
	f, _ := fs.Next()

	function := path.Base(f.Function)
	slash := strings.LastIndex(function, ".")
	if slash >= 0 && slash+1 < len(function) {
		function = function[slash+1:]
	}
	file := path.Base(f.File)
	if file == "" {
		file = "???"
	}
	return &slog.Source{
		Function: function,
		File:     path.Base(f.File),
		Line:     f.Line,
	}
}

// ShortCallerPrettyfier modify the content of the function and
// file keys in the data when ReportCaller is activated.
// INFO[0000] main.go:23 main() hello world
var ShortCallerPrettyfier = func(f *runtime.Frame) (function string, file string) {
	funcname := path.Base(f.Function)
	filename := path.Base(f.File)
	return fmt.Sprintf("%s()", funcname), fmt.Sprintf("%s:%d", filename, f.Line)
}

// attrs returns the non-zero fields of s as a slice of attrs.
// It is similar to a LogValue method, but we don't want Source
// to implement LogValuer because it would be resolved before
// the ReplaceAttr function was called.
func sourceAsGroup(s *slog.Source) slog.Value {
	var as []slog.Attr
	if s.Function != "" {
		as = append(as, slog.String("function", s.Function))
	}
	if s.File != "" {
		as = append(as, slog.String("file", s.File))
	}
	if s.Line != 0 {
		as = append(as, slog.Int("line", s.Line))
	}
	return slog.GroupValue(as...)
}
