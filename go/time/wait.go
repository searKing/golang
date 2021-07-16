// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time

import (
	"context"
	"time"

	"github.com/searKing/golang/go/runtime"
)

// Forever calls f every period for ever.
//
// Forever is syntactic sugar on top of Forever, without resetCh.
// Example: time.Second
// 2021/04/09 12:45:08 Apr  9 12:45:08
// 2021/04/09 12:45:09 Apr  9 12:45:09
// 2021/04/09 12:45:10 Apr  9 12:45:10
// 2021/04/09 12:45:11 Apr  9 12:45:11
// 2021/04/09 12:45:12 Apr  9 12:45:12
// 2021/04/09 12:45:13 Apr  9 12:45:13
// 2021/04/09 12:45:14 Apr  9 12:45:14
func Forever(f func(), period time.Duration) {
	Until(context.Background(), func(ctx context.Context) { f() }, period)
}

// ForeverWithReset calls f every period for ever.
//
// ForeverWithReset is syntactic sugar on top of UntilWithReset.
// Example: time.Second
// 2021/04/09 12:45:08 Apr  9 12:45:08
// 2021/04/09 12:45:09 Apr  9 12:45:09
// 2021/04/09 12:45:10 Apr  9 12:45:10
// 2021/04/09 12:45:11 Apr  9 12:45:11
// 2021/04/09 12:45:12 Apr  9 12:45:12
// 2021/04/09 12:45:13 Apr  9 12:45:13
// 2021/04/09 12:45:14 Apr  9 12:45:14
func ForeverWithReset(f func(), resetCh chan struct{}, period time.Duration) {
	UntilWithReset(context.Background(), func(ctx context.Context) { f() }, resetCh, period)
}

// Until loops until context is done, running f every period.
//
// Until is syntactic sugar on top of UntilWithReset, without resetCh.
func Until(ctx context.Context, f func(ctx context.Context), period time.Duration) {
	UntilWithReset(ctx, f, nil, period)
}

// UntilWithReset loops until context is done, running f every period.
//
// UntilWithReset is syntactic sugar on top of JitterUntilWithReset with zero jitter factor and
// with sliding = true (which means the timer for period starts after the f
// completes).
// Example: time.Second for period and sleep in f
// 2021/04/09 12:48:03 Apr  9 12:48:03
// 2021/04/09 12:48:05 Apr  9 12:48:05
// 2021/04/09 12:48:07 Apr  9 12:48:07
// 2021/04/09 12:48:09 Apr  9 12:48:09
// 2021/04/09 12:48:11 Apr  9 12:48:11
// 2021/04/09 12:48:13 Apr  9 12:48:13
func UntilWithReset(ctx context.Context, f func(ctx context.Context), resetCh chan struct{}, period time.Duration) {
	JitterUntilWithReset(ctx, f, resetCh, true,
		WithExponentialBackOffOptionRandomizationFactor(0),
		WithExponentialBackOffOptionMultiplier(1),
		WithExponentialBackOffOptionInitialInterval(period),
		WithExponentialBackOffOptionMaxElapsedDuration(-1))
}

// NonSlidingUntil loops until context is done, running f every
// period.
//
// NonSlidingUntil is syntactic sugar on top of NonSlidingUntilWithReset, without resetCh.
func NonSlidingUntil(ctx context.Context, f func(ctx context.Context), period time.Duration) {
	NonSlidingUntilWithReset(ctx, f, nil, period)
}

// NonSlidingUntilWithReset loops until context is done, running f every
// period.
//
// NonSlidingUntilWithReset is syntactic sugar on top of JitterUntilWithReset with zero jitter
// factor, with sliding = false (meaning the timer for period starts at the same
// time as the function starts).
// Example: time.Second for period and sleep in f
// 2021/04/09 12:45:08 Apr  9 12:45:08
// 2021/04/09 12:45:09 Apr  9 12:45:09
// 2021/04/09 12:45:10 Apr  9 12:45:10
// 2021/04/09 12:45:11 Apr  9 12:45:11
// 2021/04/09 12:45:12 Apr  9 12:45:12
// 2021/04/09 12:45:13 Apr  9 12:45:13
// 2021/04/09 12:45:14 Apr  9 12:45:14
func NonSlidingUntilWithReset(ctx context.Context, f func(ctx context.Context), resetCh chan struct{}, period time.Duration) {
	JitterUntilWithReset(ctx, f, resetCh, false,
		WithExponentialBackOffOptionRandomizationFactor(0),
		WithExponentialBackOffOptionMultiplier(1),
		WithExponentialBackOffOptionInitialInterval(period),
		WithExponentialBackOffOptionMaxElapsedDuration(-1))
}

// JitterUntil loops until context is done, running f every period.
// JitterUntil is syntactic sugar on top of JitterUntilWithReset, without resetCh.
func JitterUntil(ctx context.Context, f func(ctx context.Context), sliding bool, opts ...ExponentialBackOffOption) {
	JitterUntilWithReset(ctx, f, nil, sliding, opts...)
}

// JitterUntilWithReset loops until context is done, running f every period.
//
// period set by WithExponentialBackOffOptionInitialInterval
// jitterFactor set by WithExponentialBackOffOptionRandomizationFactor
// If jitterFactor is positive, the period is jittered before every run of f.
// If jitterFactor is not positive, the period is unchanged and not jittered.
//
// If sliding is true, the period is computed after f runs. If it is false then
// period includes the runtime for f.
// backoff is reset if resetCh has data
//
// Cancel context to stop. f may not be invoked if context is already expired.
func JitterUntilWithReset(ctx context.Context, f func(ctx context.Context), resetCh chan struct{}, sliding bool, opts ...ExponentialBackOffOption) {
	BackoffUntilWithReset(ctx, f, resetCh, NewExponentialBackOff(opts...), sliding)
}

// BackoffUntil loops until context is done, run f every duration given by BackoffManager.
// BackoffUntil is syntactic sugar on top of BackoffUntilWithReset, without resetCh.
func BackoffUntil(ctx context.Context, f func(ctx context.Context), backoff BackOff, sliding bool) {
	BackoffUntilWithReset(ctx, f, nil, backoff, sliding)
}

// BackoffUntilWithReset loops until context is done, run f every duration given by BackoffManager.
//
// If sliding is true, the period is computed after f runs. If it is false then
// period includes the runtime for f.
// backoff is reset if resetCh has data
func BackoffUntilWithReset(ctx context.Context,
	f func(ctx context.Context), resetCh chan struct{}, backoff BackOff, sliding bool) {
	var elapsed time.Duration
	var ok bool

	var drainResetCh = func() {
		// To ensure the channel is empty, check the
		// return value and drain the channel.
		for {
			select {
			case <-resetCh:
			default:
				return
			}
		}
	}
	for {
		select {
		case <-ctx.Done():
			return
		case <-resetCh:
			backoff.Reset()
			drainResetCh()
		default:
		}

		var cost Cost
		if !sliding {
			cost.Start()
			elapsed, ok = backoff.NextBackOff()
		}

		func() {
			defer runtime.DefaultPanic.Recover()
			f(ctx)
		}()
		if !sliding {
			elapsed -= cost.Elapse()
		}

		if sliding {
			elapsed, ok = backoff.NextBackOff()
		}
		if !ok {
			return
		}

		func() {
			if elapsed <= 0 {
				return
			}
			timer := time.NewTimer(elapsed)
			defer timer.Stop()

			// NOTE: b/c there is no priority selection in golang
			// it is possible for this to race, meaning we could
			// trigger t.C and stopCh, and t.C select falls through.
			// In order to mitigate we re-check stopCh at the beginning
			// of every loop to prevent extra executions of f().
			select {
			case <-ctx.Done():
				return
			case <-timer.C:
			case <-resetCh:
				backoff.Reset()
				drainResetCh()
			}
		}()
	}
}
