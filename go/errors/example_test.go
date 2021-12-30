// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors_test

import (
	"errors"
	"fmt"

	errors_ "github.com/searKing/golang/go/errors"
)

func ExampleMulti() {
	err := errors_.Multi(nil, nil)
	fmt.Println(err)
	fmt.Println("-----")
	err = errors_.Multi(fmt.Errorf("whoops"), nil)
	fmt.Println(err)
	fmt.Println("-----")
	err = errors_.Multi(fmt.Errorf("whoops"), fmt.Errorf("foo"))
	fmt.Println(err)
	fmt.Println("-----")

	// Output:
	// <nil>
	// -----
	// whoops
	// -----
	// whoops|foo
	// -----
}

func ExampleMultiFormat() {
	err := errors_.Multi(nil)
	fmt.Printf("v: %v\n", err)
	fmt.Printf("+v: %+v\n", err)
	fmt.Println("-----")

	err = errors_.Multi(fmt.Errorf("whoops"), nil)
	fmt.Printf("v: %v\n", err)
	fmt.Printf("+v: %+v\n", err)
	fmt.Println("-----")

	err = errors_.Multi(fmt.Errorf("whoops"), fmt.Errorf("foo"))
	fmt.Printf("v: %v\n", err)
	fmt.Printf("+v: %+v\n", err)
	fmt.Println("-----")

	// Output:
	// v: <nil>
	// +v: <nil>
	// -----
	// v: whoops
	// +v: whoops
	// -----
	// v: whoops|foo
	// +v: Multiple errors occurred:
	// |	whoops
	// |	foo
	// -----
}

func ExampleMark() {
	err := errors_.Mark(nil, nil)
	fmt.Println(err)
	fmt.Println("-----")
	err = errors_.Mark(fmt.Errorf("whoops"), nil)
	fmt.Println(err)
	fmt.Println("-----")
	err = errors_.Mark(fmt.Errorf("whoops"), fmt.Errorf("foo"))
	fmt.Println(err)
	fmt.Println("-----")
	err = errors_.Mark(fmt.Errorf("whoops"), fmt.Errorf("foo"), fmt.Errorf("bar"))
	fmt.Println(err)
	fmt.Println("-----")

	// Output:
	// <nil>
	// -----
	// whoops
	// -----
	// whoops
	// -----
	// whoops
	// -----
}

func ExampleMarkFormat() {
	err := errors_.Mark(nil)
	fmt.Printf("v: %v\n", err)
	fmt.Printf("+v: %+v\n", err)
	fmt.Println("-----")

	err = errors_.Mark(fmt.Errorf("whoops"), nil)
	fmt.Printf("v: %v\n", err)
	fmt.Printf("+v: %+v\n", err)
	fmt.Println("-----")

	err = errors_.Mark(fmt.Errorf("whoops"), fmt.Errorf("foo"))
	fmt.Printf("v: %v\n", err)
	fmt.Printf("+v: %+v\n", err)
	fmt.Println("-----")

	err = errors_.Mark(fmt.Errorf("whoops"), fmt.Errorf("foo"), fmt.Errorf("bar"))
	fmt.Printf("v: %v\n", err)
	fmt.Printf("+v: %+v\n", err)
	fmt.Println("-----")

	// Output:
	// v: <nil>
	// +v: <nil>
	// -----
	// v: whoops
	// +v: whoops
	// -----
	// v: whoops
	// +v: Marked errors occurred:
	// |	whoops
	// M	foo
	// -----
	// v: whoops
	// +v: Marked errors occurred:
	// |	whoops
	// M	foo
	// M	bar
	// -----
}

func ExampleMarkIs() {
	var mark = errors.New("mark")
	err := errors_.Mark(nil, mark)
	fmt.Printf("%v\n", errors.Is(err, nil))
	fmt.Printf("%v\n", errors.Is(err, mark))
	fmt.Println("-----")

	err = errors_.Mark(fmt.Errorf("whoops"), nil, mark)
	fmt.Printf("%v\n", errors.Is(err, nil))
	fmt.Printf("%v\n", errors.Is(err, mark))
	fmt.Println("-----")

	err = errors_.Mark(fmt.Errorf("whoops"), fmt.Errorf("foo"), mark)
	fmt.Printf("%v\n", errors.Is(err, nil))
	fmt.Printf("%v\n", errors.Is(err, mark))
	fmt.Println("-----")

	err = errors_.Mark(fmt.Errorf("whoops"), fmt.Errorf("foo"), fmt.Errorf("bar"), mark)
	fmt.Printf("%v\n", errors.Is(err, nil))
	fmt.Printf("%v\n", errors.Is(err, mark))
	fmt.Println("-----")

	// Output:
	// true
	// false
	// -----
	// false
	// true
	// -----
	// false
	// true
	// -----
	// false
	// true
	// -----

}
