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
	caller := runtime.GetCaller()
	fmt.Print(caller)

	// Output:
	// github.com/searKing/golang/go/runtime_test.ExampleGetCaller
}

//go:nosplit
func sp() (got, expect uintptr) {
	var sp uintptr
	sp = runtime.GetSP(0)
	return sp, uintptr(unsafe.Pointer(&sp))
}

//func ExampleGetSP() {
//	got, expect := sp()
//	if got != expect {
//		fmt.Printf("got = %#x\n", got)
//		fmt.Printf("expect = %#x\n", expect)
//	}
//	fmt.Printf("%t", got == expect)
//
//	// Output:
//	// true
//}

//go:nosplit
func spv() (got, expect int) {
	var spv = 1
	var sp *int
	sp = (*int)(unsafe.Pointer(runtime.GetSP(0)))
	return *sp, spv
}

//func ExampleGetSPV() {
//	got, expect := spv()
//	if got != expect {
//		fmt.Printf("got = %d\n", got)
//		fmt.Printf("expect = %d\n", expect)
//	}
//	fmt.Printf("%t", got == expect)
//
//	// Output:
//	// true
//}
