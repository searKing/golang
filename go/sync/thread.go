// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync

import (
	"context"
	"runtime"
	"sync"
)

//go:generate go-option -type=threadDo
type threadDo struct {
	// call the function f in the same thread or escape thread
	EscapeThread bool
}

// Thread should be used for such as calling OS services or
// non-Go library functions that depend on per-thread state, as runtime.LockOSThread().
type Thread struct {
	GoRoutine bool // Use thread as goroutine, that is without runtime.LockOSThread()

	once sync.Once
	// fCh optionally specifies a function to generate
	// a value when Get would otherwise return nil.
	// It may not be changed concurrently with calls to Get.
	fCh chan func()

	mu     sync.Mutex
	ctx    context.Context
	cancel func()
}

func (th *Thread) Shutdown() {
	th.initOnce()
	th.cancel()
}

func (th *Thread) initOnce() {
	th.once.Do(func() {
		th.mu.Lock()
		defer th.mu.Unlock()
		th.ctx, th.cancel = context.WithCancel(context.Background())
		th.fCh = make(chan func())
		go th.lockOSThreadForever()
	})
}

// Do will call the function f in the same thread or escape thread.
// f is enqueued only if ctx is not canceled and Thread is not Shutdown and Not escape
func (th *Thread) Do(ctx context.Context, f func(), opts ...ThreadDoOption) error {
	var opt threadDo
	opt.ApplyOptions(opts...)
	return th.do(ctx, f, opt.EscapeThread)
}

func (th *Thread) do(ctx context.Context, f func(), escapeThread bool) error {
	th.initOnce()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-th.ctx.Done():
		return th.ctx.Err()
	default:
		break
	}

	var r interface{}
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
	case th.fCh <- neverPanic:
		wg.Wait() // wait for f has been executed or panic
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-th.ctx.Done():
		return th.ctx.Err()
	}
}

func (th *Thread) lockOSThreadForever() {
	defer th.cancel()
	if !th.GoRoutine {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
	}
	for {
		select {
		case handler, ok := <-th.fCh:
			if !ok {
				return
			}
			if handler == nil {
				continue
			}
			handler()
		case <-th.ctx.Done():
			return
		}
	}
}
