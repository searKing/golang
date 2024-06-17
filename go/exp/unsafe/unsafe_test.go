// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package unsafe_test

import (
	"fmt"

	"github.com/searKing/golang/go/exp/unsafe"
)

func ExampleBytes() {
	fmt.Printf("%v\n", unsafe.Bytes(int8(1)))
	fmt.Printf("%v\n", unsafe.Bytes(int16(1)))
	fmt.Printf("%v\n", unsafe.Bytes(int32(1)))
	fmt.Printf("%v\n", unsafe.Bytes([0]byte{}))
	fmt.Printf("%v\n", unsafe.Bytes([2]byte{}))
	fmt.Printf("%v\n", unsafe.Bytes([2]byte{1, 2}))
	fmt.Printf("%v\n", unsafe.Bytes(struct {
		Age  int8
		Name [2]byte
	}{Age: 1, Name: [2]byte{2, 3}}))

	// Output:
	// [1]
	// [1 0]
	// [1 0 0 0]
	// []
	// [0 0]
	// [1 2]
	// [1 2 3]
}

func ExamplePlacement() {
	fmt.Printf("%+v\n", *unsafe.Placement[int8]([]byte{1}))
	fmt.Printf("%+v\n", *unsafe.Placement[int16]([]byte{1, 0}))
	fmt.Printf("%+v\n", *unsafe.Placement[int32]([]byte{1, 0, 0, 0}))
	fmt.Printf("%+v\n", unsafe.Placement[[0]byte](nil))
	fmt.Printf("%+v\n", unsafe.Placement[[0]byte]([]byte{}))
	fmt.Printf("%+v\n", unsafe.Placement[[2]byte](nil))
	fmt.Printf("%+v\n", unsafe.Placement[[2]byte]([]byte{}))
	fmt.Printf("%+v\n", *unsafe.Placement[[2]byte]([]byte{1, 2}))
	fmt.Printf("%+v\n", *unsafe.Placement[struct {
		Age  int8
		Name [2]byte
	}]([]byte{1, 2, 3}))

	// Output:
	// 1
	// 1
	// 1
	// <nil>
	// <nil>
	// <nil>
	// <nil>
	// [1 2]
	// {Age:1 Name:[2 3]}
}
