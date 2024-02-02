// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logrus

import (
	"io"
	"log/slog"

	slog_ "github.com/searKing/golang/go/log/slog"
	"github.com/sirupsen/logrus"
)

func (f Factory) Apply() error {
	logrus.SetLevel(f.fc.Level)
	logrus.SetReportCaller(f.fc.ReportCaller)
	logrus.AddHook(HookFunc(DefaultSlogHook))
	var slogOpt slog.HandlerOptions
	slogOpt.Level = ToSlogLevel(f.fc.Level)
	slogOpt.AddSource = f.fc.ReportCaller
	slogOpt.ReplaceAttr = slog_.ReplaceAttrTruncate(max(f.fc.TruncateKeySizeTo, f.fc.TruncateValueSizeTo, f.fc.TruncateMessageSizeTo))
	var rotateOpts []slog_.RotateOption
	rotateOpts = append(rotateOpts, slog_.WithRotateRotateInterval(f.fc.RotationDuration),
		slog_.WithRotateMaxCount(f.fc.RotationMaxCount),
		slog_.WithRotateMaxAge(f.fc.RotationMaxAge),
		slog_.WithRotateRotateSize(f.fc.RotationSizeInByte))

	var newer func(path string, opts *slog.HandlerOptions, options ...slog_.RotateOption) (slog.Handler, error)
	switch f.fc.Format {
	case FormatJson:
		newer = slog_.NewRotateJSONHandler
	case FormatText:
		newer = slog_.NewRotateTextHandler
	case FormatGlog:
		newer = slog_.NewRotateGlogHandler
	case FormatGlogHuman:
		newer = slog_.NewRotateGlogHumanHandler
	default:
		// nop if unsupported
		return nil
	}
	var handlers []slog.Handler

	if f.fc.MuteDirectlyOutput {
		logrus.SetOutput(io.Discard)
	}
	if f.fc.MuteDirectlyOutput {
		slogOpt2 := slogOpt
		slogOpt2.Level = ToSlogLevel(f.fc.MuteDirectlyOutputLevel)
		h, err := newer("", &slogOpt2)
		if err != nil {
			logrus.WithField("path", f.fc.Path).
				WithField("duration", f.fc.RotationDuration).
				WithField("max_count", f.fc.RotationMaxCount).
				WithField("max_age", f.fc.RotationMaxAge).
				WithField("rotate_size_in_byte", f.fc.RotationSizeInByte).
				WithField("mute_directly_output", f.fc.MuteDirectlyOutput).
				WithError(err).Error("add rotation wrapper for log")
			return err
		}
		handlers = append(handlers, h)
	}
	{
		h, err := newer(f.fc.Path, &slogOpt, rotateOpts...)
		if err != nil {
			return err
		}
		handlers = append(handlers, h)
	}
	slog.SetDefault(slog.New(slog_.MultiHandler(handlers...)))

	logrus.WithField("path", f.fc.Path).
		WithField("duration", f.fc.RotationDuration).
		WithField("max_count", f.fc.RotationMaxCount).
		WithField("max_age", f.fc.RotationMaxAge).
		WithField("rotate_size_in_byte", f.fc.RotationSizeInByte).
		WithField("mute_directly_output", f.fc.MuteDirectlyOutput).
		Info("add rotation wrapper for log")
	return nil
}
