package bytes

import (
	"bytes"
)

// Truncate shrinks s's len to n at most
func Truncate(s []byte, n int) []byte {
	if n < 0 {
		n = 0
	}
	if len(s) <= n {
		return s
	}
	return s[:n]
}

// PadLeft returns s padded to length n, padded left with repeated pad
// return s directly if pad is empty
// padding s with {{pad}} and spaces(less than len(pad)) as a prefix, as [pad]...[pad][space]...[space][s]
func PadLeft(s []byte, pad []byte, n int) []byte {
	if len(pad) == 0 {
		return s
	}

	pc, sc := ComputePad(s, pad, n)

	return append(bytes.Repeat(pad, pc), append(bytes.Repeat([]byte(" "), sc), s...)...)
}

// PadRight returns s padded to length n, padded right with repeated pad
// return s directly if pad is empty
// padding s with {{pad}} and spaces(less than len(pad))  as a suffix, as [s][space]...[space][pad]...[pad]
func PadRight(s []byte, pad []byte, n int) []byte {
	if len(pad) == 0 {
		return s
	}
	pc, sc := ComputePad(s, pad, n)

	return append(append(s, bytes.Repeat([]byte(" "), sc)...), bytes.Repeat(pad, pc)...)
}

// ComputePad returns pad's count and space's count(less than len(pad)) will be need to pad s to len n
// padCount = (n-len(s))/len(pad)
// spaceCount = (n-len(s))%len(pad)
func ComputePad(s []byte, pad []byte, n int) (padCount, spaceCount int) {
	if len(pad) == 0 {
		return 0, 0
	}

	c := n - len(s)
	if c < 0 {
		c = 0
	}

	padCount = c / len(pad)

	spaceCount = c % len(pad)
	return padCount, spaceCount
}

// Reverse returns bytes in reverse order.
func Reverse(s []byte) []byte {
	var b bytes.Buffer
	b.Grow(len(s))
	for i := len(s) - 1; i >= 0; i-- {
		b.WriteByte(s[i])
	}
	return b.Bytes()
}
