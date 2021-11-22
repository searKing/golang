// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Enumeration with an offset.
// Also includes a duplicate.

package main

import (
	"fmt"
)

//go:generate go-nulljson -type "Numbers<Value>" -nullable -protojson
type Value string

func main() {
	var numbers Numbers
	err := numbers.Scan(`"One"`)
	ckError(err)
	ck(numbers, "One")

	err = numbers.Scan(`"Two"`)
	ckError(err)
	ck(numbers, "Two")

	err = numbers.Scan(`"Three"`)
	ckError(err)
	ck(numbers, "Three")

	err = numbers.Scan(`"One"`)
	ckError(err)
	ck(numbers, "One")

	err = numbers.Scan(`"Key(127)"`)
	ckError(err)
	ck(numbers, "Key(127)")
}

func ck(numbers Numbers, str Value) {
	val := numbers.Data
	if val != str {
		panic(fmt.Sprintf("Numbers<Value>.go: %s", str))
	}
}

func ckError(err error) {
	if err != nil {
		panic(fmt.Sprintf("Numbers<Value>.go: error happened %s", err))
	}
}
