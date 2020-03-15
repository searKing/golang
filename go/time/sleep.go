// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time

import (
	"time"
)

// https://github.com/golang/go/issues/27169
// Timer to fix time: Timer.Stop documentation example easily leads to deadlocks
type Timer struct {
	*time.Timer
}

// Stop prevents the Timer from firing, with the channel drained.
// Stop ensures the channel is empty after a call to Stop.
// Stop == std [Stop + drain]
// It returns true if the call stops the timer, false if the timer has already
// expired or been stopped.
// Stop does not close the channel, to prevent a read from the channel succeeding
// incorrectly.
func (t *Timer) Stop() bool {
	if t.Timer == nil {
		panic("time: Stop called on uninitialized Timer")
	}

	active := t.Timer.Stop()
	if !active {
		// drain the channel, prevents the Timer from blocking on Send to t.C by sendTime, t.C is reused.
		// The underlying Timer is not recovered by the garbage collector until the timer fires.
		// consume the channel only once for the channel can be triggered only one time at most before Stop is called.
	L:
		for {
			select {
			case _, ok := <-t.Timer.C:
				if !ok {
					break L
				}
			default:
				break L
			}
		}
	}
	return active
}

// Reset changes the timer to expire after duration d.
// Reset can be invoked anytime, which enhances std time.Reset
// Reset == std [Stop + drain + Reset]
// It returns true if the timer had been active, false if the timer had
// expired or been stopped.
func (t *Timer) Reset(d time.Duration) bool {
	if t.Timer == nil {
		panic("time: Reset called on uninitialized Timer")
	}
	active := t.Stop()
	t.Timer.Reset(d)
	return active
}

func NewTimer(d time.Duration) *Timer {
	return &Timer{
		Timer: time.NewTimer(d),
	}
}

func WrapTimer(t *time.Timer) *Timer {
	return &Timer{
		Timer: t,
	}
}

func After(d time.Duration) <-chan time.Time {
	return NewTimer(d).C
}

func AfterFunc(d time.Duration, f func()) *Timer {
	t := &Timer{}
	t.Timer = time.AfterFunc(d, func() {
		f()
	})
	return t
}
