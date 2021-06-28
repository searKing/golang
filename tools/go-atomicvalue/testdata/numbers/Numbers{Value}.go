// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Enumeration with an offset.
// Also includes a duplicate.

package main

import (
	"fmt"
	"sync/atomic"
)

//go:generate go-atomicvalue -type "Numbers<Value>"
type Numbers atomic.Value
type Value string

func main() {
	var numbers Numbers
	numbers.Store("One")
	ck(numbers, "One")
	numbers.Store("Two")
	ck(numbers, "Two")
	numbers.Store("Three")
	ck(numbers, "Three")
	numbers.Store("One")
	ck(numbers, "One")
	numbers.Store("Key(127)")
	ck(numbers, "Key(127)")
}

func ck(numbers Numbers, str Value) {
	val := numbers.Load()
	if val != str {
		panic(fmt.Sprintf("Numbers<Value>.go: %s", str))
	}
}
