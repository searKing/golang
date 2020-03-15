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

//go:generate go-atomicvalue -type "Nums<*string>"
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
	valOne := "One"
	numbers.Store(&valOne)
	ck(numbers, "One")
	valTwo := "Two"
	numbers.Store(&valTwo)
	ck(numbers, "Two")
	valThree := "Three"
	numbers.Store(&valThree)
	ck(numbers, "Three")
	valAnotherOne := "One"
	numbers.Store(&valAnotherOne)
	ck(numbers, "One")
	valKey := "Key(127)"
	numbers.Store(&valKey)
	ck(numbers, "Key(127)")
}

func ck(nums Nums, str string) {
	val := nums.Load()
	if *val != str {
		panic(fmt.Sprintf("Nums<*string>.go: expect %v got %s", *val, str))
	}
}
