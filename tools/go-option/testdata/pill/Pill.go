// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"strings"
	time_ "time"
)

//go:generate go-option -type "Pill"
type Pill[T comparable] struct {
	// This is Name doc comment
	Name      string // This is Name line comment
	Age       string `option:",short"`
	Address   string `option:"-"`
	NameAlias string `option:"Title,"`

	genericType   GenericType[T]
	structType    time_.Time
	arrayType     [5]T
	pointerType   *[5]T
	funcType      func()
	interfaceType any
	mapType       map[string]int64
	sliceType     []int64
}
type GenericType[T any] struct{}

func NewPill[T comparable](opts ...PillOption[T]) *Pill[T] {
	return (&Pill[T]{}).ApplyOptions(opts...)
}

type Value string

func main() {
	var num *Pill[int]
	num = &Pill[int]{}
	ck(num, "")
	num = NewPill(WithPillName[int]("Name"))
	ck(num, "Name")
}

func ck[T comparable](num *Pill[T], str string) {
	name := num.Name
	if strings.Compare(name, str) != 0 {
		panic(fmt.Sprintf("Pill.go: %s", str))
	}
}
