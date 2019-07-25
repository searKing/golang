// Copyright 2019 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Enumeration with an offset.
// Also includes a duplicate.

package main

import (
	"fmt"
	"sync"
)

//go:generate go-syncpool -type "Nums<string>"
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

func ck(nums Nums, str string) {
	val := nums.Get()
	if val != str {
		panic(fmt.Sprintf("Nums<string>.go: %s", str))
	}
}
