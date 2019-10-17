package math

import "math"

// TruncPrecision returns the float value of x, with
// case n >= 0
// 	the maximum n bits precision.
// case n < 0
//	-n bits of the magnitude of x trunked
// Special cases are:
//	Trunc(±0) = ±0
//	Trunc(±Inf) = ±Inf
//	Trunc(NaN) = NaN
func TruncPrecision(x float64, n int) float64 {
	n10 := math.Pow10(n)
	return math.Copysign(math.Trunc((x+0.5/n10)*n10)/n10, x)
}

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
func Close(a, b float64) bool { return Tolerance(a, b, 1e-14) }

func VeryClose(a, b float64) bool { return Tolerance(a, b, 4e-16) }

func SoClose(a, b, e float64) bool { return Tolerance(a, b, e) }

func Alike(a, b float64) bool {
	switch {
	case math.IsNaN(a) && math.IsNaN(b):
		return true
	case a == b:
		return math.Signbit(a) == math.Signbit(b)
	}
	return false
}
