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
	mu     sync.Mutex
}

var stdFieldLogger = New(nil)

func New(l logrus.FieldLogger) *FieldLogger {
	return &FieldLogger{
		logger: l,
	}
}

func (b *FieldLogger) Write(p []byte) (n int, err error) {
	b.GetLogger().Printf("%s", string(p))
	return len(p), nil
}

func (b *FieldLogger) GetStdLogger() *log.Logger {
	return log.New(New(b.GetLogger()), "", 0)
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
