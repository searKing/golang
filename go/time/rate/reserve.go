// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rate

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

// A Reservation holds information about events that are permitted by a BurstLimiter to happen after a delay.
// A Reservation may be canceled, which may enable the BurstLimiter to permit additional events.
type Reservation struct {
	*reservation
}

// reservation is the real representation of *Reservation.
// The extra level of indirection ensures that a Reservation
// with a finalizer, that cycle in cyclic structure is not
// guaranteed to be garbage collected
// https://tip.golang.org/doc/gc-guide#Where_Go_Values_Live
type reservation struct {
	ok  bool
	lim *BurstLimiter

	// [0, tokens, burst]
	burst  int // reservation bucket size
	tokens int // tokens got(reserved) from BurstLimiter, Cancel(put back) must be called to the BurstLimiter after Wait

	// timeToAct time.Time       // now + wait, wait if bucket is not enough
	tokensGot         context.Context // chan to notify tokens is put, check if enough
	notifyTokensReady context.CancelFunc

	// test only
	canceled context.CancelFunc
}

func newReservation(gc bool) *Reservation {
	r := &Reservation{reservation: &reservation{}}
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
	if r.canceled != nil {
		r.canceled()
	}
	if r.burst <= 0 {
		return
	}
	defer func() { r.burst = 0 }() // set Reservation as empty
	r.lim.trackReservationRemove(r.reservation)
	r.lim.PutTokenN(r.tokens)
	r.tokens = 0
	return
}

// PutToken (as Complete): refill all tokens taken by the Reservation back to BurstLimiter.
// PutToken is shorthand for Cancel().
func (r *Reservation) PutToken() {
	r.Cancel()
}
