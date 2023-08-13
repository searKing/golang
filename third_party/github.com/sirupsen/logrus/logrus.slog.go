// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logrus

import (
	"log/slog"

	"github.com/sirupsen/logrus"
)

func ToSlogLevel(l logrus.Level) slog.Level {
	switch l {
	case logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel:
		return slog.LevelError
	case logrus.WarnLevel:
		return slog.LevelWarn
	case logrus.InfoLevel:
		return slog.LevelInfo
	case logrus.DebugLevel, logrus.TraceLevel:
		return slog.LevelDebug
	default:
		if l < logrus.PanicLevel {
			return slog.LevelError
		}
		return slog.LevelDebug
	}
}
