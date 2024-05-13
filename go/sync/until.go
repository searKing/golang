// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// ErrUntilClosed is returned by the Until's Do method after a call to Close.
var ErrUntilClosed = errors.New("sync: Until closed")

// Until represents a class of work and forms a namespace in
// which units of work can be executed with duplicate suppression.
// It blocks on certain Do actions and unblock when Retry is called.
type Until struct {
	mu         sync.Mutex
	done       bool
	blockingCh chan struct{}
}

// Retry unblocks all blocked pick.
func (u *Until) Retry() {
	u.mu.Lock()
	if u.done {
		u.mu.Unlock()
		return
	}
	if u.blockingCh != nil {
		close(u.blockingCh)
	}
	u.blockingCh = make(chan struct{})
	u.mu.Unlock()
}

// Do executes and returns the results of the given function.
// It may block in the following cases:
// - the current fn is nil
// - the err returned by the current fn is not nil
// When one of these situations happens, Do blocks until the Retry is called.
func (u *Until) Do(ctx context.Context, fn func() (any, error)) (val any, err error) {
	u.mu.Lock()
	if u.blockingCh == nil {
		u.blockingCh = make(chan struct{})
	}
	u.mu.Unlock()
	var ch chan struct{}

	var lastPickErr error
	for {
		u.mu.Lock()
		if u.done {
			u.mu.Unlock()
			return nil, ErrUntilClosed
		}
		if fn == nil {
			ch = u.blockingCh
		}
		if ch == u.blockingCh {
			// This could happen when either:
			// - forget (the previous if condition), or
			// - has called Do on the current Do.
			u.mu.Unlock()
			select {
			case <-ctx.Done():
				var errStr string
				if lastPickErr != nil {
					errStr = "latest single_flight_group error: " + lastPickErr.Error()
				} else {
					errStr = ctx.Err().Error()
				}
				return nil, fmt.Errorf("%s: %w", errStr, ctx.Err())
			case <-ch:
			}
			continue
		}

		ch = u.blockingCh
		u.mu.Unlock()
		val, err = fn()
		if err != nil {
			lastPickErr = err
			// continue back to the beginning of the for loop to redo.
			continue
		}
		return val, nil
	}
}

func (u *Until) Close() {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.done {
		return
	}
	u.done = true
	close(u.blockingCh)
}
