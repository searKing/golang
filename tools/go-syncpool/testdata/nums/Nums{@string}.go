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

//go:generate go-syncpool -type "Nums<*string>"
type Nums sync.Pool

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
	numbers.Put(&valOne)
	ck(numbers, "One")
	valTwo := "Two"
	numbers.Put(&valTwo)
	ck(numbers, "Two")
	valThree := "Three"
	numbers.Put(&valThree)
	ck(numbers, "Three")
	valAnotherOne := "One"
	numbers.Put(&valAnotherOne)
	ck(numbers, "One")
	valKey := "Key(127)"
	numbers.Put(&valKey)
	ck(numbers, "Key(127)")
}

func ck(nums Nums, str string) {
	val := nums.Get()
	if *val != str {
		panic(fmt.Sprintf("Nums<*string>.go: expect %v got %s", *val, str))
	}
}
