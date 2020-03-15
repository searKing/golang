// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resilience

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type TaskType struct {
	Drop      bool // task will be dropped
	Retry     bool // task will be retried if error happens
	Construct bool // task will be called after New
	Repeat    bool // Task will be executed again and again
}

func (t TaskType) String() string {
	var b strings.Builder
	if t.Drop {
		b.WriteString("drop")
	}
	if t.Retry {
		if b.String() != "" {
			b.WriteRune('-')
		}
		b.WriteString("retry")
	}
	if t.Construct {
		if b.String() != "" {
			b.WriteRune('-')
		}
		b.WriteString("construct")
	}
	if t.Repeat {
		if b.String() != "" {
			b.WriteRune('-')
		}
		b.WriteString("repeat")
	}
	return b.String()
}

type Task struct {
	Type        TaskType
	State       TaskState
	Description string // for debug
	Handle      func() error

	RepeatDuration time.Duration
	RetryDuration  time.Duration

	Ctx        context.Context
	CancelFn   context.CancelFunc
	inShutdown bool
}

//
// The returned context is always non-nil; it defaults to the
// background context.
func (t *Task) Context() context.Context {
	if t.Ctx != nil {
		return t.Ctx
	}
	return context.Background()
}

func (t *Task) String() string {
	if t == nil {
		return "empty task"
	}
	return fmt.Sprintf("%s-%s-%s", t.ID(), t.State, t.Description)
}

func (t *Task) ID() string {
	if t == nil {
		return "empty task"
	}
	return fmt.Sprintf("%s-%p", t.Type, t.Handle)
}

func (t *Task) Clone() *Task {
	if t == nil {
		return nil
	}
	task := *t
	return &task
}
