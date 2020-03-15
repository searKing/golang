// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time

import "time"

const DefaultInitDuration = 5 * time.Millisecond
const DefaultStepTimes = 2
const DefaultMaxDuration = 1 * time.Second

const ZeroDuration = 0

func NewDefaultDelay() *Delay {
	return &Delay{
		InitDuration: DefaultInitDuration,
		MaxDuration:  DefaultMaxDuration,
		DelayAgainHandler: func(delay time.Duration) time.Duration {
			return delay * DefaultStepTimes
		},
	}
}

type Delay struct {
	delay             time.Duration
	InitDuration      time.Duration
	MaxDuration       time.Duration
	DelayAgainHandler func(delay time.Duration) time.Duration
}

func (d *Delay) Update() {
	if d.delay == ZeroDuration {
		d.delay = d.InitDuration
	} else {
		if d.DelayAgainHandler != nil {
			d.delay = d.DelayAgainHandler(d.delay)
		}
	}
	if max := d.MaxDuration; d.delay > max {
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
func (d *Delay) Reset() {
	d.delay = ZeroDuration
}
func (d *Delay) Duration() time.Duration {
	return d.delay
}
