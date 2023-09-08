// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types_test

import (
	"fmt"
	"math"

	"github.com/searKing/golang/go/exp/types"
)

func ExampleConv() {

	fmt.Printf("%#x\n", types.Conv[int8, uint8](-1))           // overflow
	fmt.Printf("%d\n", types.Conv[uint8, int8](math.MaxUint8)) // overflow
	fmt.Printf("%d\n", types.Conv[float32, int32](-1))
	fmt.Printf("%f\n", types.Conv[int32, float32](-1))
	// Output:
	// 0xff
	// -1
	// -1
	// -1.000000

}

func ExampleConvBinary() {

	fmt.Printf("%#x\n", types.ConvBinary[int8, uint8](-1))           // overflow
	fmt.Printf("%d\n", types.ConvBinary[uint8, int8](math.MaxUint8)) // overflow
	fmt.Printf("%d\n", types.ConvBinary[float32, int32](-1))
	fmt.Printf("%f\n", types.ConvBinary[int32, float32](-1082130432))
	fmt.Printf("%d\n", types.ConvBinary[*stringerVal, int64]((*stringerVal)(nil)))
	fmt.Printf("%d\n", types.ConvBinary[*stringerVal, int64](nil))
	// Output:
	// 0xff
	// -1
	// -1082130432
	// -1.000000
	// 0
	// 0

}
