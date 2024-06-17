// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package unsafe

import "unsafe"

// Bytes returns a byte slice whose underlying array starts at ptr of v
// and whose length and capacity are len(v).
// Bytes(v) is equivalent to
//
//	(*[unsafe.Sizeof(v)]byte)(unsafe.Pointer(&v))[:]
//
// except that, as a special case, if ptr is nil and len is zero,
// Slice returns nil.
func Bytes[T any](v T) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(&v)), unsafe.Sizeof(v))
}

// Placement returns a T pointer whose underlying bytes
// construct objects in allocated storage.
//
// Construct a “T” object, placing it directly into your
// pre-allocated storage at memory address “s”.
//
// except that, as a special case, if s is nil or len is zero,
// Placement returns nil.
func Placement[T any](s []byte) *T {
	if len(s) == 0 {
		return nil
	}
	return (*T)(unsafe.Pointer(&s[0]))
}
