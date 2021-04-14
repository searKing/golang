// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync

import (
	"context"
	"sync"
	"time"
)

type TimeoutCond struct {
	cond ConditionVariable

	// L is held while observing or changing the condition
	L sync.Locker
}

// NewTimeoutCond returns a new TimeoutCond with Locker l.
func NewTimeoutCond(l sync.Locker) *TimeoutCond {
	return &TimeoutCond{L: l}
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
func (c *TimeoutCond) Wait() {
	c.cond.Wait(c.L)
}

// WaitPred wait until notified
// If pred is specified (2), the function only blocks if pred returns false,
// and notifications can only unblock the thread when it becomes true
// (which is specially useful to check against spurious wake-up calls). This version (2) behaves as if implemented as:
// while (!pred()) wait(lck);
func (c *TimeoutCond) WaitPred(pred func() bool) {
	c.cond.WaitPred(c.L, pred)
}

// WaitFor The execution of the current goroutine (which shall have locked lck's mutex) is blocked during rel_time,
// or until notified (if the latter happens first).
// WaitFor wait for timeout or until notified
// It behaves as if implemented as:
// return wait_until (lck, chrono::steady_clock::now() + rel_time);
func (c *TimeoutCond) WaitFor(timeout time.Duration) error {
	return c.cond.WaitFor(c.L, timeout)
}

// WaitForPred wait for timeout or until notified
// If pred is nil, do as pred returns false always,
// If pred is specified, the function only blocks if pred returns false,
// and notifications can only unblock the thread when it becomes true
// (which is especially useful to check against spurious wake-up calls).
// It behaves as if implemented as:
// return wait_until (lck, chrono::steady_clock::now() + rel_time, std::move(pred));
func (c *TimeoutCond) WaitForPred(timeout time.Duration, pred func() bool) bool {
	return c.cond.WaitForPred(c.L, timeout, pred)
}

// WaitUntil wait until notified or time point
// The execution of the current thread (which shall have locked lck's mutex) is blocked either
// until notified or until abs_time, whichever happens first.
func (c *TimeoutCond) WaitUntil(d time.Time) error {
	return c.cond.WaitUntil(c.L, d)
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
func (c *TimeoutCond) WaitUntilPred(d time.Time, pred func() bool) bool {
	return c.cond.WaitUntilPred(c.L, d, pred)
}

// WaitContext wait until notified or time point
// The execution of the current thread (which shall have locked lck's mutex) is blocked either
// until notified or until ctx done, whichever happens first.
func (c *TimeoutCond) WaitContext(ctx context.Context) error {
	return c.cond.WaitContext(ctx, c.L)
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
//  if ( wait_until(ctx,lck) == cv_status::timeout)
//    return pred();
// return true;
func (c *TimeoutCond) WaitPredContext(ctx context.Context, pred func() bool) bool {
	return c.cond.WaitPredContext(ctx, c.L, pred)
}

// Signal wakes one goroutine waiting on c, if there is any.
//
// It is allowed but not required for the caller to hold c.L
// during the call.
func (c *TimeoutCond) Signal() {
	c.cond.Signal()
}

// Broadcast wakes all goroutines waiting on c.
//
// It is allowed but not required for the caller to hold c.L
// during the call.
func (c *TimeoutCond) Broadcast() {
	c.cond.Broadcast()
}
