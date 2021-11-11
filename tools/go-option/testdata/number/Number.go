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
type Number struct {
	arrayType     [5]int64
	funcType      func()
	interfaceType interface{}
	mapType       map[string]int64
	sliceType     []int64
	name          string
}

func NewNumber(options ...NumberOption) *Number {
	return (&Number{}).ApplyOptions()
}

type Value string

func main() {
	var num *Number
	num = &Number{}
	ck(num, "")
	num = NewNumber(WithNumberName("Name"))
	ck(num, "")
}

func ck(num *Number, str string) {
	name := num.name
	if strings.Compare(name, str) != 0 {
		panic(fmt.Sprintf("Numbers.go: %s", str))
	}
}
