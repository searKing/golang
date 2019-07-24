package slice

import (
	"github.com/searKing/golang/go/util/object"
)

// LimitFunc Returns a slice consisting of the elements of this slice, truncated
// to be no longer than {@code maxSize} in length.
func LimitFunc(s interface{}, maxSize int) interface{} {
	return normalizeSlice(limitFunc(Of(s), maxSize), s)
}

// limitFunc is the same as LimitFunc.
func limitFunc(s []interface{}, maxSize int) []interface{} {
	object.RequireNonNil(s, "limitFunc called on nil slice")
	m := len(s)
	if m > maxSize {
		m = maxSize
	}
	return s[:m]
}
