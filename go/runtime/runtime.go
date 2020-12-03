package runtime

import "unsafe"

// Code borrowed from https://github.com/golang/go/blob/master/src/runtime/panic.go#L1068.

// noescape hides a pointer from escape analysis.  noescape is
// the identity function but escape analysis doesn't think the
// output depends on the input.  noescape is inlined and currently
// compiles down to zero instructions.
// USE CAREFULLY!
// 禁止逃逸,即将指针变成无意义整数
//go:nosplit
func NoEscape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	return unsafe.Pointer(x ^ 0)
}

// getargp returns the location where the caller
// writes outgoing function call arguments.
//go:nosplit
//go:noinline
func GetArgPtr(x int) uintptr {
	return uintptr(NoEscape(unsafe.Pointer(&x)))
}
