// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package rate
// The key observation and some code is borrowed from
// golang.org/x/time/rate/rate.go
package rate

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

const (
	starvationThresholdNs = 1e6
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
// Reorder Buffer
// It allows instructions to be committed in-order.
// - Allocated by `Reserve`  or `ReserveN` into account when allowing future events
// - Wait by `Wait` or `WaitN` blocks until lim permits n events to happen.
// - Allow and Wait Complete by `PutToken` or `PutTokenN`
// - Reserve Complete by `Cancel` of the Reservation self, GC Cancel supported
// See https://en.wikipedia.org/wiki/Re-order_buffer for more about Reorder buffer.
// See https://web.archive.org/web/20040724215416/http://lgjohn.okstate.edu/6253/lectures/reorder.pdf for more about Reorder buffer.
//
// The zero value is a valid BurstLimiter, but it will reject all events.
// Use NewFullBurstLimiter to create non-zero Limiters.
//
// BurstLimiter has three main methods, Allow, Reserve, and Wait.
// Most callers should use Wait for token bucket.
// Most callers should use Reserve for Reorder buffer.
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
// AllowN is shorthand for GetTokenN.
// Use this method if you intend to drop / skip events that exceed the rate limit.
// Otherwise, use Reserve or Wait.
// 当没有可用或足够的事件时，返回false
func (lim *BurstLimiter) AllowN(n int) bool {
	return lim.GetTokenN(n)
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
//
//	    // Allocate: The dispatch stage reserves space in the reorder buffer for instructions in program order.
//		r := lim.ReserveN(context.Background(), 1)
//		if !r.OK() {
//			// Not allowed to act! Did you remember to set lim.burst to be > 0 ?
//			return
//		}
//
//		// Execute: out-of-order execution
//		Act()
//
//		// Wait: The complete stage must wait for instructions to finish execution.
//		if err:= r.Wait(); err!=nil {
//		// Not allowed to act! Reservation or context canceled ?
//			return
//		}
//		// Complete: Finished instructions are allowed to write results in order into the architected registers.
//		// It allows instructions to be committed in-order.
//		defer r.PutToken()
//
//		// Execute: in-order execution
//		Act()
//
// Use this method if you wish to wait and slow down in accordance with the rate limit without dropping events.
// If you need to respect a deadline or cancel the delay, use Wait instead.
// To drop or skip events exceeding rate limit, use Allow instead.
// 当没有可用或足够的事件时，返回 Reservation，和要等待多久才能获得足够的事件。
// See https://en.wikipedia.org/wiki/Re-order_buffer for more about Reorder buffer.
// See https://web.archive.org/web/20040724215416/http://lgjohn.okstate.edu/6253/lectures/reorder.pdf for more about Reorder buffer.
func (lim *BurstLimiter) ReserveN(ctx context.Context, n int) *Reservation {
	r := lim.reserveN(ctx, n, true, true)
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
	r := lim.reserveN(ctx, n, true, false)
	if r.Ready() { // tokens already hold by the Reservation
		return nil
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
		if r.burst <= 0 {
			// remove notified
			if i == len(lim.tokensChangedListeners)-1 {
				lim.tokensChangedListeners = lim.tokensChangedListeners[:i]
			} else {
				lim.tokensChangedListeners = append(lim.tokensChangedListeners[:i], lim.tokensChangedListeners[i+1:]...)
			}
			r.notifyTokensReady()
			continue
		}

		tokensWait := r.burst - r.tokens

		// tokens in the Bucket is not enough for the Reservation
		if lim.tokens < tokensWait {
			r.tokens += lim.tokens
			lim.tokens = 0
			break
		}

		// enough
		r.tokens = r.burst
		lim.tokens -= tokensWait
		// remove notified
		if i == len(lim.tokensChangedListeners)-1 {
			lim.tokensChangedListeners = lim.tokensChangedListeners[:i]
		} else {
			lim.tokensChangedListeners = append(lim.tokensChangedListeners[:i], lim.tokensChangedListeners[i+1:]...)
		}
		r.notifyTokensReady()
		continue
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
func (lim *BurstLimiter) reserveN(ctx context.Context, n int, wait bool, gc bool) *Reservation {
	if n <= 0 {
		r := newReservation(gc)
		r.lim = lim
		r.burst = 0
		r.tokens = 0
		return r
	}

	lim.mu.Lock()
	defer lim.mu.Unlock()

	// tokens are enough
	if lim.tokens >= n && len(lim.tokensChangedListeners) == 0 {
		// get n tokens from lim
		if lim.getTokenNLocked(n) {
			r := newReservation(gc)
			r.lim = lim
			r.burst = n
			r.tokens = n
			return r
		}
	}

	// tokens are not enough

	// Decide result
	var expired bool
	if deadline, has := ctx.Deadline(); has && deadline.Before(time.Now()) {
		expired = true
	}

	addToListener := n <= lim.burst && !expired && wait

	// Prepare reservation
	r := newReservation(gc)
	r.lim = lim
	r.burst = n
	if addToListener {
		r.tokensGot, r.notifyTokensReady = context.WithCancel(context.WithValue(ctx, expectTokensKey, r))
		lim.tokensChangedListeners = append(lim.tokensChangedListeners, r.tokensGot)
	}

	return r
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
	ok  bool
	lim *BurstLimiter

	// [0, tokens, burst]
	burst  int // reservation bucket size
	tokens int // tokens got(reserved) from BurstLimiter, Cancel(put back) must be called to the BurstLimiter after Wait

	//timeToAct time.Time       // now + wait, wait if bucket is not enough
	tokensGot         context.Context // chan to notify tokens is put, check if enough
	notifyTokensReady context.CancelFunc
}

func newReservation(gc bool) *Reservation {
	r := &Reservation{}
	if gc {
		runtime.SetFinalizer(r, (*Reservation).Cancel)
	}
	return r
}

func (r *Reservation) removeGC() *Reservation {
	// no need for a finalizer anymore
	runtime.SetFinalizer(r, nil)
	return r
}

// OK returns whether the limiter can provide the requested number of tokens
// within the maximum wait time. If OK is false, Delay returns InfDuration, and
// Cancel does nothing.
func (r *Reservation) OK() bool {
	return r.burst <= r.lim.Burst()
}

// Ready returns whether the limiter can provide the requested number of tokens
// within the maximum wait time. If Ready is false, Wait returns nil directly, and
// Cancel or GC does put back the token reserved in the Reservation.
// If Ready is false, WaitN blocks until lim permits n events to happen.
func (r *Reservation) Ready() bool {
	return r.tokens >= r.burst
}

// Wait blocks before taking the reserved action
// Wait 当没有可用或足够的事件时，将阻塞等待
func (r *Reservation) Wait(ctx context.Context) error {
	if r.burst <= 0 {
		r.burst = 0
		r.tokens = r.burst
		return nil
	}

	if r.Ready() {
		return nil
	}

	var burst = r.lim.Burst()
	if r.burst > burst {
		return fmt.Errorf("rate: Wait(n=%d) exceeds limiter's burst %d", r.burst, burst)
	}
	timer := time.NewTimer(starvationThresholdNs)
	defer timer.Stop()
	for {
		// fast path
		if r.tokensGot == nil {
			// We can proceed.
			if r.lim.GetTokenN(r.burst - r.tokens) {
				r.tokens = r.burst
				return nil
			}
			// Wait if necessary
			if !timer.Stop() {
				<-timer.C
			}
			timer.Reset(starvationThresholdNs)

			select {
			case <-ctx.Done():
				// Context was canceled before we could proceed.  Cancel the
				// reservation, which may permit other events to proceed sooner.
				r.Cancel()
				return ctx.Err()
			case <-timer.C:
				break
			}
			continue
		}

		// Wait if necessary
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
	defer func() {
		// no need for a finalizer anymore
		runtime.SetFinalizer(r, nil)
	}()
	if r.burst <= 0 {
		return
	}
	defer func() { r.burst = 0 }() // set Reservation as empty
	r.lim.trackReservationRemove(r)
	r.lim.PutTokenN(r.tokens)
	r.tokens = 0
	return
}

// PutToken (as Complete): refill all tokens taken by the Reservation back to BurstLimiter.
// PutToken is shorthand for Cancel().
func (r *Reservation) PutToken() {
	r.Cancel()
}
