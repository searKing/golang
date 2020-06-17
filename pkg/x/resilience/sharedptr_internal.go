// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resilience

import (
	"context"
	"time"

	logrus2 "github.com/searKing/golang/third_party/github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus"
)

type sharedPtr struct {
	// New optionally specifies a function to generate
	// a value when Get would otherwise return nil.
	// It may not be changed concurrently with calls to Get.
	// readonly
	New func() (Ptr, error)
	*logrus2.FieldLogger

	// to judge whether Get&Construct is timeout
	// readonly
	ConstructTimeout time.Duration
	// MaxDuration for retry if tasks failed
	// readonly
	TaskMaxRetryDuration time.Duration

	// ctx is either the client or server context. It should only
	// be modified via copying the whole Request using WithContext.
	// It is unexported to prevent people from using Context wrong
	// and mutating the contexts held by callers of the same request.
	ctx context.Context

	x      Ptr
	taskC  chan *Task
	tasks  map[string]*Task
	eventC chan Event
}

func newSharedPtrSafe(ctx context.Context, new func() (Ptr, error), l logrus.FieldLogger) *sharedPtr {
	return &sharedPtr{
		New:                  new,
		FieldLogger:          logrus2.New(l),
		TaskMaxRetryDuration: DefaultResilienceTaskMaxRetryDuration,
		ConstructTimeout:     DefaultResilienceConstructTimeout,
		ctx:                  ctx,
	}
}

//
// The returned context is always non-nil; it defaults to the
// background context.
func (sp *sharedPtr) Context() context.Context {
	if sp.ctx != nil {
		return sp.ctx
	}
	return context.Background()
}

func (sp *sharedPtr) InShutdown() bool {
	select {
	case <-sp.Context().Done():
		return true
	default:
		return false
	}
}

// std::shared_ptr.release()
func (sp *sharedPtr) Release() Ptr {
	x := sp.x
	sp.x = nil
	return x
}

// std::shared_ptr.reset()
func (sp *sharedPtr) Reset(ptr Ptr) {
	sp.resetPtr(ptr)
}

// std::shared_ptr.get()
func (sp *sharedPtr) Get() Ptr {
	return sp.x
}

// reset ptr and ready to start again
func (sp *sharedPtr) resetPtr(ptr Ptr) {
	oldPtr := sp.Release()
	if oldPtr != nil {
		oldPtr.Close()
	}
	sp.x = ptr

	return
}

func (sp *sharedPtr) GetTaskC() chan *Task {
	if sp.taskC == nil {
		sp.taskC = make(chan *Task)
	}
	return sp.taskC
}

func (sp *sharedPtr) Event() chan Event {
	if sp.eventC == nil {
		sp.eventC = make(chan Event)
	}
	return sp.eventC
}

func (sp *sharedPtr) RemoveAllTask() {
	for _, id := range sp.TaskIds() {
		sp.RemoveTaskById(id, true)
	}
}

func (sp *sharedPtr) TaskIds() []string {
	var ids []string
	for id := range sp.tasks {
		ids = append(ids, id)
	}
	return ids
}

func (sp *sharedPtr) SnapTasks() map[string]*Task {
	tasks := sp.tasks
	sp.tasks = nil
	return tasks
}

func (sp *sharedPtr) GetTasks() map[string]*Task {
	return sp.tasks
}

func (sp *sharedPtr) GetTaskById(id string) (task *Task, has bool) {
	task, has = sp.tasks[id]
	return
}

func (sp *sharedPtr) AddTask(task *Task) (has bool) {
	if task == nil {
		return
	}
	if sp.tasks == nil {
		sp.tasks = make(map[string]*Task)
	}

	sp.tasks[task.ID()] = task
	return
}

func (sp *sharedPtr) RemoveTask(task *Task) {
	if task == nil {
		return
	}
	sp.RemoveTaskById(task.ID(), true)
	return
}

func (sp *sharedPtr) RemoveAllTasks() {
	for _, id := range sp.TaskIds() {
		sp.RemoveTaskById(id, true)
	}
}

func (sp *sharedPtr) RemoveTaskById(id string, cancel bool) {
	t, has := sp.tasks[id]
	if !has || t == nil {
		return
	}
	if cancel && t.CancelFn != nil {
		t.CancelFn()
	}
	delete(sp.tasks, t.ID())
	return
}
