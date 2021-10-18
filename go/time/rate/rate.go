// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package rate
// The key observation and some code (shr) is borrowed from
// time/rate/rate.go
package rate

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type expectKeyType struct{}

var expectTokensKey expectKeyType

// A BurstLimiter controls how frequently events are allowed to happen.
// It implements a "token bucket" of size b, initially full and refilled
// by PutToken or PutTokenN.

// BurstLimiter
// Informally, in any large enough time interval, the BurstLimiter limits the
// burst tokens, with a maximum burst size of b events.
// As a special case, if r == Inf (the infinite rate), b is ignored.
// See https://en.wikipedia.org/wiki/Token_bucket for more about token buckets.
//
// The zero value is a valid BurstLimiter, but it will reject all events.
// Use NewFullBurstLimiter to create non-zero Limiters.
//
// BurstLimiter has three main methods, Allow, Reserve, and Wait.
// Most callers should use Wait.
//
// Each of the three methods consumes a single token.
// They differ in their behavior when no token is available.
// If no token is available, Allow returns false.
// If no token is available, Reserve returns a reservation for a future token
// and the amount of time the caller must wait before using it.
// If no token is available, Wait blocks until one can be obtained
// or its associated context.Context is canceled.
//
// The methods AllowN, ReserveN, and WaitN consume n tokens.
type BurstLimiter struct {
	mu                     sync.Mutex
	burst                  int // bucket size, Put Must be called after Get
	tokensChangedListeners []context.Context

	tokens int // unconsumed tokens
}

// Burst returns the maximum burst size. Burst is the maximum number of tokens
// that can be consumed in a single call to Allow, Reserve, or Wait, so higher
// Burst values allow more events to happen at once.
// A zero Burst allows no events, unless limit == Inf.
func (lim *BurstLimiter) Burst() int {
	lim.mu.Lock()
	defer lim.mu.Unlock()
	return lim.burst
}

// Tokens returns the token nums unconsumed.
func (lim *BurstLimiter) Tokens() int {
	lim.mu.Lock()
	defer lim.mu.Unlock()
	return lim.tokens
}

// NewFullBurstLimiter returns a new BurstLimiter inited with full tokens that allows
// events up to burst b and permits bursts of at most b tokens.
func NewFullBurstLimiter(b int) *BurstLimiter {
	return &BurstLimiter{
		burst:  b,
		tokens: b,
	}
}

// NewEmptyBurstLimiter returns a new BurstLimiter inited with zero tokens that allows
// events up to burst b and permits bursts of at most b tokens.
func NewEmptyBurstLimiter(b int) *BurstLimiter {
	return &BurstLimiter{
		burst: b,
	}
}

// SetBurst sets a new burst size for the limiter.
func (lim *BurstLimiter) SetBurst(newBurst int) {
	lim.mu.Lock()
	defer lim.mu.Unlock()
	lim.burst = newBurst
}

// Allow is shorthand for AllowN(time.Now(), 1).
// 当没有可用或足够的事件时，返回false
func (lim *BurstLimiter) Allow() bool {
	return lim.AllowN(1)
}

// AllowN reports whether n events may happen at time now.
// Use this method if you intend to drop / skip events that exceed the rate limit.
// Otherwise use Reserve or Wait.
// 当没有可用或足够的事件时，返回false
func (lim *BurstLimiter) AllowN(n int) bool {
	return lim.reserveN(context.Background(), n, false).ok
}

// Reserve is shorthand for ReserveN(1).
// 当没有可用或足够的事件时，返回 Reservation，和要等待多久才能获得足够的事件。
func (lim *BurstLimiter) Reserve(ctx context.Context) *Reservation {
	return lim.ReserveN(ctx, 1)
}

// ReserveN returns a Reservation that indicates how long the caller must wait before n events happen.
// The BurstLimiter takes this Reservation into account when allowing future events.
// ReserveN returns false if n exceeds the BurstLimiter's burst size.
// Usage example:
//   r := lim.ReserveN(context.Background(), 1)
//   if !r.OK() {
//     // Not allowed to act! Did you remember to set lim.burst to be > 0 ?
//     return
//   }
//   if err:= r.Wait();err!=nil{
//     // Not allowed to act! Reservation or context canceled ?
//     return
//	 }
//   Act()
// Use this method if you wish to wait and slow down in accordance with the rate limit without dropping events.
// If you need to respect a deadline or cancel the delay, use Wait instead.
// To drop or skip events exceeding rate limit, use Allow instead.
// 当没有可用或足够的事件时，返回 Reservation，和要等待多久才能获得足够的事件。
func (lim *BurstLimiter) ReserveN(ctx context.Context, n int) *Reservation {
	r := lim.reserveN(ctx, n, true)
	return r
}

// Wait is shorthand for WaitN(ctx, 1).
func (lim *BurstLimiter) Wait(ctx context.Context) (err error) {
	return lim.WaitN(ctx, 1)
}

// WaitN blocks until lim permits n events to happen.
// It returns an error if n exceeds the BurstLimiter's burst size, the Context is
// canceled, or the expected wait time exceeds the Context's Deadline.
// The burst limit is ignored if the rate limit is Inf.
func (lim *BurstLimiter) WaitN(ctx context.Context, n int) (err error) {
	lim.mu.Lock()
	burst := lim.burst
	lim.mu.Unlock()

	if n > burst {
		return fmt.Errorf("rate: Wait(n=%d) exceeds limiter's burst %d", n, burst)
	}
	// Check if ctx is already cancelled
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	// Reserve
	r := lim.reserveN(ctx, n, true)
	if !r.ok {
		return fmt.Errorf("rate: Wait(n=%d) would exceed context deadline", n)
	}
	// Wait if necessary
	return r.Wait(ctx)
}

// PutToken is shorthand for PutTokenN(ctx, 1).
func (lim *BurstLimiter) PutToken() {
	lim.PutTokenN(1)
}

func (lim *BurstLimiter) PutTokenN(n int) {
	lim.mu.Lock()
	defer lim.mu.Unlock()
	lim.tokens += n
	// drop if overflowed
	if lim.tokens > lim.burst {
		lim.tokens = lim.burst
	}

	for i := 0; i < len(lim.tokensChangedListeners); i++ {
		tokensGot := lim.tokensChangedListeners[i]
		r := tokensGot.Value(expectTokensKey).(*Reservation)
		if (r.tokens) > lim.tokens {
			break
		}
		// remove notified
		if i == len(lim.tokensChangedListeners)-1 {
			lim.tokensChangedListeners = lim.tokensChangedListeners[:i]
		} else {
			lim.tokensChangedListeners = append(lim.tokensChangedListeners[:i], lim.tokensChangedListeners[i+1:]...)
		}
		lim.tokens -= r.tokens
		r.cancelFn()
	}
}

// GetToken is shorthand for GetTokenN(ctx, 1).
func (lim *BurstLimiter) GetToken() (ok bool) {
	return lim.GetTokenN(1)
}

// GetTokenN returns true if token is got
func (lim *BurstLimiter) GetTokenN(n int) (ok bool) {
	lim.mu.Lock()
	defer lim.mu.Unlock()
	return lim.getTokenNLocked(n)
}

// getTokenNLocked returns true if token is got
// advance calculates and returns an updated state for lim resulting from the passage of time.
// lim is not changed.
// getTokenNLocked requires that lim.mu is held.
func (lim *BurstLimiter) getTokenNLocked(n int) (ok bool) {
	if n <= 0 {
		return true
	}
	if lim.tokens >= n {
		lim.tokens -= n
		return true
	}
	return false
}

// reserveN is a helper method for AllowN, ReserveN, and WaitN.
// maxFutureReserve specifies the maximum reservation wait duration allowed.
// reserveN returns Reservation, not *Reservation, to avoid allocation in AllowN and WaitN.
func (lim *BurstLimiter) reserveN(ctx context.Context, n int, wait bool) *Reservation {
	lim.mu.Lock()
	defer lim.mu.Unlock()

	// tokens are enough
	if lim.tokens >= n && len(lim.tokensChangedListeners) == 0 {
		// get n tokens from lim
		if lim.getTokenNLocked(n) {
			return &Reservation{
				ok:     true,
				lim:    lim,
				tokens: 0, // tokens if consumed already,don't wait
			}
		}
	}

	// tokens are not enough

	// Decide result
	var expired bool
	if deadline, has := ctx.Deadline(); has && deadline.Before(time.Now()) {
		expired = true
	}

	ok := n <= lim.burst && !expired && wait

	// Prepare reservation
	r := Reservation{
		ok:     ok,
		lim:    lim,
		tokens: n,
	}
	if ok {
		r.tokensGot, r.cancelFn = context.WithCancel(context.WithValue(ctx, expectTokensKey, &r))
		lim.tokensChangedListeners = append(lim.tokensChangedListeners, r.tokensGot)
	}

	return &r
}

func (lim *BurstLimiter) trackReservationRemove(r *Reservation) {
	lim.mu.Lock()
	defer lim.mu.Unlock()

	for i, tokensGot := range lim.tokensChangedListeners {
		if r == tokensGot.Value(expectTokensKey).(*Reservation) {
			if i == len(lim.tokensChangedListeners)-1 {
				lim.tokensChangedListeners = lim.tokensChangedListeners[:i]
				return
			}
			lim.tokensChangedListeners = append(lim.tokensChangedListeners[:i], lim.tokensChangedListeners[i+1:]...)
			return
		}
	}
	return
}

// A Reservation holds information about events that are permitted by a BurstLimiter to happen after a delay.
// A Reservation may be canceled, which may enable the BurstLimiter to permit additional events.
type Reservation struct {
	ok     bool
	lim    *BurstLimiter
	tokens int // tokens number to consumed(reserved) this time
	//timeToAct time.Time       // now + wait, wait if bucket is not enough
	tokensGot context.Context // chan to notify tokens is put, check if enough
	cancelFn  context.CancelFunc
}

// OK returns whether the limiter can provide the requested number of tokens
// within the maximum wait time.  If OK is false, Delay returns InfDuration, and
// Cancel does nothing.
func (r *Reservation) OK() bool {
	return r.ok
}

// Wait blocks before taking the reserved action
// Wait 当没有可用或足够的事件时，将阻塞等待
func (r *Reservation) Wait(ctx context.Context) error {
	for {
		// Wait if necessary
		if r.tokensGot == nil {
			// We can proceed.
			if r.lim.GetTokenN(r.tokens) {
				return nil
			}
		}
		select {
		case <-r.tokensGot.Done():
			// We can proceed.
			return nil
		case <-ctx.Done():
			// Context was canceled before we could proceed.  Cancel the
			// reservation, which may permit other events to proceed sooner.
			r.Cancel()
			return ctx.Err()
		}
	}
}

// Cancel indicates that the reservation holder will not perform the reserved action
// and reverses the effects of this Reservation on the rate limit as much as possible,
// considering that other reservations may have already been made.
func (r *Reservation) Cancel() {
	if !r.ok {
		return
	}

	if r.tokens == 0 {
		return
	}
	r.lim.trackReservationRemove(r)
	return
}
