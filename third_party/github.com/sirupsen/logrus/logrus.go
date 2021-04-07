// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logrus

import (
	"log"

	"github.com/sirupsen/logrus"

	log_ "github.com/searKing/golang/go/log"
)

// AsStdLogger returns *log.Logger from logrus.FieldLogger
func AsStdLogger(l logrus.FieldLogger, level logrus.Level, prefix string, flag int) *log.Logger {
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
	return log.New(f, prefix, flag)
}
