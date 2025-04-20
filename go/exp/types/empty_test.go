// Copyright 2025 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types_test

import (
	"fmt"

	"github.com/searKing/golang/go/exp/types"
)

func ExampleIsEmptyStruct() {
	type EmptyAlias = struct{}
	type UserEmpty struct{}
	fmt.Printf("%#v\n", types.IsEmptyStruct[struct{}]())
	fmt.Printf("%#v\n", types.IsEmptyStruct[EmptyAlias]())
	fmt.Printf("%#v\n", types.IsEmptyStruct[UserEmpty]())

	// Output:
	// true
	// true
	// false
}
