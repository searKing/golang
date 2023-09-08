// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"unsafe"

	"github.com/searKing/golang/go/exp/constraints"
)

// Conv converts a number[T] into an unlimited untyped number, and assigned back to number[S].
// static_cast<S>(v)
// https://go.dev/ref/spec#Constants
func Conv[T, S constraints.Number](v T) S {
	return S(v)
}

// ConvBinary writes a number[T] into a byte slice, and reads the byte slice into a number[S].
// binary.Write( T -> []byte ) => binary.Read( []byte, S )
func ConvBinary[T, S any](v T) S {
	return *(*S)(unsafe.Pointer(&v))
}
