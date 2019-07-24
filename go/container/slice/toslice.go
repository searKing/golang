package slice

// ToSliceFunc returns an array containing the elements of this stream.
func ToSliceFunc(s interface{}) interface{} {
	return toSliceFunc(Of(s))
}

// toSliceFunc is the same as ToSliceFunc
func toSliceFunc(s []interface{}) []interface{} {
	return s
}
