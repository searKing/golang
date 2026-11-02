// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package context_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"testing/synctest"
	"time"

	context_ "github.com/searKing/golang/go/context"
)

// afterFuncContext is a context that is not a *context.cancelCtx and
// implements AfterFunc, so that registrations made by context.AfterFunc
// on it can be observed.
type afterFuncContext struct {
	mu         sync.Mutex
	afterFuncs map[*byte]func()
	done       chan struct{}
	err        error
}

func newAfterFuncContext() *afterFuncContext {
	return &afterFuncContext{
		afterFuncs: make(map[*byte]func()),
		done:       make(chan struct{}),
	}
}

func (c *afterFuncContext) Deadline() (time.Time, bool) { return time.Time{}, false }
func (c *afterFuncContext) Done() <-chan struct{}       { return c.done }
func (c *afterFuncContext) Value(key any) any           { return nil }

func (c *afterFuncContext) Err() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.err
}

func (c *afterFuncContext) AfterFunc(f func()) func() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	k := new(byte)
	c.afterFuncs[k] = f
	return func() bool {
		c.mu.Lock()
		defer c.mu.Unlock()
		_, ok := c.afterFuncs[k]
		delete(c.afterFuncs, k)
		return ok
	}
}

func (c *afterFuncContext) cancel(err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.err != nil {
		return
	}
	c.err = err
	for _, f := range c.afterFuncs {
		go f()
	}
	clear(c.afterFuncs)
	close(c.done)
}

func (c *afterFuncContext) numAfterFuncs() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.afterFuncs)
}

func checkCanceled(t *testing.T, ctx context.Context, wantCause error) {
	t.Helper()
	select {
	case <-ctx.Done():
	default:
		t.Fatal("ctx not canceled")
	}
	if got, want := ctx.Err(), context.Canceled; got != want {
		t.Errorf("ctx.Err() = %v, want %v", got, want)
	}
	if got := context.Cause(ctx); got != wantCause {
		t.Errorf("context.Cause(ctx) = %v, want %v", got, wantCause)
	}
}

func checkNotCanceled(t *testing.T, ctx context.Context) {
	t.Helper()
	select {
	case <-ctx.Done():
		t.Fatalf("ctx canceled unexpectedly: %v", context.Cause(ctx))
	default:
	}
}

func TestWithCancelCauseContexts(t *testing.T) {
	errParent := errors.New("parent")
	errCause := errors.New("cause")
	errCancel := errors.New("cancel")

	tests := []struct {
		name string
		// run creates parent and causes, calls WithCancelCauseContexts,
		// triggers cancellation and returns the expected cause.
		run func(t *testing.T) (ctx context.Context, cancel context.CancelCauseFunc, want error)
	}{{
		name: "cancel",
		run: func(t *testing.T) (context.Context, context.CancelCauseFunc, error) {
			cause, stop := context.WithCancelCause(context.Background())
			t.Cleanup(func() { stop(nil) })
			ctx, cancel := context_.WithCancelCauseContexts(context.Background(), cause)
			cancel(errCancel)
			return ctx, cancel, errCancel
		},
	}, {
		name: "parent",
		run: func(t *testing.T) (context.Context, context.CancelCauseFunc, error) {
			parent, stopParent := context.WithCancelCause(context.Background())
			cause, stop := context.WithCancelCause(context.Background())
			t.Cleanup(func() { stop(nil) })
			ctx, cancel := context_.WithCancelCauseContexts(parent, cause)
			stopParent(errParent)
			return ctx, cancel, errParent
		},
	}, {
		name: "cause",
		run: func(t *testing.T) (context.Context, context.CancelCauseFunc, error) {
			cause1, stop1 := context.WithCancelCause(context.Background())
			t.Cleanup(func() { stop1(nil) })
			cause2, stop2 := context.WithCancelCause(context.Background())
			ctx, cancel := context_.WithCancelCauseContexts(context.Background(), cause1, cause2)
			stop2(errCause)
			return ctx, cancel, errCause
		},
	}, {
		name: "cause without cause",
		run: func(t *testing.T) (context.Context, context.CancelCauseFunc, error) {
			cause, stop := context.WithCancel(context.Background())
			ctx, cancel := context_.WithCancelCauseContexts(context.Background(), cause)
			stop()
			return ctx, cancel, context.Canceled
		},
	}, {
		name: "first cancellation wins",
		run: func(t *testing.T) (context.Context, context.CancelCauseFunc, error) {
			cause, stop := context.WithCancelCause(context.Background())
			ctx, cancel := context_.WithCancelCauseContexts(context.Background(), cause)
			cancel(errCancel)
			synctest.Wait()
			stop(errCause)
			return ctx, cancel, errCancel
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				ctx, cancel, want := tt.run(t)
				defer cancel(nil)
				synctest.Wait()
				checkCanceled(t, ctx, want)
			})
		})
	}
}

func TestWithCancelCauseContextsNotCanceled(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		cause, stop := context.WithCancelCause(context.Background())
		defer stop(nil)
		ctx, cancel := context_.WithCancelCauseContexts(context.Background(), nil, context.Background(), cause)
		defer cancel(nil)
		synctest.Wait()
		checkNotCanceled(t, ctx)
	})
}

func TestWithCancelCauseContextsNoCauses(t *testing.T) {
	errCancel := errors.New("cancel")
	ctx, cancel := context_.WithCancelCauseContexts(context.Background())
	checkNotCanceled(t, ctx)
	cancel(errCancel)
	checkCanceled(t, ctx, errCancel)
}

func TestWithCancelCauseContextsCauseAlreadyDone(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		errCause := errors.New("cause")
		cause, stop := context.WithCancelCause(context.Background())
		stop(errCause)

		other := newAfterFuncContext()
		ctx, cancel := context_.WithCancelCauseContexts(context.Background(), other, cause)
		defer cancel(nil)

		// Cancellation must be visible synchronously, as with
		// context.WithCancelCause on an already canceled parent.
		checkCanceled(t, ctx, errCause)

		synctest.Wait()
		if got := other.numAfterFuncs(); got != 0 {
			t.Errorf("cause has %d registered AfterFuncs, want 0", got)
		}
	})
}

func TestWithCancelCauseContextsParentAlreadyDone(t *testing.T) {
	errParent := errors.New("parent")
	parent, stop := context.WithCancelCause(context.Background())
	stop(errParent)

	ctx, cancel := context_.WithCancelCauseContexts(parent, context.Background())
	defer cancel(nil)
	checkCanceled(t, ctx, errParent)
}

// TestWithCancelCauseContextsReleasesCauses checks that resources registered
// with causes are released however the returned context is canceled.
func TestWithCancelCauseContextsReleasesCauses(t *testing.T) {
	tests := []struct {
		name   string
		cancel func(parent, cause *afterFuncContext, cancel context.CancelCauseFunc)
	}{{
		name:   "cancel",
		cancel: func(_, _ *afterFuncContext, cancel context.CancelCauseFunc) { cancel(nil) },
	}, {
		name:   "parent",
		cancel: func(parent, _ *afterFuncContext, _ context.CancelCauseFunc) { parent.cancel(context.Canceled) },
	}, {
		name:   "cause",
		cancel: func(_, cause *afterFuncContext, _ context.CancelCauseFunc) { cause.cancel(context.Canceled) },
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				parent := newAfterFuncContext()
				causes := []*afterFuncContext{newAfterFuncContext(), newAfterFuncContext()}
				ctx, cancel := context_.WithCancelCauseContexts(parent, causes[0], causes[1])
				defer cancel(nil)
				for i, c := range causes {
					if got := c.numAfterFuncs(); got != 1 {
						t.Fatalf("causes[%d] has %d registered AfterFuncs, want 1", i, got)
					}
				}

				tt.cancel(parent, causes[0], cancel)
				synctest.Wait()

				if ctx.Err() == nil {
					t.Fatal("ctx not canceled")
				}
				for i, c := range causes {
					if got := c.numAfterFuncs(); got != 0 {
						t.Errorf("causes[%d] has %d registered AfterFuncs, want 0", i, got)
					}
				}
			})
		})
	}
}

func TestWithCancelCauseContextsConcurrent(t *testing.T) {
	const n = 16
	causes := make([]context.Context, n)
	stops := make([]context.CancelCauseFunc, n)
	for i := range n {
		causes[i], stops[i] = context.WithCancelCause(context.Background())
	}
	ctx, cancel := context_.WithCancelCauseContexts(context.Background(), causes...)
	defer cancel(nil)

	var wg sync.WaitGroup
	for _, stop := range stops {
		wg.Go(func() { stop(errors.New("cause")) })
	}
	wg.Go(func() { cancel(errors.New("cancel")) })
	wg.Wait()

	<-ctx.Done()
	if context.Cause(ctx) == nil {
		t.Error("context.Cause(ctx) = nil, want non-nil")
	}
}
