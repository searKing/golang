// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync

import (
	"context"
	"sync"

	"github.com/searKing/golang/go/time/rate"
)

type Walk struct {
	Burst int // Burst will be set to 1 if less than 1

	wg sync.WaitGroup

	walkError     error
	walkErrorOnce sync.Once
	walkErrorMu   sync.Mutex
}

// WalkFunc is the type of the function called for each task processed
// by Walk. The path argument contains the argument to Walk as a task.
//
// In the case of an error, the info argument will be nil. If an error
// is returned, processing stops.
type WalkFunc func(task any) error

// Walk will consume all tasks parallel and block until ctx.Done() or taskChan is closed.
// Walk returns a channel that's closed when work done on behalf of this
// walk should be canceled. Done may return nil if this walk can
// never be canceled. Successive calls to Done return the same value.
// The close of the Done channel may happen asynchronously,
// after the cancel function returns.
func (p *Walk) Walk(ctx context.Context, taskChan <-chan any, procFn WalkFunc) (doneC <-chan struct{}) {
	done := make(chan struct{})
	if p.Burst <= 0 {
		p.Burst = 1
	}
	p.wg.Add(1)

	go func() {
		defer close(done)
		defer p.wg.Done()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		limiter := rate.NewFullBurstLimiter(p.Burst)
		for {
			// Wait blocks until lim permits 1 events to happen.
			if err := limiter.Wait(ctx); err != nil {
				p.TrySetError(err)
				return
			}

			select {
			case task, ok := <-taskChan:
				if !ok {
					return
				}
				p.wg.Add(1)
				go func() {
					defer limiter.PutToken()
					defer p.wg.Done()
					if err := procFn(task); err != nil {
						p.TrySetError(err)
						return
					}
				}()
			case <-ctx.Done():
				p.TrySetError(ctx.Err())
				return
			}
		}
	}()
	return done
}

// Wait blocks until the WaitGroup counter is zero.
func (p *Walk) Wait() error {
	p.wg.Wait()
	return p.Error()
}

func (p *Walk) Error() error {
	p.walkErrorMu.Lock()
	defer p.walkErrorMu.Unlock()
	return p.walkError
}
func (p *Walk) TrySetError(err error) {
	p.walkErrorOnce.Do(func() {
		p.walkErrorMu.Lock()
		defer p.walkErrorMu.Unlock()
		p.walkError = err
	})
}
