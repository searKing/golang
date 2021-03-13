// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logrus

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type HookFunc func(entry *logrus.Entry) error

func (f HookFunc) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (f HookFunc) Fire(entry *logrus.Entry) error {
	return f(entry)
}

type HookMap map[logrus.Level]func(entry *logrus.Entry) error

func (hooks HookMap) Levels() []logrus.Level {
	var levels []logrus.Level
	for lvl := range hooks {
		levels = append(levels, lvl)
	}
	return levels
}

func (hooks HookMap) Fire(entry *logrus.Entry) error {
	hook, has := hooks[entry.Level]
	if !has {
		return fmt.Errorf("level[%s] not registered", entry.Level)
	}
	return hook(entry)
}
