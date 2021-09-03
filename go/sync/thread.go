// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync

import (
	"context"
	"runtime"
	"sync"
)

// Thread should be used for such as  calling OS services or
// non-Go library functions that depend on per-thread state, as runtime.LockOSThread().
type Thread struct {
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
	th.mu.Lock()
	defer th.mu.Unlock()
	if th.cancel != nil {
		th.cancel()
	}
}

// Do will calls the function f in the same thread
// return true if f is enqueued to call
func (th *Thread) Do(f func()) {
	th.once.Do(func() {
		th.mu.Lock()
		defer th.mu.Unlock()
		th.ctx, th.cancel = context.WithCancel(context.Background())
		th.fCh = make(chan func())
		go th.lockOSThreadForever()
	})

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
	case <-th.ctx.Done():
		wg.Done()
	}
}

func (th *Thread) lockOSThreadForever() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
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
