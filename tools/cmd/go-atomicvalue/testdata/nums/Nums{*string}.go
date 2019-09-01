// Copyright 2019 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Enumeration with an offset.
// Also includes a duplicate.

package main

import (
	"fmt"
	"sync/atomic"
)

//go:generate go-atomicvalue -type "Nums<string>"
type Nums atomic.Value

const (
	_ = iota
	One
	Two
	Three
	AnotherOne = One // Duplicate; note that AnotherOne doesn't appear below.
)

func main() {
	var numbers Nums
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

func ck(nums Nums, str string) {
	val := nums.Load()
	if val != str {
		panic(fmt.Sprintf("Nums<string>.go: %s", str))
	}
}
