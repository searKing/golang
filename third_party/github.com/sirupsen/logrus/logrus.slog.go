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

// DefaultSlogHook returns a logrus Hook by [slog.Default],
// followed even if [slog.SetDefault] changes default slog logger.
func DefaultSlogHook(entry *logrus.Entry) error {
	return SlogHook(slog.Default().Handler())(entry)
}

// SlogHook returns a logrus Hook by [slog.Handler]
func SlogHook(h slog.Handler) func(entry *logrus.Entry) error {
	return func(entry *logrus.Entry) error {
		var pc uintptr
		// caller of entry.Caller
		if entry.Caller != nil {
			pc = entry.Caller.PC + 1
		}

		var attrs []slog.Attr
		for k, v := range entry.Data {
			attrs = append(attrs, slog.Any(k, v))
		}

		r := slog.NewRecord(entry.Time, ToSlogLevel(entry.Level), entry.Message, pc)
		r.AddAttrs(attrs...)
		return h.Handle(entry.Context, r)
	}
}
