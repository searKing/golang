// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Enumeration with an offset.
// Also includes a duplicate.

package main

import (
	"fmt"
	"time"
)

//go:generate go-nulljson -type "Times<time.Time>"

func main() {
	var times Times
	err := times.Scan(time.Time{}.Add(time.Second))
	ckError(err)
	ck(times, time.Time{}.Add(time.Second))

	err = times.Scan(time.Time{}.Add(2 * time.Second))
	ckError(err)
	ck(times, time.Time{}.Add(2*time.Second))

	err = times.Scan(time.Time{}.Add(3 * time.Second))
	ckError(err)
	ck(times, time.Time{}.Add(3*time.Second))
}

func ck(nums Times, t time.Time) {
	val := nums.Data
	if val != t {
		panic(fmt.Sprintf("Times<time.Time>.go: %s", t))
	}
}

func ckError(err error) {
	if err != nil {
		panic(fmt.Sprintf("Times<time.Time>.go: error happened %s", err))
	}
}
