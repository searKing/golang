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
// Forever is syntactic sugar on top of Until.
func Forever(f func(), period time.Duration) {
	Until(context.Background(), func(ctx context.Context) {
		f()
	}, period)
}

// Until loops until context is done, running f every period.
//
// Until is syntactic sugar on top of JitterUntil with zero jitter factor and
// with sliding = true (which means the timer for period starts after the f
// completes).
func Until(ctx context.Context, f func(ctx context.Context), period time.Duration) {
	JitterUntil(ctx, f, true,
		WithExponentialBackOffOptionRandomizationFactor(0.5),
		WithExponentialBackOffOptionMultiplier(1),
		WithExponentialBackOffOptionInitialInterval(period))
}

// NonSlidingUntil loops until context is done, running f every
// period.
//
// NonSlidingUntil is syntactic sugar on top of JitterUntil with zero jitter
// factor, with sliding = false (meaning the timer for period starts at the same
// time as the function starts).
func NonSlidingUntil(ctx context.Context, f func(ctx context.Context), period time.Duration) {
	JitterUntil(ctx, f, false,
		WithExponentialBackOffOptionRandomizationFactor(0.5),
		WithExponentialBackOffOptionMultiplier(1),
		WithExponentialBackOffOptionInitialInterval(period))
}

// JitterUntil loops until context is done, running f every period.
//
// period set by WithExponentialBackOffOptionInitialInterval
// jitterFactor set by WithExponentialBackOffOptionRandomizationFactor
// If jitterFactor is positive, the period is jittered before every run of f.
// If jitterFactor is not positive, the period is unchanged and not jittered.
//
// If sliding is true, the period is computed after f runs. If it is false then
// period includes the runtime for f.
//
// Cancel context to stop. f may not be invoked if context is already expired.
func JitterUntil(ctx context.Context, f func(ctx context.Context), sliding bool, opts ...ExponentialBackOffOption) {
	BackoffUntil(ctx, f, NewExponentialBackOff(opts...), sliding)
}

// BackoffUntil loops until context is done, run f every duration given by BackoffManager.
//
// If sliding is true, the period is computed after f runs. If it is false then
// period includes the runtime for f.
func BackoffUntil(ctx context.Context, f func(ctx context.Context), backoff BackOff, sliding bool) {
	var elapsed time.Duration
	var ok bool
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		if !sliding {
			elapsed, ok = backoff.NextBackOff()
		}

		func() {
			defer runtime.DefaultPanic.Recover()
			f(ctx)
		}()

		if sliding {
			elapsed, ok = backoff.NextBackOff()
		}

		func() {
			if !ok || elapsed < 0 {
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
			}
		}()
	}
}
