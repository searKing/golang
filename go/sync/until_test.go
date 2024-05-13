// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync_test

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	sync_ "github.com/searKing/golang/go/sync"
)

const goroutineCount = 5

func TestUntil_DoTimeout(t *testing.T) {
	var until sync_.Until
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	if _, err := until.Do(ctx, func() (any, error) { return "bar", fmt.Errorf("must error") }); !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("until.Do returned error %v, want DeadlineExceeded", err)
	}
}

func TestUntil_DoBlocking(t *testing.T) {
	var until sync_.Until
	// All goroutines should block because picker is nil in until.
	var finishedCount uint64
	for i := goroutineCount; i > 0; i-- {
		go func() {
			if _, err := until.Do(context.Background(), nil); err == nil {
				t.Errorf("until.Do returned nil error")
			}
			atomic.AddUint64(&finishedCount, 1)
		}()
	}
	time.Sleep(50 * time.Millisecond)
	if c := atomic.LoadUint64(&finishedCount); c != 0 {
		t.Errorf("finished goroutines count: %v, want 0", c)
	}
	until.Close()
}

func TestUntil_Do(t *testing.T) {
	var until sync_.Until
	// All goroutines should block because picker is nil in until.
	var finishedCount uint64
	for i := goroutineCount; i > 0; i-- {
		go func() {
			if tr, err := until.Do(context.Background(), func() (any, error) {
				return "bar", nil
			}); err != nil || tr != "bar" {
				t.Errorf("until.Do returned non-nil error: %v", err)
			}
			atomic.AddUint64(&finishedCount, 1)
		}()
	}
	time.Sleep(50 * time.Millisecond)
	if c := atomic.LoadUint64(&finishedCount); c != goroutineCount {
		t.Errorf("finished goroutines count: %v, want 0", c)
	}
	until.Close()
}

func TestUntil_DoBlockingAndGot(t *testing.T) {
	var mu sync.Mutex
	var fn = func() (any, error) {
		return "bar", fmt.Errorf("must error")
	}

	var fnHolder = func() (any, error) {
		mu.Lock()
		defer mu.Unlock()
		return fn()
	}

	var until sync_.Until
	// All goroutines should block because picker is nil in until.
	var finishedCount uint64
	for i := goroutineCount; i > 0; i-- {
		go func() {
			if tr, err := until.Do(context.Background(), fnHolder); err != nil || tr != "bar" {
				t.Errorf("until.Do returned non-nil error: %v", err)
			}
			atomic.AddUint64(&finishedCount, 1)
		}()
	}
	time.Sleep(50 * time.Millisecond)
	if c := atomic.LoadUint64(&finishedCount); c != 0 {
		t.Errorf("finished goroutines count: %v, want 0", c)
	}
	fn = func() (any, error) {
		return "bar", nil
	}
	until.Retry()

	time.Sleep(50 * time.Millisecond)
	if c := atomic.LoadUint64(&finishedCount); c != goroutineCount {
		t.Errorf("finished goroutines count: %v, want 0", c)
	}
	until.Close()
}
