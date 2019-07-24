package slice

// EmptyFunc returns an empty sequential {@code slice}.
func EmptyFunc(s interface{}) interface{} {
	return normalizeSlice(emptyFunc(), s)
}

// emptyFunc is the same as EmptyFunc
func emptyFunc() []interface{} {
	return []interface{}{}
}
