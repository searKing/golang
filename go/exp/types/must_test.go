// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types_test

import (
	"errors"
	"fmt"

	types_ "github.com/searKing/golang/go/exp/types"
)

func ExampleMust() {
	sentinel := errors.New("boom")
	defer func() {
		r := recover()
		if r != nil {
			fmt.Printf("panicked with %T (%v)", r, r)
		}
	}()
	n := types_.Must(1, nil)
	s := types_.Must("hello", nil)
	fmt.Println(n, s)
	_ = types_.Must(struct{}{}, sentinel)
	// Output:
	// 1 hello
	// panicked with *errors.errorString (boom)
}
