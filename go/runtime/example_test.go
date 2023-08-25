// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime_test

import (
	"fmt"

	runtime_ "github.com/searKing/golang/go/runtime"
)

func ExampleGetCaller() {
	caller := runtime_.GetCaller(1)
	fmt.Print(caller)

	// Output:
	// github.com/searKing/golang/go/runtime_test.ExampleGetCaller
}

func ExampleGetShortCallerFuncFileLine() {
	caller, file, line := runtime_.GetShortCallerFuncFileLine(1)
	fmt.Printf("%s() %s:%d", caller, file, line)

	// Output:
	// ExampleGetShortCallerFuncFileLine() example_test.go:22
}
