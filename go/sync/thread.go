// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync

import (
	"context"
	"runtime"
	"sync"
)

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

// Do will call the function f in the same thread
// f is enqueued only if ctx is not canceled and Thread is not Shutdown
func (th *Thread) Do(ctx context.Context, f func()) error {
	th.initOnce()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-th.ctx.Done():
		return th.ctx.Err()
	default:
		break
	}

	var wg sync.WaitGroup
	defer wg.Wait()
	wg.Add(1)
	var r interface{}
	defer func() {
		if r != nil {
			panic(r)
		}
	}()
	monitor := func() {
		defer wg.Done()
		defer func() {
			r = recover()
		}()
		f()
	}
	select {
	case th.fCh <- monitor:
		return nil
	case <-ctx.Done():
		wg.Done()
		return ctx.Err()
	case <-th.ctx.Done():
		wg.Done()
		return th.ctx.Err()
	}
}

func (th *Thread) lockOSThreadForever() {
	if th.GoRoutine {
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
