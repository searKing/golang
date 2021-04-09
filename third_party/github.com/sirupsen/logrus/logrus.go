// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logrus

import (
	"io"
	"log"

	"github.com/sirupsen/logrus"

	log_ "github.com/searKing/golang/go/log"
)

func StandardWriter(level logrus.Level) io.Writer {
	return Writer(logrus.StandardLogger(), level)
}

func Writer(l logrus.FieldLogger, level logrus.Level) io.Writer {
	var f log_.PrintfFunc
	switch level {
	case logrus.PanicLevel:
		f = l.Panicf
	case logrus.FatalLevel:
		f = l.Fatalf
	case logrus.ErrorLevel:
		f = l.Errorf
	case logrus.WarnLevel:
		f = l.Warnf
	case logrus.InfoLevel:
		f = l.Infof
	case logrus.DebugLevel:
		f = l.Debugf
	case logrus.TraceLevel:
		if l, ok := l.(logrus.Ext1FieldLogger); ok {
			f = l.Tracef
		} else {
			f = l.Printf
		}
	default:
		f = l.Printf
	}
	return f
}

// AsStdLogger returns *log.Logger from logrus.FieldLogger
func AsStdLogger(l logrus.FieldLogger, level logrus.Level, prefix string, flag int) *log.Logger {
	return log.New(Writer(l, level), prefix, flag)
}

// AsStdLoggerWithLevel is only a helper of AsStdLogger
func AsStdLoggerWithLevel(l logrus.FieldLogger, level logrus.Level) *log.Logger {
	return AsStdLogger(l, level, "", 0)
}
