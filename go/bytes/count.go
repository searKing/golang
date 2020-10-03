package bytes

import "bytes"

// CountIndex counts the number of non-overlapping instances of sep in s.
// Index returns the index of the first instance of sep in s, or -1 if sep is not present in s.
func CountIndex(s, sep []byte) (c, index int) {
	n := 0
	lastIndex := -1
	for {
		i := bytes.Index(s, sep)
		if i == -1 {
			return n, lastIndex
		}
		n++
		lastIndex = i
		s = s[i+len(sep):]
	}
}
