// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logrus

import (
	"log"
	"sync"

	"github.com/sirupsen/logrus"
)

type FieldLogger struct {
	logger logrus.FieldLogger
	level  logrus.Level
	mu     sync.Mutex
}

var stdFieldLogger = New(nil)

func New(l logrus.FieldLogger) *FieldLogger {
	return &FieldLogger{
		logger: l,
		level:  logrus.InfoLevel,
	}
}

func (b *FieldLogger) Clone() *FieldLogger {
	b.mu.Lock()
	defer b.mu.Unlock()
	return &FieldLogger{
		logger: b.logger,
		level:  b.level,
	}
}

func (b *FieldLogger) Write(p []byte) (n int, err error) {
	b.mu.Lock()
	var level = b.level
	b.mu.Unlock()
	switch level {
	case logrus.PanicLevel:
		b.GetLogger().Panicf("%s", string(p))
	case logrus.FatalLevel:
		b.GetLogger().Fatalf("%s", string(p))
	case logrus.ErrorLevel:
		b.GetLogger().Errorf("%s", string(p))
	case logrus.WarnLevel:
		b.GetLogger().Warnf("%s", string(p))
	case logrus.InfoLevel:
		b.GetLogger().Infof("%s", string(p))
	case logrus.DebugLevel:
		b.GetLogger().Debugf("%s", string(p))
	case logrus.TraceLevel:
		logger := b.GetLogger()
		if logger_, ok := logger.(logrus.Ext1FieldLogger); ok {
			logger_.Tracef("%s", string(p))
		}
	}
	return len(p), nil
}

func (b *FieldLogger) GetStdLogger() *log.Logger {
	return log.New(New(b.GetLogger()), "", 0)
}

func (b *FieldLogger) GetStdLoggerWithLevel(level logrus.Level) *log.Logger {
	logger := b.Clone()
	logger.level = level
	return log.New(logger, "", 0)
}

func (b *FieldLogger) GetLogger() logrus.FieldLogger {
	if b == nil {
		return stdFieldLogger.GetLogger()
	}
	if b.logger != nil {
		return b.logger
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.logger == nil {
		b.logger = logrus.StandardLogger()
		b.logger.Warning("No logger was set, defaulting to standard logger.")
	}
	return b.logger
}

func (b *FieldLogger) SetStdLogger(l *log.Logger) {
	if l == nil {
		return
	}
	logger := logrus.New()
	logger.Out = l.Writer()
	b.SetLogger(logger)
}

func (b *FieldLogger) SetLogger(l logrus.FieldLogger) {
	if b == nil {
		return
	}
	if l == nil {
		return
	}

	b.mu.Lock()
	defer b.mu.Unlock()
	b.logger = l
}
