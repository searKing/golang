// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time

import (
	"math/rand"
	"time"
)

const DefaultBaseDuration = 5 * time.Millisecond
const DefaultMaxDuration = 1 * time.Second

const ZeroDuration = 0

type DelayHandler interface {
	Delay(attempt int, cap, base, last time.Duration) (delay time.Duration)
}

// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as HTTP handlers. If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler that calls f.
type DelayHandlerFunc func(attempt int, cap, base, last time.Duration) (delay time.Duration)

// ServeHTTP calls f(w, r).
func (f DelayHandlerFunc) Delay(attempt int, cap, base, last time.Duration) (delay time.Duration) {
	return f(attempt, cap, base, last)
}

// Exponential backoff is an algorithm that uses feedback to multiplicatively decrease the rate of some process,
// in order to gradually find an acceptable rate
// see https://cloud.google.com/iot/docs/how-tos/exponential-backoff
// see https://amazonaws-china.com/cn/blogs/architecture/exponential-backoff-and-jitter/

// Work(calls) of Competing Clients(less is better)
// None > Exponential > DecorrelatedJitter > EqualJitter > FullJitter
// Looking at the amount of client work, the number of calls is approximately the same for “Full” and “Equal” jitter,
// and higher for “Decorrelated”. Both cut down work substantially relative to both the no-jitter approaches.

// Completion Time(ms) of Competing Clients(less is better)
// Exponential > EqualJitter > FullJitter > DecorrelatedJitter > None

// none backed-off
// sleep = min(cap, base)
func NoneBackOffDelayHandler(_ int, cap, base, _ time.Duration) time.Duration {
	delay := base
	if delay > cap {
		delay = cap
	}
	return delay
}

// exponentially backed-off
// sleep = min(cap, base * 2 ** attempt)
func ExponentialDelayHandler(attempt int, cap, base, _ time.Duration) time.Duration {
	delay := base * (2 << attempt)
	if delay > cap {
		delay = cap
	}
	return delay
}

// exponentially backed-off with full jitter
// sleep = random_between(0, min(cap, base * 2 ** attempt))
func FullJitterDelayHandler(attempt int, cap, base, _ time.Duration) time.Duration {
	delay := base * (2 << attempt)
	if delay > cap {
		delay = cap
	}
	return time.Duration(rand.Float64() * float64(delay))
}

// exponentially backed-off with equal jitter
// temp = min(cap, base * 2 ** attempt)
// sleep = temp/2 + random_between(temp/2, min(cap, base * 2 ** attempt))
func EqualJitterDelayHandler(attempt int, cap, base, _ time.Duration) time.Duration {
	temp := base * (2 << attempt)
	if temp > cap {
		temp = cap
	}
	delay := temp / 2
	return delay + time.Duration(float64(delay)+rand.Float64()*float64(delay))
}

// exponentially backed-off with decorrelated jitter
// sleep = min(cap, random_between(base, sleep * 3))
func DecorrelatedJitterDelayHandler(_ int, cap, base, last time.Duration) time.Duration {
	delay := base + time.Duration(rand.Float64()*float64(last*3-base))
	if delay > cap {
		delay = cap
	}
	return delay
}

func NewDelay(base, cap time.Duration, h DelayHandler) *Delay {
	return &Delay{
		Base:    base,
		Cap:     cap,
		Handler: h,
	}
}

func NewDefaultExponentialDelay() *Delay {
	return &Delay{
		Base:    DefaultBaseDuration,
		Cap:     DefaultMaxDuration,
		Handler: DelayHandlerFunc(ExponentialDelayHandler),
	}
}

type Delay struct {
	attempt int
	delay   time.Duration
	Base    time.Duration
	Cap     time.Duration
	Handler DelayHandler
}

func (d *Delay) Update() {
	defer func() { d.attempt++ }()
	h := d.Handler
	if h == nil {
		h = DelayHandlerFunc(NoneBackOffDelayHandler)
	}

	d.delay = h.Delay(d.attempt, d.Cap, d.Base, d.delay)
	if max := d.Cap; d.delay > max {
		d.delay = max
	}
}

func (d *Delay) Sleep() {
	d.Update()
	time.Sleep(d.delay)
}

func (d *Delay) Delay() <-chan time.Time {
	d.Update()
	return After(d.delay)
}

func (d *Delay) DelayFunc(f func()) *Timer {
	d.Update()
	return AfterFunc(d.delay, f)
}

// Reset to initial state.
func (d *Delay) Reset() {
	d.delay = ZeroDuration
	d.attempt = 0
}

// Gets duration to wait before retrying the operation or {@link #STOP} to
// indicate that no retries should be made.
func (d *Delay) NextBackOff() (backoff time.Duration, ok bool) {
	d.Update()
	return d.delay, ok
}
