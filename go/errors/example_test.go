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
	err1 := errors.New("err1")
	err2 := errors.New("err2")
	err3 := errors.New("err3")
	err := errors.Join(err1, err2)
	fmt.Println(err)
	if errors.Is(err, err1) {
		fmt.Println("err is err1")
	}
	if errors.Is(err, err2) {
		fmt.Println("err is err2")
	}
	if errors.Is(err, nil) {
		fmt.Println("err is nil")
	}
	if errors.Is(err, err3) {
		fmt.Println("err is err3")
	}
	// Output:
	// err1
	// err2
	// err is err1
	// err is err2
}

func ExampleMark() {
	err := errors.New("err")
	mark1 := errors.New("mark1")
	mark2 := errors.New("mark2")
	mark3 := errors.New("mark3")
	me := errors_.Mark(err, mark1, mark2)
	fmt.Println(me)
	if errors.Is(me, err) {
		fmt.Println("err is err")
	}
	if errors.Is(me, mark1) {
		fmt.Println("err is mark1")
	}
	if errors.Is(me, mark2) {
		fmt.Println("err is mark2")
	}
	if errors.Is(me, mark3) {
		fmt.Println("err is mark3")
	}
	if errors.Is(me, nil) {
		fmt.Println("err is nil")
	}
	// Output:
	// err
	// err is err
	// err is mark1
	// err is mark2
}
