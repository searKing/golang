// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Enumeration with an offset.
// Also includes a duplicate.

package main

import (
	"fmt"
	"sync"
	"time"
)

//go:generate go-syncmap -type "Times<int, time.Time>"
type Times sync.Map

const (
	_ = iota
	One
	Two
	Three
	AnotherOne = One // Duplicate; note that AnotherOne doesn't appear below.
)

func main() {
	var times Times
	times.Store(One, time.Time{}.Add(time.Second))
	times.Store(Two, time.Time{}.Add(2*time.Second))
	times.Store(Three, time.Time{}.Add(3*time.Second))
	times.Store(AnotherOne, time.Time{}.Add(time.Second))
	ck(times, One, time.Time{}.Add(time.Second))
	ck(times, Two, time.Time{}.Add(2*time.Second))
	ck(times, Three, time.Time{}.Add(3*time.Second))
	ck(times, AnotherOne, time.Time{}.Add(time.Second))
	ck(times, 127, time.Now())
}

func ck(nums Times, num int, t time.Time) {
	val, loaded := nums.Load(num)
	if num < One || num > Three {
		if loaded {
			panic(fmt.Sprintf("Times<int,time.Time>.go: %s", t))
		}
		return
	}
	if !loaded || val != t {
		panic(fmt.Sprintf("Times<int,time.Time>.go: %s", t))
	}
}
