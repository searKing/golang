// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync

import (
	"context"
	"errors"
	"expvar"
	"runtime"
	"sync"

	expvar_ "github.com/searKing/golang/go/expvar"
)

var threadStatsOnce sync.Once
var threadStats *expvar.Map
var osThreadLeak, goroutineLeak, handlerLeak expvar_.Leak

//go:generate go-option -type=threadDo
type threadDo struct {
	// call the function f in the same thread or escape thread
	EscapeThread bool
}

// Thread should be used for such as calling OS services or
// non-Go library functions that depend on per-thread state, as runtime.LockOSThread().
type Thread struct {
	GoRoutine bool // Use thread as goroutine, that is without runtime.LockOSThread()

	// The Leak is published as a variable directly.
	GoroutineLeak *expvar_.Leak // represents whether goroutine is leaked, take effects if not nil
	OSThreadLeak  *expvar_.Leak // represents whether runtime.LockOSThread is leaked, take effects  if not nil
	HandlerLeak   *expvar_.Leak // represents whether handler in Do is blocked is leaked, take effects  if not nil

	once sync.Once
	// fCh optionally specifies a function to generate
	// a value when Get would otherwise return nil.
	// It may not be changed concurrently with calls to Get.
	fCh chan func()

	mu     sync.Mutex
	ctx    context.Context
	cancel context.CancelCauseFunc
}

// WatchStats bind Leak var to "sync.Thread"
func (t *Thread) WatchStats() {
	threadStatsOnce.Do(func() {
		threadStats = expvar.NewMap("sync.Thread")
		threadStats.Set("goroutine_leak", &goroutineLeak)
		threadStats.Set("os_thread_leak", &osThreadLeak)
		threadStats.Set("handler_leak", &handlerLeak)
	})
	t.GoroutineLeak = &goroutineLeak
	t.OSThreadLeak = &osThreadLeak
	t.HandlerLeak = &handlerLeak
}

// ErrThreadClosed is returned by the Thread's Do methods after a call to `Shutdown`.
var ErrThreadClosed = errors.New("sync: Thread closed")

func (t *Thread) Shutdown() {
	t.initOnce()
	t.cancel(ErrThreadClosed)
}

func (t *Thread) initOnce() {
	t.once.Do(func() {
		t.mu.Lock()
		defer t.mu.Unlock()

		t.ctx, t.cancel = context.WithCancelCause(context.Background())
		t.fCh = make(chan func())
		go t.lockOSThreadForever()
	})
}

// Do will call the function f in the same thread or escape thread.
// f is enqueued only if ctx is not canceled and Thread is not Shutdown and Not escape
func (t *Thread) Do(ctx context.Context, f func(), opts ...ThreadDoOption) error {
	var opt threadDo
	opt.ApplyOptions(opts...)
	return t.do(ctx, f, opt.EscapeThread)
}

func (t *Thread) do(ctx context.Context, f func(), escapeThread bool) error {
	t.initOnce()
	if t.HandlerLeak != nil {
		t.HandlerLeak.Add(1)
		defer t.HandlerLeak.Done()
	}
	select {
	case <-ctx.Done():
		return context.Cause(ctx)
	case <-t.ctx.Done():
		return context.Cause(t.ctx)
	default:
		break
	}

	var r any
	defer func() {
		if r != nil {
			panic(r) // rethrow panic if panic in f
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	neverPanic := func() {
		defer wg.Done() // Mark f is called or panic
		defer func() {
			r = recover()
		}()
		f()
	}

	if escapeThread {
		neverPanic()
		return nil
	}

	select {
	case t.fCh <- neverPanic:
		wg.Wait() // wait for f has been executed or panic
		return nil
	case <-ctx.Done():
		return context.Cause(ctx)
	case <-t.ctx.Done():
		return context.Cause(t.ctx)
	}
}

func (t *Thread) lockOSThreadForever() {
	defer t.cancel(ErrThreadClosed)
	if t.GoroutineLeak != nil {
		t.GoroutineLeak.Add(1)
		defer t.GoroutineLeak.Done()
	}
	if !t.GoRoutine {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		if t.OSThreadLeak != nil {
			t.OSThreadLeak.Add(1)
			defer t.OSThreadLeak.Done()
		}
	}
	for {
		select {
		case handler, ok := <-t.fCh:
			if !ok {
				return
			}
			if handler == nil {
				continue
			}
			handler()
		case <-t.ctx.Done():
			return
		}
	}
}
