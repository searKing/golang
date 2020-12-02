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
