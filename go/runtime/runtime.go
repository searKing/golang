// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	_ "runtime"
	"unsafe"

	"github.com/searKing/golang/go/reflect"
)

// ReadNextSP returns the value, that is SP pointed to where the caller
// value Next SP pointed to will be reset randomly when ReadNextSP returns
//go:linkname ReadNextSP github.com/searKing/golang/go/runtime.readNextSP
func ReadNextSP() int

// readNextSP returns the value, that is SP pointed to where the caller
// value Next SP pointed to will be reset randomly when ReadNextSP returns
//go:nosplit
//go:noinline
func readNextSP() (x int) {
	return
}

// GetSP returns the location, that is SP where the caller
//go:linkname GetSP github.com/searKing/golang/go/runtime.getSP
//go:nosplit
//go:noinline
func GetSP(uintptr) uintptr

// getSP returns the location, that is SP where the caller
//go:nosplit
//go:noinline
func getSP(x uintptr) uintptr {
	// x is an argument mainly so that we can return its address.
	//return uintptr(noescape(unsafe.Pointer(&x))) + unsafe.Sizeof(struct{ _, _ uintptr }{})
	//return uintptr(noescape(unsafe.Pointer(&x))) + unsafe.Sizeof(uintptr(0))*2
	return uintptr(noescape(unsafe.Pointer(&x))) + reflect.PtrSize*2
}

// getNextSP returns the location, that is SP after PUSH 0 where the caller
//go:nosplit
//go:noinline
func getNextSP() (x uintptr) {
	// x is an argument mainly so that we can return its address.
	return uintptr(noescape(unsafe.Pointer(&x)))
}

// noescape hides a pointer from escape analysis.  noescape is
// the identity function but escape analysis doesn't think the
// output depends on the input.  noescape is inlined and currently
// compiles down to zero instructions.
// USE CAREFULLY!
//go:nosplit
func noescape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	return unsafe.Pointer(x ^ 0)
}
