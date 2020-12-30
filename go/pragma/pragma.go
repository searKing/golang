// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package pragma provides types that can be embedded into a struct to
// statically enforce or prevent certain language properties.
// The key observation and some code (shr) is borrowed from https://github.com/protocolbuffers/protobuf-go/blob/v1.25.0/internal/pragma/pragma.go
package pragma

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

// NoUnkeyedLiterals can be embedded in a struct to prevent unkeyed literals.
type NoUnkeyedLiterals struct{}

// DoNotImplement can be embedded in an interface to prevent trivial
// implementations of the interface.
//
// This is useful to prevent unauthorized implementations of an interface
// so that it can be extended in the future for any protobuf language changes.
type doNotImplement struct{}
type DoNotImplement interface{ ProtoInternal(doNotImplement) }

// DoNotCompare can be embedded in a struct to prevent comparability.
type DoNotCompare [0]func()

// DoNotCopy can be embedded in a struct to help prevent shallow copies.
// This does not rely on a Go language feature, but rather a special case
// within the vet checker.
//
// See https://golang.org/issues/8005.
type DoNotCopy [0]sync.Mutex

// CopyChecker holds back pointer to itself to detect object copying.
// Deprecated. use DoNotCopy instead, check by go vet.
// methods Copied or Check return not copied if none of methods Copied or Check have bee called before
type CopyChecker uintptr

// Copied returns true if this object is copied
func (c *CopyChecker) Copied() bool {
	return uintptr(*c) != uintptr(unsafe.Pointer(c)) &&
		!atomic.CompareAndSwapUintptr((*uintptr)(c), 0, uintptr(unsafe.Pointer(c))) &&
		uintptr(*c) != uintptr(unsafe.Pointer(c))
}

// Check panic is c is copied
func (c *CopyChecker) Check() {
	if c.Copied() {
		panic("object is copied")
	}
}
