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

// ConditionVariable is an object able to block the calling thread until notified to resume.
// see http://www.cplusplus.com/reference/condition_variable/condition_variable/
type ConditionVariable struct {
	subject Subject
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
//	c.L.Lock()
//	for !condition() {
//	    c.Wait()
//	}
//	... make use of condition ...
//	c.L.Unlock()
//
// Wait wait until notified
func (c *ConditionVariable) Wait(lck sync.Locker) {
	_ = c.WaitContext(context.Background(), lck)
}

// WaitPred wait until notified
// If pred is specified (2), the function only blocks if pred returns false,
// and notifications can only unblock the thread when it becomes true
// (which is specially useful to check against spurious wake-up calls). This version (2) behaves as if implemented as:
// while (!pred()) wait(lck);
func (c *ConditionVariable) WaitPred(lck sync.Locker, pred func() bool) {
	_ = c.WaitPredContext(context.Background(), lck, pred)
}

// WaitFor The execution of the current goroutine (which shall have locked lck's mutex) is blocked during rel_time,
// or until notified (if the latter happens first).
// WaitFor wait for timeout or until notified
// It behaves as if implemented as:
// return wait_until (lck, chrono::steady_clock::now() + rel_time);
func (c *ConditionVariable) WaitFor(lck sync.Locker, timeout time.Duration) error {
	return c.WaitUntil(lck, time.Now().Add(timeout))
}

// WaitForPred wait for timeout or until notified
// If pred is nil, do as pred returns false always,
// If pred is specified, the function only blocks if pred returns false,
// and notifications can only unblock the thread when it becomes true
// (which is especially useful to check against spurious wake-up calls).
// It behaves as if implemented as:
// return wait_until (lck, chrono::steady_clock::now() + rel_time, std::move(pred));
func (c *ConditionVariable) WaitForPred(lck sync.Locker, timeout time.Duration, pred func() bool) bool {
	return c.WaitUntilPred(lck, time.Now().Add(timeout), pred)
}

// WaitUntil wait until notified or time point
// The execution of the current thread (which shall have locked lck's mutex) is blocked either
// until notified or until abs_time, whichever happens first.
func (c *ConditionVariable) WaitUntil(lck sync.Locker, d time.Time) error {
	ctx, cancelCtx := context.WithDeadline(context.Background(), d)
	defer cancelCtx()
	return c.WaitContext(ctx, lck)
}

// WaitUntilPred wait until notified or ctx done
// The execution of the current thread (which shall have locked lck's mutex) is blocked either
// until notified or until ctx, whichever happens first.
// If pred is nil, do as pred returns false always,
// If pred is specified, the function only blocks if pred returns false,
// and notifications can only unblock the thread when it becomes true
// (which is especially useful to check against spurious wake-up calls).
// It behaves as if implemented as:
// while (!pred())
//
//	if ( wait_until(lck,abs_time) == cv_status::timeout)
//	  return pred();
//
// return true;
func (c *ConditionVariable) WaitUntilPred(lck sync.Locker, d time.Time, pred func() bool) bool {
	ctx, cancelCtx := context.WithDeadline(context.Background(), d)
	defer cancelCtx()
	return c.WaitPredContext(ctx, lck, pred)
}

// WaitContext wait until notified or time point
// The execution of the current thread (which shall have locked lck's mutex) is blocked either
// until notified or until ctx done, whichever happens first.
func (c *ConditionVariable) WaitContext(ctx context.Context, lck sync.Locker) error {
	eventC, cancel := c.subject.Subscribe()
	defer cancel()
	lck.Unlock()
	defer lck.Lock()

	select {
	case <-eventC:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// WaitPredContext wait until notified or ctx done
// The execution of the current thread (which shall have locked lck's mutex) is blocked either
// until notified or until ctx, whichever happens first.
// If pred is nil, do as pred returns false always,
// If pred is specified, the function only blocks if pred returns false,
// and notifications can only unblock the thread when it becomes true
// (which is especially useful to check against spurious wake-up calls).
// It behaves as if implemented as:
// while (!pred())
//
//	if ( wait_until(ctx,lck) == cv_status::timeout)
//	  return pred();
//
// return true;
func (c *ConditionVariable) WaitPredContext(ctx context.Context, lck sync.Locker, pred func() bool) bool {
	if pred == nil {
		pred = func() bool { return false }
	}
	for !pred() {
		if errors.Is(c.WaitContext(ctx, lck), context.DeadlineExceeded) {
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
