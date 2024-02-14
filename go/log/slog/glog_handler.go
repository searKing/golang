// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slog

import (
	"context"
	"io"
	"log/slog"
	"sync"
)

var _ slog.Handler = (*GlogHandler)(nil)

type GlogHandler struct {
	*commonHandler
}

// NewGlogHandler creates a GlogHandler that writes to w,
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
func NewGlogHandler(w io.Writer, opts *slog.HandlerOptions) *GlogHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{
			AddSource: true,
		}
	}
	return &GlogHandler{
		&commonHandler{
			ReplaceLevelString: func(l slog.Level) string {
				switch {
				case l < slog.LevelInfo:
					return "D"
				case l < slog.LevelWarn:
					return "I"
				case l < slog.LevelError:
					return "W"
				default:
					return "E"
				}
			},
			opts:           *opts,
			AttrSep:        ", ",
			DisableQuote:   true,
			SourcePrettier: ShortSource,
			sharedVar:      &sharedVar{once: &sync.Once{}},
			mu:             &sync.Mutex{},
			w:              w,
		},
	}
}

// NewGlogHumanHandler creates a human-readable GlogHandler that writes to w,
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
//	[INFO ] [20081103 11:57:31.739339] [24395] [google.cc:2341](main) Command line: ./some_prog
//	[INFO ] [20081103 11:57:31.739403] [24395] [google.cc:2342](main) Process id 24395
func NewGlogHumanHandler(w io.Writer, opts *slog.HandlerOptions) *GlogHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{
			AddSource: true,
		}
	}
	return &GlogHandler{
		&commonHandler{
			opts:           *opts,
			AttrSep:        ", ",
			DisableQuote:   true,
			PadLevelText:   true,
			HumanReadable:  true,
			SourcePrettier: ShortSource,
			WithFuncName:   true,
			sharedVar:      &sharedVar{once: &sync.Once{}},
			mu:             &sync.Mutex{},
			w:              w,
		},
	}
}

// Enabled reports whether the handler handles records at the given level.
// The handler ignores records whose level is lower.
func (h *GlogHandler) Enabled(_ context.Context, level slog.Level) bool {
	return h.commonHandler.enabled(level)
}

// WithAttrs returns a new GlogHandler whose attributes consists
// of h's attributes followed by attrs.
func (h *GlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &GlogHandler{commonHandler: h.commonHandler.withAttrs(attrs)}
}

func (h *GlogHandler) WithGroup(name string) slog.Handler {
	return &GlogHandler{commonHandler: h.commonHandler.withGroup(name)}
}

// Handle formats its argument Record as a JSON object on a single line.
//
// If the Record's time is zero, the time is omitted.
// Otherwise, the key is "time"
// and the value is output as with GlogDate.
//
// If the Record's level is zero, the level is omitted.
// Otherwise, the key is "level"
// and the value of [Level.String] is output.
//
// If the AddSource option is set and source information is available,
// the key is "source", and the value is a record of type [Source].
//
// The message's key is "msg".
//
// To modify these or other attributes, or remove them from the output, use
// [HandlerOptions.ReplaceAttr].
//
// Values are formatted as with an [encoding/json.Encoder] with SetEscapeHTML(false),
// with two exceptions.
//
// First, an Attr whose Value is of type error is formatted as a string, by
// calling its Error method. Only errors in Attrs receive this special treatment,
// not errors embedded in structs, slices, maps or other data structures that
// are processed by the encoding/json package.
//
// Second, an encoding failure does not cause Handle to return an error.
// Instead, the error message is formatted as a string.
//
// Each call to Handle results in a single serialized call to io.Writer.Write.
//
// Header formats a log header as defined by the C++ implementation.
// It returns a buffer containing the formatted header and the user's file and line number.
// The depth specifies how many stack frames above lives the source line to be identified in the log message.
//
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
//
// NOTE: although the microseconds are useful for comparing events on
// a single machine, clocks on different machines may not be well
// synchronized.  Hence, use caution when comparing the low bits of
// timestamps from different machines.
func (h *GlogHandler) Handle(_ context.Context, r slog.Record) error {
	return h.commonHandler.handle(r)
}
