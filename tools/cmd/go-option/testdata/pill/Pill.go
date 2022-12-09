// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"strings"
)

////go:generate go-option -type "Pill"
//type Pill[T comparable] struct {
//	// This is Name doc comment
//	Name      string // This is Name line comment
//	Age       string `option:",short"`
//	Address   string `option:"-"`
//	NameAlias string `option:"Title,"`
//
//	arrayType     [5]T
//	funcType      func()
//	interfaceType interface{}
//	mapType       map[string]int
//	sliceType     []int64
//}

func NewPill[T comparable](options ...PillOption[T]) *Pill[T] {
	return (&Pill[T]{}).ApplyOptions()
}

type Value string

func main() {
	var num *Pill[int]
	num = &Pill[int]{}
	ck(num, "")
	num = NewPill(WithPillName[int]("Name"))
	ck(num, "")
}

func ck[T comparable](num *Pill[T], str string) {
	name := num.Name
	if strings.Compare(name, str) != 0 {
		panic(fmt.Sprintf("Pill.go: %s", str))
	}
}
