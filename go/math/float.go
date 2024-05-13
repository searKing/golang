// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math

import "math"

// Epsilon very small
var Epsilon = 1e-6
var EpsilonClose = 1e-14
var EpsilonVeryClose = 1e-16

// TruncPrecision returns the float value of x, with
// case n >= 0
//
//	the maximum n bits precision.
//
// case n < 0
//
//	-n bits of the magnitude of x trunked
//
// Special cases are:
//
//	Trunc(±0) = ±0
//	Trunc(±Inf) = ±Inf
//	Trunc(NaN) = NaN
func TruncPrecision(x float64, n int) float64 {
	n10 := math.Pow10(n)
	return math.Copysign(math.Trunc((math.Abs(x)+0.5/n10)*n10)/n10, x)
}

// Tolerance returns true if |a-b| < e
// e usually can be set with  Epsilon(1e-6)
func Tolerance(a, b, e float64) bool {
	// Multiplying by e here can underflow denormal values to zero.
	// Check a==b so that at least if a and b are small and identical
	// we say they match.
	if a == b {
		return true
	}
	d := a - b
	if d < 0 {
		d = -d
	}

	// note: b is correct (expected) value, a is actual value.
	// make error tolerance a fraction of b, not a.
	if b != 0 {
		e = e * b
		if e < 0 {
			e = -e
		}
	}
	return d < e
}

// Close returns true if |a-b| < 1e14
func Close(a, b float64) bool { return Tolerance(a, b, EpsilonClose) }

// VeryClose returns true if |a-b| < 1e16
func VeryClose(a, b float64) bool { return Tolerance(a, b, EpsilonVeryClose) }

// SoClose is an alias of Tolerance
func SoClose(a, b, e float64) bool { return Tolerance(a, b, e) }

// Alike returns true if a,b is the same exactly (no tolerance) or both NaN
func Alike(a, b float64) bool {
	switch {
	case math.IsNaN(a) && math.IsNaN(b):
		return true
	case a == b:
		return math.Signbit(a) == math.Signbit(b)
	}
	return false
}
