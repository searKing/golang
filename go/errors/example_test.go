// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors_test

import (
	"fmt"

	"github.com/searKing/golang/go/errors"
)

func ExampleMulti() {
	err := errors.Multi(nil, nil)
	fmt.Println(err)
	err = errors.Multi(fmt.Errorf("whoops"), nil)
	fmt.Println(err)
	err = errors.Multi(fmt.Errorf("whoops"), fmt.Errorf("foo"))
	fmt.Println(err)

	// Output:
	// <nil>
	// whoops
	// whoops|foo
}

func ExampleFormat() {
	err := errors.Multi(nil)
	fmt.Printf("v: %v\n", err)
	fmt.Printf("+v: %+v\n", err)

	err = errors.Multi(fmt.Errorf("whoops"), nil)
	fmt.Printf("v: %v\n", err)
	fmt.Printf("+v: %+v\n", err)

	err = errors.Multi(fmt.Errorf("whoops"), fmt.Errorf("foo"))
	fmt.Printf("v: %v\n", err)
	fmt.Printf("+v: %+v\n", err)

	// Output:
	// v: <nil>
	// +v: <nil>
	// v: whoops
	// +v: whoops
	// v: whoops|foo
	// +v: Multiple errors occurred:
	// |	whoops
	// |	foo
}
