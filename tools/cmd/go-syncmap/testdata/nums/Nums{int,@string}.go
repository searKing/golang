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

//go:generate go-syncmap -type "Nums<int, *string>"
type Nums sync.Map

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
	numbers.Store(One, &valOne)
	valTwo := "Two"
	numbers.Store(Two, &valTwo)
	valThree := "Three"
	numbers.Store(Three, &valThree)
	valAnotherOne := "One"
	numbers.Store(AnotherOne, &valAnotherOne)
	ck(numbers, One, "One")
	ck(numbers, Two, "Two")
	ck(numbers, Three, "Three")
	ck(numbers, AnotherOne, "One")
	ck(numbers, 127, "Key(127)")
}

func ck(nums Nums, num int, str string) {
	val, loaded := nums.Load(num)
	if num < One || num > Three {
		if loaded {
			panic(fmt.Sprintf("Nums<int,*string>.go: %s", str))
		}
		return
	}
	if !loaded || *val != str {
		panic(fmt.Sprintf("Nums<int,*string>.go: key %d, expect %v, got %s", num, *val, str))
	}
}
