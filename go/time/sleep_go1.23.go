// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.23

package time

import (
	"time"
)

// Timer to fix time: Timer.Stop documentation example easily leads to deadlocks
// https://github.com/golang/go/issues/27169
//
// Deprecated: Use [time.Timer] instead since go1.23.
// https://github.com/golang/go/issues/37196
// https://github.com/golang/go/issues/14383
type Timer = time.Timer

// Deprecated: Use [time.NewTimer] instead since go1.23.
func NewTimer(d time.Duration) *Timer {
	return time.NewTimer(d)
}

// Deprecated: Use [time.Timer] instead since go1.23.
func WrapTimer(t *time.Timer) *Timer {
	return t
}

// Deprecated: Use [time.After] instead since go1.23.
func After(d time.Duration) <-chan time.Time {
	return NewTimer(d).C
}

// Deprecated: Use [time.AfterFunc] instead since go1.23.
func AfterFunc(d time.Duration, f func()) *Timer {
	t := time.AfterFunc(d, func() {
		f()
	})
	return t
}
