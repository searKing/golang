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

//go:generate go-syncpool -type "Times<time.Time>"
type Times sync.Pool

func main() {
	var times Times
	times.Put(time.Time{}.Add(time.Second))
	ck(times, time.Time{}.Add(time.Second))
	times.Put(time.Time{}.Add(2 * time.Second))
	ck(times, time.Time{}.Add(2*time.Second))
	times.Put(time.Time{}.Add(3 * time.Second))
	ck(times, time.Time{}.Add(3*time.Second))
}

func ck(nums Times, t time.Time) {
	val := nums.Get()
	if val != t {
		panic(fmt.Sprintf("Times<time.Time>.go: %s", t))
	}
}
