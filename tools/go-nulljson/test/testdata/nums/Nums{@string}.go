// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Enumeration with an offset.
// Also includes a duplicate.

package main

import (
	"fmt"
)

//go:generate go-nulljson -type "Nums<*string>"

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
	err := numbers.Scan(&valOne)
	ckError(err)
	ck(numbers, "One")

	valTwo := "Two"
	err = numbers.Scan(&valTwo)
	ckError(err)
	ck(numbers, "Two")

	valThree := "Three"
	err = numbers.Scan(&valThree)
	ckError(err)
	ck(numbers, "Three")

	valAnotherOne := "One"
	err = numbers.Scan(&valAnotherOne)
	ckError(err)
	ck(numbers, "One")
	valKey := "Key(127)"
	err = numbers.Scan(&valKey)
	ckError(err)
	ck(numbers, "Key(127)")
}

func ck(nums Nums, str string) {
	val := nums.Data
	if *val != str {
		panic(fmt.Sprintf("Nums<*string>.go: expect %v got %s", *val, str))
	}
}

func ckError(err error) {
	if err != nil {
		panic(fmt.Sprintf("Nums<*string>.go: error happened %s", err))
	}
}
