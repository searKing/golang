package slice

import (
	"github.com/searKing/golang/go/util/object"
	"sort"
)

// SortedFunc returns a slice consisting of the distinct elements (according to
// {@link Object#equals(Object)}) of this slice.
// s: Accept Array、Slice、String(as []byte if ifStringAsRune else []rune)
func SortedFunc(s interface{}, f func(interface{}, interface{}) int) interface{} {
	return normalizeSlice(sortedFunc(Of(s), f), s)
}

// sortedFunc is the same as SortedFunc except that if
// truth==false, the sense of the predicate function is
// inverted.
func sortedFunc(s []interface{}, f func(interface{}, interface{}) int) []interface{} {
	object.RequireNonNil(s, "distinctFunc called on nil slice")
	object.RequireNonNil(f, "distinctFunc called on nil callfn")

	less := func(i, j int) bool {
		if f(s[i], s[j]) < 0 {
			return true
		}
		return false
	}
	sort.Slice(s, less)
	return s
}
