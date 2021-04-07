// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resilience

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/searKing/golang/go/sync/atomic"
	"github.com/sirupsen/logrus"
)

const (
	// DefaultConnectTimeout is the default timeout to establish a connection to
	// a ZooKeeper node.
	DefaultResilienceConstructTimeout = 0
	// DefaultSessionTimeout is the default timeout to keep the current
	// ZooKeeper session alive during a temporary disconnect.
	DefaultResilienceTaskMaxRetryDuration = 15 * time.Second

	DefaultTaskRetryTimeout      = 1 * time.Second
	DefaultTaskRescheduleTimeout = 1 * time.Second
)

var (
	ErrEmptyValue       = fmt.Errorf("empty value")
	ErrAlreadyShutdown  = fmt.Errorf("already shutdown")
	ErrNotReady         = fmt.Errorf("not ready")
	ErrAlreadyAddedTask = fmt.Errorf("task is already added")
)

type SharedPtr struct {
	ErrorLog logrus.FieldLogger

	sp *sharedPtr
	mu sync.Mutex

	backgroundStopped atomic.Bool
	WatchStopped      atomic.Bool
}

func NewSharedPtr(ctx context.Context, new func() (Ptr, error)) *SharedPtr {
	return &SharedPtr{
		sp: newSharedPtrSafe(ctx, new),
	}
}

//
// The returned context is always non-nil; it defaults to the
// background context.
func (g *SharedPtr) Context() context.Context {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.sp.Context()
}

func (g *SharedPtr) logger() logrus.FieldLogger {
	if g.ErrorLog != nil {
		return g.ErrorLog
	}
	return logrus.StandardLogger()
}

func (g *SharedPtr) InShutdown() bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.sp.InShutdown()
}

func (g *SharedPtr) GetTaskById(id string) (task *Task, has bool) {
	//g.mu.Lock()
	//defer g.mu.Unlock()
	return g.sp.GetTaskById(id)
}

func (g *SharedPtr) GetTasks() map[string]*Task {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.sp.GetTasks()
}

func (g *SharedPtr) AddTask(task *Task) error {
	if task == nil {
		go g.logger().WithError(ErrEmptyValue).
			Error("task is nil to add, ignore it...")
		return ErrEmptyValue
	}
	if task.Handle == nil {
		go g.logger().WithField("task", task.String()).WithError(ErrEmptyValue).
			Error("task is nonsense to add, ignore it...")
		return ErrEmptyValue
	}
	if g.InShutdown() {
		go g.logger().WithField("task", task.String()).
			Error("resilience is shutdown  already, ignore it...")
		return ErrAlreadyShutdown
	}
	task.Ctx, task.CancelFn = context.WithCancel(g.Context())

	if _, addedTask := g.GetTaskById(task.ID()); addedTask {
		go g.logger().WithField("task", task.String()).
			Error("task is added already, ignore it...")
		return ErrAlreadyAddedTask
	}

	go g.backgroundTask()
	go func() {
		go g.logger().WithField("task", task.String()).Info("new task is adding...")
		g.getTaskC() <- task
	}()
	return nil
}
func (g *SharedPtr) AddTaskFunc(taskType TaskType, handle func() error, descriptions ...string) error {
	return g.AddTask(&Task{
		Description:    strings.Join(descriptions, ""),
		Type:           taskType,
		Handle:         handle,
		RetryDuration:  DefaultTaskRetryTimeout,
		RepeatDuration: DefaultTaskRescheduleTimeout,
		Ctx:            g.Context(),
	})
}
func (g *SharedPtr) AddTaskFuncAsConstruct(handle func() error, descriptions ...string) error {
	return g.AddTaskFunc(TaskType{Construct: true}, handle, descriptions...)
}
func (g *SharedPtr) AddTaskFuncAsConstructRepeat(handle func() error, descriptions ...string) error {
	return g.AddTaskFunc(TaskType{Construct: true, Repeat: true}, handle, descriptions...)
}
func (g *SharedPtr) AddTaskFuncAsRepeat(handle func() error, descriptions ...string) error {
	return g.AddTaskFunc(TaskType{Repeat: true}, handle, descriptions...)
}
func (g *SharedPtr) AddTaskFuncAsRetry(handle func() error, descriptions ...string) error {
	return g.AddTaskFunc(TaskType{Retry: true}, handle, descriptions...)
}

func (g *SharedPtr) RemoveTaskById(id string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.sp.RemoveTaskById(id, true)
	return
}

func (g *SharedPtr) RemoveAllTask() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.sp.RemoveAllTasks()
}

func (g *SharedPtr) TaskIds() []string {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.sp.TaskIds()
}

func (g *SharedPtr) Watch() chan<- Event {
	eventC := g.event()

	go func() {
		swapped := g.WatchStopped.CAS(false, true)
		if !swapped {
			return
		}

		defer func() {
			g.WatchStopped.Store(false)
		}()
	L:
		for {
			select {
			case <-g.Context().Done():
				break L
			case event, ok := <-eventC:
				if !ok {
					break L
				}
				switch event {
				case EventNew, EventExpired:
					if event == EventExpired {
						g.resetPtr()
					}
					// New x
					_, err := g.GetWithRetry()
					if err != nil {
						go g.logger().WithField("event", event).WithError(err).Warn("handle event failed...")
						continue
					}
					go g.logger().WithField("event", event).Infof("handle event success...")
				case EventClose:
					g.resetPtr()
				}
			}
		}
	}()
	return eventC
}

func (g *SharedPtr) Ready() error {
	if g == nil {
		return ErrEmptyValue
	}

	if g.InShutdown() {
		return ErrAlreadyShutdown
	}
	x := g.Get()
	if x != nil {
		return x.Ready()
	}
	return ErrEmptyValue
}

// std::shared_ptr.get() until ptr is ready & std::shared_ptr.make_unique() if necessary
func (g *SharedPtr) GetUntilReady() (Ptr, error) {
	err := Retry(g.Context(), g.logger(), g.sp.TaskMaxRetryDuration, g.sp.ConstructTimeout, func() error {
		x := g.Get()
		if x != nil {
			// check  if x is ready
			if err := x.Ready(); err != nil {
				// until ready
				return err
			}
			return nil
		}

		// New x
		x, err := g.allocate()
		if err != nil {
			return err
		}
		if x == nil {
			return ErrEmptyValue
		}
		return ErrNotReady
	})
	return g.Get(), err
}

// std::shared_ptr.get() & std::shared_ptr.make_unique() if necessary
func (g *SharedPtr) GetWithRetry() (Ptr, error) {
	// if allocated, return now
	if x := g.Get(); x != nil {
		return x, nil
	}

	// New x
	err := Retry(g.Context(), g.logger(), g.sp.TaskMaxRetryDuration, g.sp.ConstructTimeout, func() error {
		_, err := g.allocate()
		return err
	})
	if err != nil {
		return nil, err
	}

	return g.Get(), nil
}

// std::shared_ptr.release()
func (g *SharedPtr) Release() Ptr {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.sp.Release()
}

// std::shared_ptr.reset()
func (g *SharedPtr) Reset() {
	g.RemoveAllTask()
	g.resetPtr()
}

// std::shared_ptr.get()
func (g *SharedPtr) Get() Ptr {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.sp.Get()
}

// reset ptr and ready to start again
func (g *SharedPtr) resetPtr() {
	x := g.Release()
	if x != nil {
		x.Close()
	}
	return
}

func (g *SharedPtr) allocate() (Ptr, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	x := g.sp.Get()
	if x != nil {
		return x, nil
	}
	if g.sp.New == nil {
		return nil, nil
	}
	x, err := g.sp.New()
	if err != nil {
		return g.sp.Get(), err
	}
	g.sp.Reset(x)
	go g.backgroundTask()
	g.recoveryTaskLocked()
	return g.sp.Get(), nil
}

func (g *SharedPtr) recoveryTaskLocked() {
	for _, t := range g.sp.SnapTasks() {
		if t == nil {
			continue
		}

		task := t.Clone()
		select {
		case <-g.sp.Context().Done():
			return
		case <-task.Context().Done():
			continue
		default:
		}

		if !task.Type.Drop {
			task.State = TaskStateDormancy
		}

		if task.State == TaskStateDormancy {
			task.State = TaskStateNew
			g.logger().WithField("task", task.String()).Info("recover task is adding...")
			go func() {
				g.sp.GetTaskC() <- task
			}()
		}
	}
}

func (g *SharedPtr) backgroundTask() {
	swapped := g.backgroundStopped.CAS(false, true)
	if !swapped {
		return
	}
	defer func() {
		g.backgroundStopped.Store(false)
	}()
L:
	for {
		select {
		case <-g.Context().Done():
			break L
		case task, ok := <-g.getTaskC():
			if !ok {
				break L
			}
			if task == nil {
				continue
			}
			if task.State != TaskStateNew {
				go g.logger().WithField("task", task.String()).
					Warn("task is received with unexpected state, ignore duplicate schedule...")
				continue
			}
			// verify whether task is duplicated
			// store task
			_, addedTask := g.GetTaskById(task.ID())

			go g.logger().WithField("task", task.String()).
				Infof("task is received, try to schedule...")

			if addedTask {
				go g.logger().WithField("task", task.String()).
					Warn("task is added already, ignore duplicate schedule...")
				continue
			}

			if task.Type.Construct {
				if _, err := g.GetWithRetry(); err != nil {
					go g.logger().WithField("task", task.String()).
						Warn("task is added but not scheduled, new has not been called yet...")
					continue
				}
			} else {
				if _, err := g.GetUntilReady(); err != nil {
					go g.logger().WithField("task", task.String()).
						Warn("task is added but not scheduled, not ready yet...")
					continue
				}
			}

			addTask := func() {
				g.mu.Lock()
				defer g.mu.Unlock()
				g.sp.AddTask(task)
			}

			deleteTask := func(cancel bool) {
				g.mu.Lock()
				defer g.mu.Unlock()
				g.sp.RemoveTaskById(task.ID(), cancel)
			}

			// Handle task
			addTask()
			go func() {
				if task.State != TaskStateNew {
					return
				}
				task.State = TaskStateRunning
				go g.logger().WithField("task", task.String()).Info("task is running now...")

				// execute the task and refresh the state
				func() {
					defer func() {
						if r := recover(); r != nil {
							task.State = TaskStateDoneErrorHappened
							go g.logger().WithField("task", task.String()).WithField("recovery", r).
								Error("task is done failed...")
						}
					}()
					if task.Handle == nil {
						task.State = TaskStateDoneNormally
						return
					}
					if err := task.Handle(); err != nil {
						task.State = TaskStateDoneErrorHappened
						go g.logger().WithField("task", task.String()).WithError(err).
							Warnf("task is done failed...")
						return
					}
					go g.logger().WithField("task", task.String()).
						Info("task is done successfully...")
					task.State = TaskStateDoneNormally
				}()

				// handle completed execution and refresh the state
				func() {
					waitBeforeRepeat := func() {
						go g.logger().WithField("task", task.String()).
							Warnf("task is rescheduled to repeat in %s...", task.RepeatDuration)
						<-time.After(task.RepeatDuration)
					}
					waitBeforeRetry := func() {
						go g.logger().WithField("task", task.String()).
							Warnf("task is rescheduled to recover in %s...", task.RetryDuration)
						<-time.After(task.RetryDuration)
					}
					select {
					case <-task.Context().Done():
						task.State = TaskStateDeath
						return
					default:

						// Drop
						if task.Type.Drop && !task.Type.Retry {
							task.State = TaskStateDeath
							return
						}

						if task.Type.Drop && task.Type.Retry {
							if task.State == TaskStateDoneErrorHappened {
								waitBeforeRetry()
								task.State = TaskStateNew
								return
							}
							task.State = TaskStateDeath
							return
						}

						// Repeat
						if task.Type.Repeat && !task.Type.Construct {
							waitBeforeRepeat()
							task.State = TaskStateNew
							return
						}

						if task.Type.Repeat && task.Type.Construct {
							go g.logger().WithField("task", task.String()).
								Warnf("task is rescheduled and restart all tasks...")
							deleteTask(false) // don't recover this task, this task will be added later
							g.resetPtr()
							go func() {
								_, _ = g.GetWithRetry()
							}()
							waitBeforeRepeat()
							task.State = TaskStateNew
							return
						}

						// Construct && !Drop && !Repeat
						if task.Type.Construct {
							//Retry
							if task.Type.Retry && task.State == TaskStateDoneErrorHappened {
								go g.logger().WithField("task", task.String()).
									Warnf("task is rescheduled and restart all tasks...")
								deleteTask(false) // don't recover this task, this task will be added later
								g.resetPtr()
								go func() {
									_, _ = g.GetWithRetry()
								}()
								waitBeforeRetry()
								task.State = TaskStateNew
								return
							}
							task.State = TaskStateDormancy
							return
						}

						// !Dop && !Repeat && !Construct
						if task.Type.Retry {
							if task.State == TaskStateDoneErrorHappened {
								go g.logger().WithField("task", task.String()).
									Warnf("task is rescheduled ...")
								deleteTask(false) // don't recover this task, this task will be added later
								waitBeforeRetry()
								task.State = TaskStateNew
								return
							}
							task.State = TaskStateDormancy
							return
						}
						task.State = TaskStateDormancy
						return
					}

				}()
				// complete the task's life cycle
				func() {
					select {
					case <-task.Context().Done():
						go g.logger().WithField("task", task.String()).
							Info("task is canceled, go to death now...")
						deleteTask(false) // canceled already
						return
					default:
						switch task.State {
						case TaskStateNew:
							deleteTask(false)
							if task.State == TaskStateRunning {
								go g.logger().WithField("task", task.String()).
									Infof("task is rescheduled now...")
							}
							go g.logger().WithField("task", task.String()).
								Infof("task is rescheduled now...")
							g.getTaskC() <- task
						case TaskStateDormancy:
							go g.logger().WithField("task", task.String()).
								Infof("task is done,  go to dormancy...")
						case TaskStateDeath:
							go g.logger().WithField("task", task.String()).
								Infof("task is dead,  go to death...")
							deleteTask(true)
						default:
							go g.logger().WithField("task", task.String()).
								Info("task is with unexpect state, go to death now...")
							deleteTask(true)
						}
					}
				}()

			}()
		}
	}
}

func (g *SharedPtr) getTaskC() chan *Task {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.sp.GetTaskC()
}
func (g *SharedPtr) event() chan Event {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.sp.Event()
}
