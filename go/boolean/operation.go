// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package boolean

import (
	"github.com/searKing/golang/go/container/slice"
	"github.com/searKing/golang/go/util/object"
)

// https://en.wikipedia.org/wiki/Boolean_operation

func xor(a, b bool) bool {
	return a != b
}
func xnor(a, b bool) bool {
	return !xor(a, b)
}
func or(a, b bool) bool {
	return a || b
}
func and(a, b bool) bool {
	return a && b
}

func RelationFunc(a, b interface{}, f func(a, b interface{}) interface{}, c ...interface{}) interface{} {
	object.RequireNonNil(f)
	if c == nil || len(c) == 0 {
		return f(a, b)
	}
	return RelationFunc(f(a, b), c[0], f, c[1:]...)
}
func BoolFunc(a bool, b bool, f func(a, b bool) bool, c ...bool) bool {

	return RelationFunc(a, b, func(a, b interface{}) interface{} {
		return f(a.(bool), b.(bool))
	}, slice.Of(c)...).(bool)
	object.RequireNonNil(f)
	if c == nil || len(c) == 0 {
		return f(a, b)
	}
	return BoolFunc(f(a, b), c[0], f, c[1:]...)
}

// XOR return a^b^c...
func XOR(a bool, b bool, c ...bool) bool {
	return BoolFunc(a, b, xor, c...)
}

// XNOR return a xnor b xnor c...
func XNOR(a bool, b bool, c ...bool) bool {
	return BoolFunc(a, b, xnor, c...)
}

// OR return a|b|c...
func OR(a bool, b bool, c ...bool) bool {
	return BoolFunc(a, b, or, c...)
}

// AND return a&b&c...
func AND(a bool, b bool, c ...bool) bool {
	return BoolFunc(a, b, and, c...)
}
