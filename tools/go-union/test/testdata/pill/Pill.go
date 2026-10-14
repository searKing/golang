// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Enumeration with an offset.
// Also includes a duplicate.

package main

import (
	"fmt"
	time_ "time"
)

type AliasString string
type AliasStruct StructType

//go:generate go-union -type "Pill"
type Pill[T comparable] struct {
	AliasStruct AliasStruct
	// This is Name doc comment
	Name      string // This is Name line comment
	Age       int
	Address   string `union:"-"`
	NameAlias AliasString

	genericFuncType        GenericFuncType[T]
	genericStructType      GenericStructType[T]
	genericEmptyStructType GenericStructTypeEmpty[T]
	pointerType            *[5]T
	structType             time_.Time
	arrayType              [5]T
	funcType               func()
	chanType               chan int
	interfaceType          any
	stringType             string
	mapType                map[string]int64
	indentSliceType        []int64
	selectorSliceType      []time_.Time
	starSelectorSliceType  []*time_.Time
}
type StructTypeEmpty struct{}
type StructType struct {
	Name string
	List []int
}

type GenericStructTypeEmpty[T any] struct{}
type GenericStructType[T any] struct {
	Name string
	List []int
}

type GenericFuncType[T any] func()

func main() {
	var num Pill[int]
	num.Name = "Name"
	u := num.Union()
	fmt.Printf("Union: %v\n", u)

	num = Pill[int]{
		Age: 18,
	}
	fmt.Printf("Union: %v\n", u)
}
