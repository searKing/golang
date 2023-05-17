// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math

import (
	"github.com/searKing/golang/go/exp/types"
	"golang.org/x/exp/constraints"
)

// SerialNumber is a constraint that permits any unsigned integer type.
// If future releases of Go add new predeclared unsigned integer types,
// this constraint will be modified to include them.
type SerialNumber = constraints.Unsigned

// IsNewer implements RFC 1982: Serial Number Arithmetic
// See also: https://datatracker.ietf.org/doc/html/rfc1982#section-2
// s1 < s2 and (s1 + 1) > (s2 + 1)
// See also: https://chromium.googlesource.com/external/webrtc/trunk/webrtc/+/f54860e9ef0b68e182a01edc994626d21961bc4b/modules/include/module_common_types.h
func IsNewer[T SerialNumber](t1, t2 T) bool {
	// kBreakpoint is the half-way mark for the type U. For instance, for a
	// uint16_t it will be 0x8000, and for a uint32_t, it will be 0x8000000.
	kBreakpoint := (MaxInt[T]() >> 1) + 1
	if t1 == t2 {
		return false
	}
	// Distinguish between elements that are exactly kBreakpoint apart.
	// If t1>t2 and |t1-t2| = kBreakpoint: IsNewer(t1,t2)=true,
	// IsNewer(t2,t1)=false
	// rather than having IsNewer(t1,t2) = IsNewer(t2,t1) = false.
	if t1 > t2 {
		return t1-t2 <= kBreakpoint
	}

	return t2-t1 > kBreakpoint
}

// Latest return newer of s1, s2
func Latest[T SerialNumber](s1, s2 T) T {
	if IsNewer(s1, s2) {
		return s1
	}
	return s2
}

// Unwrap unwrap a number to a larger type.
// The numbers will never be unwrapped to a negative value.
func Unwrap[T SerialNumber](lastValue int64, value T) int64 {
	kMaxPlusOne := int64(MaxInt[T]()) + 1
	croppedLast := T(lastValue)
	delta := int64(value - croppedLast)
	if IsNewer(value, croppedLast) {
		if delta < 0 {
			delta += kMaxPlusOne // Wrap forwards.
		}
	} else if delta > 0 && (lastValue+delta-kMaxPlusOne) >= 0 {
		// If value is older but delta is positive, this is a backwards
		// wrap-around. However, don't wrap backwards past 0 (unwrapped).
		delta -= kMaxPlusOne
	}
	return lastValue + delta
}

// Unwrapper is an utility class to unwrap a number to a larger type.
// The numbers will never be unwrapped to a negative value.
type Unwrapper[T SerialNumber] struct {
	lastValue *int64
}

// UnwrapWithoutUpdate returns the unwrapped value, but don't update the internal state.
func (u *Unwrapper[T]) UnwrapWithoutUpdate(value T) int64 {
	if u.lastValue == nil {
		return int64(value)
	}
	kMaxPlusOne := int64(MaxInt[T]()) + 1
	croppedLast := T(*u.lastValue)
	delta := int64(value - croppedLast)
	if IsNewer(value, croppedLast) {
		if delta < 0 {
			delta += kMaxPlusOne // Wrap forwards.
		}
	} else if delta > 0 && (*u.lastValue+delta-kMaxPlusOne) >= 0 {
		// If value is older but delta is positive, this is a backwards
		// wrap-around. However, don't wrap backwards past 0 (unwrapped).
		delta -= kMaxPlusOne
	}
	return *u.lastValue + delta
}

// UpdateLast only update the internal state to the specified last (unwrapped) value.
func (u *Unwrapper[T]) UpdateLast(lastValue int64) {
	u.lastValue = types.Pointer(lastValue)
}

// Unwrap the value and update the internal state.
func (u *Unwrapper[T]) Unwrap(value T) int64 {
	unwrapped := u.UnwrapWithoutUpdate(value)
	u.UpdateLast(unwrapped)
	return unwrapped
}
