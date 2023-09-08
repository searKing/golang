// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types_test

import (
	"fmt"

	"github.com/searKing/golang/go/exp/types"
)

type stringerVal struct {
	s string
}

func (s stringerVal) String() string {
	return s.s
}

func ExampleAny() {

	fmt.Printf("%#v\n", (*stringerVal)(nil)) // typed nil implements [fmt.Stringer]
	fmt.Printf("%#v\n", types.Any((*stringerVal)(nil)))

	// Output:
	// (*types_test.stringerVal)(nil)
	// <nil>
}
