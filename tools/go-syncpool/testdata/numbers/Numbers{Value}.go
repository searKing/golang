// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Enumeration with an offset.
// Also includes a duplicate.

package main

import (
	"fmt"
	"sync"
)

//go:generate go-syncpool -type "Numbers<Value>"
type Numbers sync.Pool
type Value string

func main() {
	var numbers Numbers
	numbers.Put("One")
	ck(numbers, "One")
	numbers.Put("Two")
	ck(numbers, "Two")
	numbers.Put("Three")
	ck(numbers, "Three")
	numbers.Put("One")
	ck(numbers, "One")
	numbers.Put("Key(127)")
	ck(numbers, "Key(127)")
}

func ck(numbers Numbers, str Value) {
	val := numbers.Get()
	if val != str {
		panic(fmt.Sprintf("Numbers<Value>.go: %s", str))
	}
}
