// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time_test

import (
	"context"
	"math"
	"testing"
	"time"

	time_ "github.com/searKing/golang/go/time"
)

func TestFixedBackOff(t *testing.T) {
	nonSliding := time_.NonSlidingBackOff(time.Second)
	tests := []struct {
		name    string
		backoff time_.BackOff
		min     time.Duration
		max     time.Duration
		ok      bool
	}{
		{"*stop", &time_.StopBackOff{}, 0, 0, false},
		{"*non-sliding", &nonSliding, time.Second, time.Second, true},
		{"jitter", time_.JitterBackOff(time.Second, 0.5), 500 * time.Millisecond, 1500*time.Millisecond + 1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for range 3 {
				got, ok := tt.backoff.NextBackOff()
				if ok != tt.ok || got < tt.min || got > tt.max {
					t.Fatalf("NextBackOff() = %v, %t; want [%v, %v], %t", got, ok, tt.min, tt.max, tt.ok)
				}
			}
		})
	}
}

func TestBackoffUntilNonSliding(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var n int
	backoff := time_.NonSlidingBackOff(time.Millisecond)
	time_.BackoffUntil(ctx, func(ctx context.Context) {
		if n++; n == 3 {
			cancel()
		}
	}, &backoff, true)
	if n != 3 {
		t.Errorf("BackoffUntil runs f %d times; want 3", n)
	}
}

func TestExponentialBackOff_MaxInterval(t *testing.T) {
	tests := []struct {
		name       string
		multiplier float64
		want       []time.Duration
	}{
		// 0.5, 0.75, ..., 50.6, 60 rather than 75.9
		{name: "default", multiplier: 1.5, want: []time.Duration{500 * time.Millisecond, 750 * time.Millisecond, 1125 * time.Millisecond}},
		{name: "fractional multiplier", multiplier: 2.5, want: []time.Duration{500 * time.Millisecond, 1250 * time.Millisecond, 3125 * time.Millisecond}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := time_.NewDefaultExponentialBackOff(
				time_.WithExponentialBackOffOptionMultiplier(tt.multiplier),
				time_.WithExponentialBackOffOptionMaxElapsedDuration(-1))
			for i, want := range tt.want {
				if got := b.GetCurrentInterval(); got != want {
					t.Errorf("#%d interval = %v, want %v", i, got, want)
				}
				b.NextBackOff()
			}
			for i := 0; i < 100; i++ {
				if got := b.GetCurrentInterval(); got > time_.DefaultMaxInterval {
					t.Fatalf("interval = %v, want no more than %v", got, time_.DefaultMaxInterval)
				}
				b.NextBackOff()
			}
			if got := b.GetCurrentInterval(); got != time_.DefaultMaxInterval {
				t.Errorf("interval = %v, want %v", got, time_.DefaultMaxInterval)
			}
		})
	}
}

func TestExponentialBackOff_NoLimitSaturated(t *testing.T) {
	b := time_.NewExponentialBackOff(
		time_.WithExponentialBackOffOptionInitialInterval(math.MaxInt64/2),
		time_.WithExponentialBackOffOptionMultiplier(3))
	for i := 0; i < 3; i++ {
		b.NextBackOff()
		if got := b.GetCurrentInterval(); got != math.MaxInt64 {
			t.Fatalf("#%d interval = %v, want %v", i, got, time.Duration(math.MaxInt64))
		}
	}
}
