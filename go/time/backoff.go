// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time

import (
	"math/rand"
	"time"
)

// see https://cloud.google.com/iot/docs/how-tos/exponential-backoff
// Implementation of {@link BackOff} that increases the back off period for each retry attempt using
// a randomization function that grows exponentially.
//
// <p>{@link #NextBackOff()} is calculated using the following formula:
//
// <pre>
// randomized_interval =
// retry_interval // (random value in range [1 - randomization_factor, 1 + randomization_factor])
// </pre>
//
// <p>In other words {@link #NextBackOff()} will range between the randomization factor
// percentage below and above the retry interval. For example, using 2 seconds as the base retry
// interval and 0.5 as the randomization factor, the actual back off period used in the next retry
// attempt will be between 1 and 3 seconds.
//
// <p><b>Note:</b> max_interval caps the retry_interval and not the randomized_interval.
//
// <p>If the time elapsed since an {@link ExponentialBackOff} instance is created goes past the
// max_elapsed_time then the method {@link #NextBackOff()} starts returning {@link
// BackOff#STOP}. The elapsed time can be reset by calling {@link #reset()}.
//
// <p>Example: The default retry_interval is .5 seconds, default randomization_factor is 0.5,
// default multiplier is 1.5 and the default max_interval is 1 minute. For 10 tries the sequence
// will be (values in seconds) and assuming we go over the max_elapsed_time on the 10th try:
//
// <pre>
// request#     retry_interval     randomized_interval
//
// 1             0.5                [0.25,   0.75]
// 2             0.75               [0.375,  1.125]
// 3             1.125              [0.562,  1.687]
// 4             1.687              [0.8435, 2.53]
// 5             2.53               [1.265,  3.795]
// 6             3.795              [1.897,  5.692]
// 7             5.692              [2.846,  8.538]
// 8             8.538              [4.269, 12.807]
// 9            12.807              [6.403, 19.210]
// 10           19.210              {@link BackOff#STOP}
// </pre>
//
// <p>Implementation is not thread-safe.
//

const (

	// The default initial interval value (0.5 seconds).
	DefaultInitialInterval = 500 * time.Millisecond

	// The default randomization factor (0.5 which results in a random period ranging between 50%
	// below and 50% above the retry interval).
	DefaultRandomizationFactor = 0.5

	// The default multiplier value (1.5 which is 50% increase per back off).
	DefaultMultiplier = 1.5

	// The default maximum back off time (1 minute).
	DefaultMaxInterval = time.Minute

	// The default maximum elapsed time (15 minutes).
	DefaultMaxElapsedDuration = 15 * time.Minute
)

// Code borrowed from https://github.com/googleapis/google-http-java-client/blob/master/google-http-client/
// src/main/java/com/google/api/client/util/BackOff.java
type BackOff interface {
	// Reset to initial state.
	Reset()

	// Gets duration to wait before retrying the operation to
	// indicate that no retries should be made.
	// ok indicates that no more retries should be made.
	// Example usage:
	// var backOffDuration, ok = backoff.NextBackOff();
	// if (!ok) {
	// 	// do not retry operation
	// } else {
	//	// sleep for backOffDuration milliseconds and retry operation
	// }
	NextBackOff() (backoff time.Duration, ok bool)
}

//  Fixed back-off policy whose back-off time is always zero, meaning that the operation is retried
//  immediately without waiting.

type ZeroBackOff struct {
}

func (o *ZeroBackOff) Reset() {
}
func (o *ZeroBackOff) NextBackOff() (backoff time.Duration, ok bool) {
	return 0, true
}

// Fixed back-off policy that always returns {@code #STOP} for {@link #NextBackOff()},
// meaning that the operation should not be retried.

type StopBackOff struct {
}

func (o *StopBackOff) Reset() {
}
func (o *StopBackOff) NextBackOff() (backoff time.Duration, ok bool) {
	return 0, false
}

//go:generate go-option -type "ExponentialBackOff"

// Code borrowed from https://github.com/googleapis/google-http-java-client/blob/master/google-http-client/
// src/main/java/com/google/api/client/util/ExponentialBackOff.java
type ExponentialBackOff struct {
	// The current retry interval.
	currentInterval time.Duration
	// The initial retry interval.
	initialInterval time.Duration

	// The randomization factor to use for creating a range around the retry interval.
	// A randomization factor of 0.5 results in a random period ranging between 50% below and 50%
	// above the retry interval.
	randomizationFactor float64

	// The value to multiply the current interval with for each retry attempt.
	multiplier float64

	// The maximum value of the back off period. Once the retry interval reaches this
	// value it stops increasing.
	maxInterval time.Duration

	// The system time in nanoseconds. It is calculated when an ExponentialBackOffPolicy instance is
	// created and is reset when {@link #reset()} is called.
	startTime time.Time

	// The maximum elapsed time after instantiating {@link ExponentialBackOff} or calling {@link
	// #reset()} after which {@link #NextBackOff()} returns {@link BackOff#STOP}.
	maxElapsedDuration time.Duration
}

func NewExponentialBackOff(opts ...ExponentialBackOffOption) *ExponentialBackOff {
	o := &ExponentialBackOff{
		initialInterval:     DefaultInitialInterval,
		randomizationFactor: DefaultRandomizationFactor,
		multiplier:          DefaultMultiplier,
		maxInterval:         DefaultMaxInterval,
		maxElapsedDuration:  DefaultMaxElapsedDuration,
	}
	o.ApplyOptions(opts...)
	o.Reset()
	return o
}

// Sets the interval back to the initial retry interval and restarts the timer.
func (o *ExponentialBackOff) Reset() {
	o.currentInterval = o.initialInterval
	o.startTime = time.Now()
}

/**
 * {@inheritDoc}
 *
 * <p>This method calculates the next back off interval using the formula: randomized_interval =
 * retry_interval +/- (randomization_factor * retry_interval)
 *
 * <p>Subclasses may override if a different algorithm is required.
 */
func (o *ExponentialBackOff) NextBackOff() (backoff time.Duration, ok bool) {
	// Make sure we have not gone over the maximum elapsed time.
	if o.GetElapsedDuration() > o.maxElapsedDuration {
		return 0, false
	}
	randomizedInterval := o.GetRandomValueFromInterval(o.randomizationFactor, o.currentInterval)
	o.incrementCurrentInterval()
	return randomizedInterval, true
}

/**
 * Returns a random value from the interval [randomizationFactor * currentInterval,
 * randomizationFactor * currentInterval].
 */
func (o *ExponentialBackOff) GetRandomValueFromInterval(
	randomizationFactor float64, currentInterval time.Duration) time.Duration {
	return Jitter(currentInterval, randomizationFactor)
}

// Returns the initial retry interval.
func (o *ExponentialBackOff) GetInitialInterval() time.Duration {
	return o.initialInterval
}

// Returns the randomization factor to use for creating a range around the retry interval.
// A randomization factor of 0.5 results in a random period ranging between 50% below and 50%
// above the retry interval.
func (o *ExponentialBackOff) GetRandomizationFactor() float64 {
	return o.randomizationFactor
}

// Returns the current retry interval.
func (o *ExponentialBackOff) GetCurrentInterval() time.Duration {
	return o.currentInterval
}

// Returns the value to multiply the current interval with for each retry attempt.
func (o *ExponentialBackOff) GetMultiplier() float64 {
	return o.multiplier
}

// Returns the maximum value of the back off period. Once the current interval
// reaches this value it stops increasing.
func (o *ExponentialBackOff) GetMaxInterval() time.Duration {
	return o.maxInterval
}

// Returns the maximum elapsed time.
// If the time elapsed since an {@link ExponentialBackOff} instance is created goes past the
// max_elapsed_time then the method {@link #NextBackOff()} starts returning STOP.
// The elapsed time can be reset by calling
func (o *ExponentialBackOff) GetMaxElapsedDuration() time.Duration {
	return o.maxElapsedDuration
}

// Returns the elapsed time since an {@link ExponentialBackOff} instance is
// created and is reset when {@link #reset()} is called.
// The elapsed time is computed using {@link System#nanoTime()}.
func (o *ExponentialBackOff) GetElapsedDuration() time.Duration {
	return time.Now().Sub(o.startTime)
}

// Increments the current interval by multiplying it with the multiplier.
func (o *ExponentialBackOff) incrementCurrentInterval() {
	// Check for overflow, if overflow is detected set the current interval to the max interval.
	if o.currentInterval*time.Duration(o.multiplier) >= o.maxInterval {
		o.currentInterval = o.maxInterval
	} else {
		o.currentInterval *= time.Duration(o.multiplier)
	}
}

// Jitter returns a time.Duration between
// [duration - maxFactor*duration, duration + maxFactor*duration].
//
// This allows clients to avoid converging on periodic behavior.
func Jitter(duration time.Duration, maxFactor float64) time.Duration {
	delta := time.Duration(maxFactor) * duration
	minInterval := duration - delta
	maxInterval := duration + delta
	// Get a random value from the range [minInterval, maxInterval].
	// The formula used below has a +1 because if the minInterval is 1 and the maxInterval is 3 then
	// we want a 33% chance for selecting either 1, 2 or 3.
	return minInterval + time.Duration(rand.Float64()*float64(maxInterval-minInterval+1))
}
