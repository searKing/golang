package slice

// TakeUntilFunc returns, if this slice is ordered, a slice consisting of the longest
// prefix of elements taken from this slice that unmatch the given predicate.
// Otherwise returns, if this slice is unordered, a slice consisting of a
// subset of elements taken from this slice that unmatch the given predicate.
func TakeUntilFunc(s interface{}, f func(interface{}) bool) interface{} {
	return normalizeSlice(takeUntilFunc(Of(s), f, false), s)
}

// takeUntilFunc is the same as TakeUntilFunc.
func takeUntilFunc(s []interface{}, f func(interface{}) bool, truth bool) []interface{} {
	return takeWhileFunc(s, f, !truth)
}
