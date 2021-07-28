// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time

import "time"

func WithExponentialBackOffOptionInitialInterval(duration time.Duration) ExponentialBackOffOption {
	return ExponentialBackOffOptionFunc(func(o *ExponentialBackOff) {
		o.initialInterval = duration
	})
}

func WithExponentialBackOffOptionRandomizationFactor(factor float64) ExponentialBackOffOption {
	return ExponentialBackOffOptionFunc(func(o *ExponentialBackOff) {
		o.randomizationFactor = factor
	})
}

func WithExponentialBackOffOptionMultiplier(multiplier float64) ExponentialBackOffOption {
	return ExponentialBackOffOptionFunc(func(o *ExponentialBackOff) {
		o.multiplier = multiplier
	})
}

func WithExponentialBackOffOptionMaxInterval(maxInterval time.Duration) ExponentialBackOffOption {
	return ExponentialBackOffOptionFunc(func(o *ExponentialBackOff) {
		o.maxInterval = maxInterval
	})
}

func WithExponentialBackOffOptionMaxElapsedDuration(duration time.Duration) ExponentialBackOffOption {
	return ExponentialBackOffOptionFunc(func(o *ExponentialBackOff) {
		o.maxElapsedDuration = duration
	})
}

func WithExponentialBackOffOptionMaxElapsedCount(count int) ExponentialBackOffOption {
	return ExponentialBackOffOptionFunc(func(o *ExponentialBackOff) {
		o.maxElapsedCount = count
	})
}

func WithExponentialBackOffOptionNoLimit() ExponentialBackOffOption {
	return ExponentialBackOffOptionFunc(func(o *ExponentialBackOff) {
		o.initialInterval = DefaultInitialInterval
		o.randomizationFactor = DefaultRandomizationFactor
		o.multiplier = DefaultMultiplier
		o.maxInterval = -1
		o.maxElapsedDuration = -1
		o.maxElapsedCount = -1
	})
}

func WithExponentialBackOffOptionGRPC() ExponentialBackOffOption {
	return ExponentialBackOffOptionFunc(func(o *ExponentialBackOff) {
		o.initialInterval = 1.0 * time.Second
		o.randomizationFactor = 0.2
		o.multiplier = 1.6
		o.maxInterval = 120 * time.Second
		o.maxElapsedDuration = -1
		o.maxElapsedCount = -1
	})
}
