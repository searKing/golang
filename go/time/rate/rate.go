// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The key observation and some code (shr) is borrowed from
// golang.org/x/time/rate/rate.go
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
// at rate r tokens per second.
// Informally, in any large enough time interval, the BurstLimiter limits the
// rate to r tokens per second, with a maximum burst size of b events.
// As a special case, if r == Inf (the infinite rate), b is ignored.
// See https://en.wikipedia.org/wiki/Token_bucket for more about token buckets.
//
// The zero value is a valid BurstLimiter, but it will reject all events.
// Use NewBurstLimiter to create non-zero Limiters.
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
	burst              int // 桶的大小，且必须归还
	tokensChangedChans []context.Context

	mu     sync.Mutex
	tokens int // unconsumed tokens
}

// Burst returns the maximum burst size. Burst is the maximum number of tokens
// that can be consumed in a single call to Allow, Reserve, or Wait, so higher
// Burst values allow more events to happen at once.
// A zero Burst allows no events, unless limit == Inf.
func (lim *BurstLimiter) Burst() int {
	return lim.burst
}

// NewBurstLimiter returns a new BurstLimiter that allows events up to rate r and permits
// bursts of at most b tokens.
func NewBurstLimiter(b int) *BurstLimiter {
	return &BurstLimiter{
		burst:  b,
		tokens: b,
	}
}

// Allow is shorthand for AllowN(time.Now(), 1).
func (lim *BurstLimiter) Allow() bool {
	return lim.AllowN(1)
}

// AllowN reports whether n events may happen at time now.
// Use this method if you intend to drop / skip events that exceed the rate limit.
// Otherwise use Reserve or Wait.
func (lim *BurstLimiter) AllowN(n int) bool {
	return lim.reserveN(context.Background(), n, false).ok
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

func (r *Reservation) Wait(ctx context.Context) error {
	// Wait if necessary
	if r.tokensGot == nil {
		// We can proceed.
		r.lim.GetTokenN(r.tokens)
		return nil
	}
	select {
	case <-r.tokensGot.Done():
		// We can proceed.
		r.lim.GetTokenN(r.tokens)
		return nil
	case <-ctx.Done():
		// Context was canceled before we could proceed.  Cancel the
		// reservation, which may permit other events to proceed sooner.
		r.Cancel()
		return ctx.Err()
	}
}

// InfDuration is the duration returned by Delay when a Reservation is not OK.
const InfDuration = time.Duration(1<<63 - 1)

// Cancel indicates that the reservation holder will not perform the reserved action
// and reverses the effects of this Reservation on the rate limit as much as possible,
// considering that other reservations may have already been made.
func (r *Reservation) Cancel() {
	if !r.ok {
		return
	}

	r.lim.mu.Lock()
	defer r.lim.mu.Unlock()

	if r.tokens == 0 {
		return
	}

	for i, tokensGot := range r.lim.tokensChangedChans {
		if r == tokensGot.Value(expectTokensKey).(*Reservation) {
			if i == len(r.lim.tokensChangedChans)-1 {
				r.lim.tokensChangedChans = r.lim.tokensChangedChans[:i]
				return
			}
			r.lim.tokensChangedChans = append(r.lim.tokensChangedChans[:i], r.lim.tokensChangedChans[i+1:]...)
			return
		}
	}
	return
}

// Reserve is shorthand for ReserveN(1).
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
//   time.Sleep(r.Delay())
//   Act()
// Use this method if you wish to wait and slow down in accordance with the rate limit without dropping events.
// If you need to respect a deadline or cancel the delay, use Wait instead.
// To drop or skip events exceeding rate limit, use Allow instead.
func (lim *BurstLimiter) ReserveN(ctx context.Context, n int) *Reservation {
	r := lim.reserveN(ctx, n, true)
	return &r
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
	if n > lim.burst {
		return fmt.Errorf("rate: Wait(n=%d) exceeds limiter's burst %d", n, lim.burst)
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

	for i := 0; i < len(lim.tokensChangedChans); i++ {
		tokensGot := lim.tokensChangedChans[i]
		r := tokensGot.Value(expectTokensKey).(*Reservation)
		if (r.tokens) > lim.tokens {
			break
		}
		// remove notified
		lim.tokensChangedChans = append(lim.tokensChangedChans[:i], lim.tokensChangedChans[i+1:]...)
		lim.tokens -= r.tokens
		r.cancelFn()
	}
}

// GetToken is shorthand for GetTokenN(ctx, 1).
func (lim *BurstLimiter) GetToken() {
	lim.GetTokenN(1)
}

func (lim *BurstLimiter) GetTokenN(n int) {
	lim.mu.Lock()
	defer lim.mu.Unlock()
	lim.tokens -= n
	// drop if overflowed
	if lim.tokens > lim.burst {
		lim.tokens = lim.burst
	}
	if lim.tokens < 0 {
		lim.tokens = 0
	}
}

// reserveN is a helper method for AllowN, ReserveN, and WaitN.
// maxFutureReserve specifies the maximum reservation wait duration allowed.
// reserveN returns Reservation, not *Reservation, to avoid allocation in AllowN and WaitN.
func (lim *BurstLimiter) reserveN(ctx context.Context, n int, wait bool) Reservation {
	lim.mu.Lock()
	defer lim.mu.Unlock()

	// tokens are enough
	if lim.tokens >= n && len(lim.tokensChangedChans) == 0 {
		return Reservation{
			ok:     true,
			lim:    lim,
			tokens: n,
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
		lim.tokensChangedChans = append(lim.tokensChangedChans, r.tokensGot)
	}

	return r
}
