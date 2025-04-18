// Copyright 2025 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package atomic

import (
	"sync/atomic"

	"github.com/searKing/golang/go/pragma"
)

// A Value is a generic wrapper for atomic.Value, which provides an atomic load and store of value.
// The zero value for a Value returns zero from [Value.Load].
// Once [Value.Store] has been called, a Value must not be copied.
//
// A Value must not be copied after first use.
type Value[T any] struct {
	_ pragma.DoNotCopy
	v atomic.Value
}

// Load returns the value set by the most recent Store.
// It returns zero if there has been no call to Store for this Value.
func (v *Value[T]) Load() T {
	value := v.v.Load()
	if value == nil {
		var zeroT T
		return zeroT
	}
	return (value).(T)
}

// Store sets the value of the [Value] v to val.
func (v *Value[T]) Store(value T) {
	v.v.Store(value)
}

// CompareAndSwap executes the compare-and-swap operation for the [Value].
func (v *Value[T]) CompareAndSwap(old, new T) (swapped bool) {
	return v.v.CompareAndSwap(old, new)
}

// Swap stores new into Value and returns the previous value. It returns zero if
// the Value is empty.
func (v *Value[T]) Swap(new T) (old T) {
	value := v.v.Swap(new)
	if value == nil {
		var zeroT T
		return zeroT
	}
	return (value).(T)
}
