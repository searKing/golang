// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime_test

import (
	"fmt"
	"unsafe"

	"github.com/searKing/golang/go/runtime"
)

func ExampleGetCaller() {
	caller := runtime.GetCaller(1)
	fmt.Print(caller)

	// Output:
	// github.com/searKing/golang/go/runtime_test.ExampleGetCaller
}

func ExampleGetShortCallerFuncFileLine() {
	caller, file, line := runtime.GetShortCallerFuncFileLine(1)
	fmt.Printf("%s() %s:%d", caller, file, line)

	// Output:
	// ExampleGetShortCallerFuncFileLine() run_example.go:63
}

//go:nosplit
func cf() (got, expect uintptr) {
	var spv = 99
	var sp *int
	sp = (*int)(unsafe.Pointer(runtime.GetEIP(unsafe.Sizeof(uintptr(0)))))
	//sp = (*int)(unsafe.Pointer(runtime.GetEIP(0)))
	*sp = 0xFFFFFFF

	return uintptr(unsafe.Pointer(sp)), uintptr(unsafe.Pointer(&spv))
}
func ExampleGetCallFrame() {
	got, expect := cf()
	if got != expect {
		fmt.Printf("got = %#x\n", got)
		fmt.Printf("expect = %#x\n", expect)
	}
	fmt.Printf("%t", got == expect)

	// Output:
	// true

}
