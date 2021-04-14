// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync

import (
	"context"
	"errors"
	"sync"
	"time"
)

type ConditionVariable struct {
	subject Subject

	// L is held while observing or changing the condition
	L sync.Locker
}

// NewConditionVariable returns a new ConditionVariable with Locker l.
// see http://www.cplusplus.com/reference/condition_variable/condition_variable/
func NewConditionVariable(l sync.Locker) *ConditionVariable {
	return &ConditionVariable{L: l}
}

// Wait atomically unlocks c.L and suspends execution
// of the calling goroutine. After later resuming execution,
// Wait locks c.L before returning. Unlike in other systems,
// Wait cannot return unless awoken by Broadcast or Signal.
//
// Because c.L is not locked when Wait first resumes, the caller
// typically cannot assume that the condition is true when
// Wait returns. Instead, the caller should Wait in a loop:
//
//    c.L.Lock()
//    for !condition() {
//        c.Wait()
//    }
//    ... make use of condition ...
//    c.L.Unlock()
//
// Wait wait until notified
func (c *ConditionVariable) Wait() {
	eventC, cancel := c.subject.Subscribe()
	defer cancel()
	c.L.Unlock()
	defer c.L.Lock()
	<-eventC
}

// WaitPred wait until notified
// If pred is specified (2), the function only blocks if pred returns false,
// and notifications can only unblock the thread when it becomes true
// (which is specially useful to check against spurious wake-up calls). This version (2) behaves as if implemented as:
// while (!pred()) wait(lck);
func (c *ConditionVariable) WaitPred(pred func() bool) {
	if pred == nil {
		pred = func() bool { return false }
	}
	for !pred() {
		c.Wait()
	}
}

// WaitFor The execution of the current goroutine (which shall have locked lck's mutex) is blocked during rel_time,
// or until notified (if the latter happens first).
// WaitFor wait for timeout or until notified
// It behaves as if implemented as:
// return wait_until (lck, chrono::steady_clock::now() + rel_time);
func (c *ConditionVariable) WaitFor(timeout time.Duration) error {
	return c.WaitUntil(time.Now().Add(timeout))
}

// WaitForPred wait for timeout or until notified
// If pred is nil, do as pred returns false always,
// If pred is specified, the function only blocks if pred returns false,
// and notifications can only unblock the thread when it becomes true
// (which is especially useful to check against spurious wake-up calls).
// It behaves as if implemented as:
// return wait_until (lck, chrono::steady_clock::now() + rel_time, std::move(pred));
func (c *ConditionVariable) WaitForPred(timeout time.Duration, pred func() bool) bool {
	return c.WaitUntilPred(time.Now().Add(timeout), pred)
}

// WaitUntil wait until notified or time point
// The execution of the current thread (which shall have locked lck's mutex) is blocked either
// until notified or until abs_time, whichever happens first.
func (c *ConditionVariable) WaitUntil(d time.Time) error {
	eventC, cancel := c.subject.Subscribe()
	defer cancel()

	ctx, cancelCtx := context.WithDeadline(context.Background(), d)
	defer cancelCtx()

	select {
	case <-eventC:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// WaitUntilPred wait until notified or time point
// The execution of the current thread (which shall have locked c.L's mutex) is blocked either
// until notified or until abs_time, whichever happens first.
// If pred is nil, do as pred returns false always,
// If pred is specified, the function only blocks if pred returns false,
// and notifications can only unblock the thread when it becomes true
// (which is especially useful to check against spurious wake-up calls).
// It behaves as if implemented as:
// while (!pred())
//  if ( wait_until(lck,abs_time) == cv_status::timeout)
//    return pred();
// return true;
func (c *ConditionVariable) WaitUntilPred(d time.Time, pred func() bool) bool {
	if pred == nil {
		pred = func() bool { return false }
	}
	for !pred() {
		if errors.Is(c.WaitUntil(d), context.DeadlineExceeded) {
			return pred()
		}
	}
	return true
}

// Signal wakes one goroutine waiting on c, if there is any.
//
// It is allowed but not required for the caller to hold c.L
// during the call.
func (c *ConditionVariable) Signal() {
	go func() {
		_ = c.subject.PublishSignal(context.Background(), nil)
	}()
}

// Broadcast wakes all goroutines waiting on c.
//
// It is allowed but not required for the caller to hold c.L
// during the call.
func (c *ConditionVariable) Broadcast() {
	go func() {
		_ = c.subject.PublishBroadcast(context.Background(), nil)
	}()
}
