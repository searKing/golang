// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reflect

import "reflect"

// IsNil reports whether its argument v is a nil interface value or an untyped nil.
// Note that IsNil is not always equivalent to a regular comparison with nil in Go.
// It is equivalent to:
//
//	var typedNil any = (v's underlying type)(nil)
//	return v == nil || v == typedNil
//
// For example, if v was created by set with `var p *int` or calling IsNil((*int)(nil)),
// i==nil will be false but [IsNil] will be true.
func IsNil(v any) bool {
	if v == nil {
		return true
	}
	return IsNilValue(reflect.ValueOf(v))
}

// UnTypeNil returns its argument v or nil if and only if v is a nil interface value or an untyped nil.
func UnTypeNil(v any) any {
	// convert typed nil into untyped nil
	if IsNil(v) {
		return nil
	}
	return v
}

// Deprecated: Use [IsNil] instead.
func IsNilObject(i any) (result bool) {
	return IsNil(i)
}
