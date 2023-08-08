// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Enumeration with an offset.
// Also includes a duplicate.

package main

import (
	"fmt"
	"strings"
	time_ "time"
)

//go:generate go-option -type "Number"
type Number[T comparable] struct {
	// This is Name doc comment
	Name      string // This is Name line comment
	Age       string `option:",short"`
	Address   string `option:"-"`
	NameAlias string `option:"Title,"`

	genericType   GenericType[T]
	pointerType   *[5]T
	structType    time_.Time
	arrayType     [5]T
	funcType      func()
	interfaceType interface{}
	mapType       map[string]int64
	sliceType     []int64
	name          string
}

type GenericType[T any] struct{}

func NewNumber[T comparable](opts ...NumberOption[T]) *Number[T] {
	return (&Number[T]{}).ApplyOptions(opts...)
}

type Value string

func main() {
	var num *Number[int]
	num = &Number[int]{}
	ck(num, "")
	num = NewNumber(WithNumberName[int]("Name"))
	ck(num, "Name")
}

func ck[T comparable](num *Number[T], str string) {
	name := num.name
	if strings.Compare(name, str) != 0 {
		panic(fmt.Sprintf("Numbers.go: %s", str))
	}
}
