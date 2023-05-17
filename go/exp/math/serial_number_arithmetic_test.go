// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math_test

import (
	"fmt"
	"testing"
	"testing/quick"

	math_ "github.com/searKing/golang/go/exp/math"
	"golang.org/x/exp/constraints"
)

func TestIsNewer(t *testing.T) {
	tests := []struct {
		x, y uint8
		want bool
	}{
		{0, 0, false},
		{0xFF, 0xFF, false},
		{0, 1, false},
		{0xFE, 0xFF, false},
		{0, 0xFF, true},
		{0xFF, 0, false},
		{0, 0x7F, false},
		{0, 0x7E, false},
		{0, 0x80, false},
		{1, 0, true},
		{44, 0, true},
		{100, 0, true},
		{100, 44, true},
		{200, 100, true},
		{255, 200, true},
		{0, 255, true},
		{100, 255, true},
		{0, 200, true},
		{44, 200, true},
		{126, 255, true},
		{127, 255, false},
		{125, 254, true},
		{126, 254, false},
		{58, 25, true},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("math_.IsNewer(%v, %v)", tt.x, tt.y), func(t *testing.T) {
			{
				got := math_.IsNewer(tt.x, tt.y)
				if got != tt.want {
					t.Errorf("math_.IsNewer(%v, %v) = %v, want %v", tt.x, tt.y, got, tt.want)
				}
			}
		})
	}

	if err := quick.CheckEqual(math_.IsNewer[uint8], checkIsNewer[uint8], nil); err != nil {
		t.Error(err)
	}
}

func TestIsNewerUint64(t *testing.T) {
	tests := []struct {
		x, y uint64
		want bool
	}{
		{0, 0, false},
		{0xFFFFFFFF, 0xFFFFFFFF, false},
		{0, 1, false},
		{0xFFFFFFFE, 0xFFFFFFFF, false},
		{0, 0xFFFFFFFFFFFFFFFF, true},
		{0xFFFFFFFFFFFFFFFF, 0, false},
		{0, 0x7FFFFFFFFFFFFFFF, false},
		{0, 0x7FFFFFFFFFFFFFFE, false},
		{0, 0x8000000000000000, false},
		{0x7FFFFFFFFFFFFFFE, 0xFFFFFFFFFFFFFFFF, true},
		{0x7FFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF, false},
		{0x7FFFFFFFFFFFFFFD, 0xFFFFFFFFFFFFFFFE, true},
		{0x7FFFFFFFFFFFFFFE, 0xFFFFFFFFFFFFFFFE, false},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("math_.IsNewer(%v, %v)", tt.x, tt.y), func(t *testing.T) {
			{
				got := math_.IsNewer(tt.x, tt.y)
				if got != tt.want {
					t.Errorf("math_.IsNewer(%v, %v) = %v, want %v", tt.x, tt.y, got, tt.want)
				}
			}
		})
	}
}

func TestDistance(t *testing.T) {
	tests := []struct {
		x, y uint16
		want int16
	}{
		{0, 0, 0},
		{0x7FFF, 0, 32767},
		{0x0001, 0, 1},
		{0x0000, 0, 0},
		{0xFFFF, 0, -1},
		{0xFFFE, 0, -2},
		{0x8000, 0, -32768},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("math_.distance(%v, %v)", tt.x, tt.y), func(t *testing.T) {
			{
				got := int16(distance(tt.x, tt.y))
				if got != tt.want {
					t.Errorf("math_.distance(%v, %v) = %v, want %v", tt.x, tt.y, got, tt.want)
				}
			}
		})
	}
}

func TestUnwrap(t *testing.T) {
	tests := []struct {
		last  int64
		value uint8
		want  int64
	}{
		{0, 0, 0},
		{255, 255, 255},
		{255, 0, 256},
		{255, 1, 257},
		{126, 255, 255},
		{127, 255, 255},
		{125, 254, 254},
		{126, 254, 254},
		{125, 0, 0},
		{126, 0, 0},
		{127, 0, 0},
		{128, 0, 0},
		{129, 0, 256},
		{256, 0, 256},
		{257, 0, 256},
		{255, 0, 256},
		{256, 0, 256},
		{256, 1, 257},
		{256, 2, 258},
		{256, 3, 259},
		{512, 0, 512},
		{512, 1, 513},
		{512, 2, 514},
		{512, 3, 515},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("math_.Unwrap(%v, %v)", tt.last, tt.value), func(t *testing.T) {
			{
				got := math_.Unwrap(tt.last, tt.value)
				var u = math_.Unwrapper[uint8]{}
				u.UpdateLast(tt.last)
				got = u.Unwrap(tt.value)
				if got != tt.want {
					t.Errorf("math_.Unwrap(%v, %v) = %v, want %v", tt.last, tt.value, got, tt.want)
				}
			}
		})
	}
	if err := quick.CheckEqual(math_.Unwrap[uint8], func(lastValue int64, value uint8) int64 {
		var u = math_.Unwrapper[uint8]{}
		u.UpdateLast(lastValue)
		return u.Unwrap(value)
	}, nil); err != nil {
		t.Error(err)
	}
}

// IsNewer implements RFC 1982: Serial Number Arithmetic
// See also: https://datatracker.ietf.org/doc/html/rfc1982#section-2
// s1 < s2 and (s1 + 1) > (s2 + 1)
// See also: https://chromium.googlesource.com/external/webrtc/trunk/webrtc/+/f54860e9ef0b68e182a01edc994626d21961bc4b/modules/include/module_common_types.h
func checkIsNewer[T constraints.Unsigned](value T, preValue T) (newer bool) {
	// kBreakpoint is the half-way mark for the type U. For instance, for a
	// uint16_t it will be 0x8000, and for a uint32_t, it will be 0x8000000.
	kBreakpoint := (math_.MaxInt[T]() >> 1) + 1
	// Distinguish between elements that are exactly kBreakpoint apart.
	// If t1>t2 and |t1-t2| = kBreakpoint: IsNewer(t1,t2)=true,
	// IsNewer(t2,t1)=false
	// rather than having IsNewer(t1,t2) = IsNewer(t2,t1) = false.
	if value-preValue == kBreakpoint {
		return value > preValue
	}
	return (value != preValue) && (T(distance(value, preValue)) < kBreakpoint)
}

// distance = (signed)(i1 - i2)
// If distance is 0, the numbers are equal.
// If it is < 0, then s1 is "less than" or "before" s2.
// Simple, clean and efficient, and fully defined. However, not without surprises.
func distance[T constraints.Unsigned](s1, s2 T) int64 {
	return int64(s1) - int64(s2)
}
