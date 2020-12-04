// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	_ "runtime"
	"unsafe"

	"github.com/searKing/golang/go/reflect"
)

// GetEIP returns the location, that is EIP after CALL
//go:linkname GetEIP github.com/searKing/golang/go/runtime.getEIP
//go:nosplit
//go:noinline
func GetEIP(uintptr) uintptr

// getEIP returns the location, that is EIP after CALL
// -> arg+argsize-1(FP)
// arg includes returns and arguments
// call frame stack <-> argsize+tmpsize+framesize
// tmp is for EIP AND EBP
//go:nosplit
//go:noinline
func getEIP(x uintptr) uintptr {
	// x is an argument mainly so that we can return its address.
	// plus reflect.PtrSize *2 for shrink call frame to zero, that is EIP
	// ATTENTION NO BSP ON STACK FOR NO SUB FUNC CALL IN THIS FUNCTION, so plus 1: EIP,VAR
	return uintptr(noescape(unsafe.Pointer(&x))) + reflect.PtrSize + x
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
