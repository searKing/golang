// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slog

import "log/slog"

// https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_(Select_Graphic_Rendition)_parameters
// https://en.wikipedia.org/wiki/ANSI_escape_code#8-bit
const (
	reset   = "\x1b[0m"
	black   = "\x1b[30m"
	red     = "\x1b[31m"
	green   = "\x1b[32m"
	yellow  = "\x1b[33m"
	blue    = "\x1b[34m"
	magenta = "\x1b[35m"
	cyan    = "\x1b[36m"
	white   = "\x1b[37m"
	gray    = "\x1b[90m"
)

func levelColor(level slog.Level) string {
	switch level {
	case slog.LevelDebug:
		return gray
	case slog.LevelWarn:
		return yellow
	case slog.LevelError:
		return red
	default:
		return blue
	}
}
