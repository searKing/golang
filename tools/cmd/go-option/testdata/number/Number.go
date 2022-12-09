// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Enumeration with an offset.
// Also includes a duplicate.

package main

import (
	"fmt"
	"strings"
)

//go:generate go-option -type "Number"
type Number[T comparable] struct {
	arrayType     [5]T
	funcType      func()
	interfaceType interface{}
	mapType       map[string]int64
	sliceType     []int64
	name          string
}

func NewNumber[T comparable](options ...NumberOption[T]) *Number[T] {
	return (&Number[T]{}).ApplyOptions()
}

type Value string

func main() {
	var num *Number[int]
	num = &Number[int]{}
	ck(num, "")
	num = NewNumber(WithNumberName[int]("Name"))
	ck(num, "")
}

func ck[T comparable](num *Number[T], str string) {
	name := num.name
	if strings.Compare(name, str) != 0 {
		panic(fmt.Sprintf("Numbers.go: %s", str))
	}
}
