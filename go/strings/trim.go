package strings

import (
	"fmt"
	"io"
)

// SplitPrefixNumber slices s into number prefix and unparsed and
// returns a slice of those substrings.
// If s does not start with number, SplitPrefixNumber returns
// a slice of length 1 whose only element is s.
func SplitPrefixNumber(s string) []string {
	unparsed := TrimPrefixNumber(s)
	return []string{s[:len(s)-len(unparsed)], unparsed}
}

// TrimPrefixNumber returns s without the leading number prefix string.
// If s doesn't start with number prefix, s is returned unchanged.
func TrimPrefixNumber(s string) string {
	unparsedFloat := TrimPrefixFloat(s)
	unparsedInt := TrimPrefixInteger(s)
	if len(unparsedFloat) < len(unparsedInt) {
		return unparsedFloat
	}
	return unparsedInt
}

// TrimPrefixFloat returns s without the leading float prefix string.
// If s doesn't start with float prefix, s is returned unchanged.
func TrimPrefixFloat(s string) string {
	var value float64
	var unparsed string
	count, err := fmt.Sscanf(s, `%v%s`, &value, &unparsed)

	if (err != nil && err != io.EOF) || (count == 0) {
		return s
	}
	return unparsed
}

// TrimPrefixInteger returns s without the leading integer prefix string.
// If s doesn't start with integer prefix, s is returned unchanged.
func TrimPrefixInteger(s string) string {
	var value int64
	var unparsed string
	count, err := fmt.Sscanf(s, `%v%s`, &value, &unparsed)

	if (err != nil && err != io.EOF) || (count == 0) {
		return s
	}
	return unparsed
}

// TrimPrefixComplex returns s without the leading complex prefix string.
// If s doesn't start with complex prefix, s is returned unchanged.
func TrimPrefixComplex(s string) string {
	var value complex128
	var unparsed string
	count, err := fmt.Sscanf(s, `%v%s`, &value, &unparsed)

	if (err != nil && err != io.EOF) || (count == 0) {
		return s
	}
	return unparsed
}
