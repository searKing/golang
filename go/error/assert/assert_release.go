// +build !debug

package assert

// Assert panics ignored.
func Assert(cond bool, a ...interface{}) bool {
	return cond
}

// Assertln panics ignored.
func Assertln(cond bool, a ...interface{}) bool {
	return cond
}

// Assertf panics ignored.
func Assertf(cond bool, format string, a ...interface{}) bool {
	return cond
}
