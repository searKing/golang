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

func ExampleMultiErrorFormat() {
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

func ExampleMultiErrorIs() {
	var asWrapErr = &wrapError{msg: "wrap"}

	var errorStringErr = errors.New("errorString")
	var wrapErr = &wrapError{msg: "wrap"}
	err := errors_.Multi(nil, errorStringErr)
	fmt.Printf("err = %+v\n", err)
	fmt.Printf("errors.Is(err, nil) = %v\n", errors.Is(err, nil))
	fmt.Printf("errors.Is(err, errorStringErr) = %v\n", errors.Is(err, errorStringErr))
	asWrapErr = &wrapError{msg: "nil"}
	fmt.Printf("errors.As(err, &asWrapErr) = %v\n", errors.As(err, &asWrapErr))
	fmt.Printf("wrapErr = %v\n", error(asWrapErr))
	fmt.Println("-----")

	err = errors_.Multi(fmt.Errorf("whoops"), nil, errorStringErr, wrapErr)
	fmt.Printf("err = %+v\n", err)
	fmt.Printf("errors.Is(err, nil) = %v\n", errors.Is(err, nil))
	fmt.Printf("errors.Is(err, errorStringErr) = %v\n", errors.Is(err, errorStringErr))
	asWrapErr = &wrapError{}
	fmt.Printf("errors.As(err, &asWrapErr) = %v\n", errors.As(err, &asWrapErr))
	fmt.Printf("wrapErr = %v\n", asWrapErr)
	fmt.Println("-----")

	err = errors_.Multi(fmt.Errorf("whoops"), fmt.Errorf("foo"), errorStringErr, wrapErr)
	fmt.Printf("err = %+v\n", err)
	fmt.Printf("errors.Is(err, nil) = %v\n", errors.Is(err, nil))
	fmt.Printf("errors.Is(err, errorStringErr) = %v\n", errors.Is(err, errorStringErr))
	asWrapErr = &wrapError{}
	fmt.Printf("errors.As(err, &asWrapErr) = %v\n", errors.As(err, &asWrapErr))
	fmt.Printf("wrapErr = %v\n", asWrapErr)
	fmt.Println("-----")

	err = errors_.Multi(fmt.Errorf("whoops"), fmt.Errorf("foo"), fmt.Errorf("bar"), errorStringErr, wrapErr)
	fmt.Printf("err = %+v\n", err)
	fmt.Printf("errors.Is(err, nil) = %v\n", errors.Is(err, nil))
	fmt.Printf("errors.Is(err, errorStringErr) = %v\n", errors.Is(err, errorStringErr))
	asWrapErr = &wrapError{}
	fmt.Printf("errors.As(err, &asWrapErr) = %v\n", errors.As(err, &asWrapErr))
	fmt.Printf("wrapErr = %v\n", asWrapErr)
	fmt.Println("-----")

	err = errors_.Multi(err, fmt.Errorf("alice"), fmt.Errorf("bob"), wrapErr)
	fmt.Printf("err = %+v\n", err)
	fmt.Printf("errors.Is(err, nil) = %v\n", errors.Is(err, nil))
	fmt.Printf("errors.Is(err, errorStringErr) = %v\n", errors.Is(err, errorStringErr))
	asWrapErr = &wrapError{}
	fmt.Printf("errors.As(err, &asWrapErr) = %v\n", errors.As(err, &asWrapErr))
	fmt.Printf("wrapErr = %v\n", asWrapErr)
	fmt.Println("-----")

	// Output:
	// err = errorString
	// errors.Is(err, nil) = false
	// errors.Is(err, errorStringErr) = true
	// errors.As(err, &asWrapErr) = false
	// wrapErr = nil
	// -----
	// err = Multiple errors occurred:
	// |	whoops
	// |	errorString
	// |	wrap
	// errors.Is(err, nil) = false
	// errors.Is(err, errorStringErr) = true
	// errors.As(err, &asWrapErr) = true
	// wrapErr = wrap
	// -----
	// err = Multiple errors occurred:
	// |	whoops
	// |	foo
	// |	errorString
	// |	wrap
	// errors.Is(err, nil) = false
	// errors.Is(err, errorStringErr) = true
	// errors.As(err, &asWrapErr) = true
	// wrapErr = wrap
	// -----
	// err = Multiple errors occurred:
	// |	whoops
	// |	foo
	// |	bar
	// |	errorString
	// |	wrap
	// errors.Is(err, nil) = false
	// errors.Is(err, errorStringErr) = true
	// errors.As(err, &asWrapErr) = true
	// wrapErr = wrap
	// -----
	// err = Multiple errors occurred:
	// |	Multiple errors occurred:
	// |	whoops
	// |	foo
	// |	bar
	// |	errorString
	// |	wrap
	// |	alice
	// |	bob
	// |	wrap
	// errors.Is(err, nil) = false
	// errors.Is(err, errorStringErr) = true
	// errors.As(err, &asWrapErr) = true
	// wrapErr = wrap
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
	err = errors_.Mark(err, fmt.Errorf("alice"), fmt.Errorf("bob"))
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
	// whoops
	// -----
}

func ExampleMarkErrorFormat() {
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

	err = errors_.Mark(err, fmt.Errorf("alice"), fmt.Errorf("bob"))
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
	// v: whoops
	// +v: Marked errors occurred:
	// |	Marked errors occurred:
	// |	whoops
	// M	foo
	// M	bar
	// M	alice
	// M	bob
	// -----
}

func ExampleMarkErrorIs() {
	var asWrapErr = &wrapError{msg: "wrap"}

	var errorStringErr = errors.New("errorString")
	var wrapErr = &wrapError{msg: "wrap"}
	err := errors_.Mark(nil, errorStringErr)
	fmt.Printf("err = %+v\n", err)
	fmt.Printf("errors.Is(err, nil) = %v\n", errors.Is(err, nil))
	fmt.Printf("errors.Is(err, errorStringErr) = %v\n", errors.Is(err, errorStringErr))
	asWrapErr = &wrapError{msg: "nil"}
	fmt.Printf("errors.As(err, &asWrapErr) = %v\n", errors.As(err, &asWrapErr))
	fmt.Printf("wrapErr = %v\n", error(asWrapErr))
	fmt.Println("-----")

	err = errors_.Mark(fmt.Errorf("whoops"), nil, errorStringErr, wrapErr)
	fmt.Printf("err = %+v\n", err)
	fmt.Printf("errors.Is(err, nil) = %v\n", errors.Is(err, nil))
	fmt.Printf("errors.Is(err, errorStringErr) = %v\n", errors.Is(err, errorStringErr))
	asWrapErr = &wrapError{}
	fmt.Printf("errors.As(err, &asWrapErr) = %v\n", errors.As(err, &asWrapErr))
	fmt.Printf("wrapErr = %v\n", asWrapErr)
	fmt.Println("-----")

	err = errors_.Mark(fmt.Errorf("whoops"), fmt.Errorf("foo"), errorStringErr, wrapErr)
	fmt.Printf("err = %+v\n", err)
	fmt.Printf("errors.Is(err, nil) = %v\n", errors.Is(err, nil))
	fmt.Printf("errors.Is(err, errorStringErr) = %v\n", errors.Is(err, errorStringErr))
	asWrapErr = &wrapError{}
	fmt.Printf("errors.As(err, &asWrapErr) = %v\n", errors.As(err, &asWrapErr))
	fmt.Printf("wrapErr = %v\n", asWrapErr)
	fmt.Println("-----")

	err = errors_.Mark(fmt.Errorf("whoops"), fmt.Errorf("foo"), fmt.Errorf("bar"), errorStringErr, wrapErr)
	fmt.Printf("err = %+v\n", err)
	fmt.Printf("errors.Is(err, nil) = %v\n", errors.Is(err, nil))
	fmt.Printf("errors.Is(err, errorStringErr) = %v\n", errors.Is(err, errorStringErr))
	asWrapErr = &wrapError{}
	fmt.Printf("errors.As(err, &asWrapErr) = %v\n", errors.As(err, &asWrapErr))
	fmt.Printf("wrapErr = %v\n", asWrapErr)
	fmt.Println("-----")

	err = errors_.Mark(err, fmt.Errorf("alice"), fmt.Errorf("bob"), wrapErr)
	fmt.Printf("err = %+v\n", err)
	fmt.Printf("errors.Is(err, nil) = %v\n", errors.Is(err, nil))
	fmt.Printf("errors.Is(err, errorStringErr) = %v\n", errors.Is(err, errorStringErr))
	asWrapErr = &wrapError{}
	fmt.Printf("errors.As(err, &asWrapErr) = %v\n", errors.As(err, &asWrapErr))
	fmt.Printf("wrapErr = %v\n", asWrapErr)
	fmt.Println("-----")

	// Output:
	// err = <nil>
	// errors.Is(err, nil) = true
	// errors.Is(err, errorStringErr) = false
	// errors.As(err, &asWrapErr) = false
	// wrapErr = nil
	// -----
	// err = Marked errors occurred:
	// |	whoops
	// M	errorString
	// M	wrap
	// errors.Is(err, nil) = false
	// errors.Is(err, errorStringErr) = true
	// errors.As(err, &asWrapErr) = true
	// wrapErr = wrap
	// -----
	// err = Marked errors occurred:
	// |	whoops
	// M	foo
	// M	errorString
	// M	wrap
	// errors.Is(err, nil) = false
	// errors.Is(err, errorStringErr) = true
	// errors.As(err, &asWrapErr) = true
	// wrapErr = wrap
	// -----
	// err = Marked errors occurred:
	// |	whoops
	// M	foo
	// M	bar
	// M	errorString
	// M	wrap
	// errors.Is(err, nil) = false
	// errors.Is(err, errorStringErr) = true
	// errors.As(err, &asWrapErr) = true
	// wrapErr = wrap
	// -----
	// err = Marked errors occurred:
	// |	Marked errors occurred:
	// |	whoops
	// M	foo
	// M	bar
	// M	errorString
	// M	wrap
	// M	alice
	// M	bob
	// M	wrap
	// errors.Is(err, nil) = false
	// errors.Is(err, errorStringErr) = true
	// errors.As(err, &asWrapErr) = true
	// wrapErr = wrap
	// -----
}

type wrapError struct {
	msg string
}

func (e *wrapError) Error() string {
	return e.msg
}
