// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"unsafe"

	"github.com/searKing/golang/go/reflect"
)

// GetEIP returns the location, that is EIP after CALL
// -> arg+argsize-1(FP)
// arg includes returns and arguments
// call frame stack <-> argsize+tmpsize+framesize
// tmp is for EIP AND EBP
//
// Deprecated: GetEIP returns an address on the stack rather than EIP since the register-based
// calling convention of Go 1.17, and the stack may move at any time.
// Use [runtime.Caller] or [runtime.Callers] to get program counters instead.
//
//go:nosplit
//go:noinline
func GetEIP(x uintptr) uintptr {
	// x is an argument mainly so that we can return its address.
	// plus reflect.PtrSize *2 for shrink call frame to zero, that is EIP
	// ATTENTION NO BSP ON STACK FOR NO SUB FUNC CALL IN THIS FUNCTION, so plus 1: EIP,VAR
	return uintptr(noescape(unsafe.Pointer(&x))) + reflect.PtrSize + x
}

// noescape hides a pointer from escape analysis. It is the identity function
// but escape analysis doesn't think the output depends on the input.
// noescape is inlined and currently compiles down to zero instructions.
// USE CAREFULLY!
// This was copied from the runtime; see issues 23382 and 7921.
//
//go:nosplit
//go:nocheckptr
func noescape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	return unsafe.Pointer(x ^ 0)
}
