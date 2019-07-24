// +build debug

package assert

import "fmt"

// Assert panics if cond is false.
func Assert(cond bool, a ...interface{}) bool {
	if !cond {
		panic(fmt.Sprint(a...))
	}
	return cond
}

// Assertln panics if cond is false.
func Assertln(cond bool, a ...interface{}) bool {
	if !cond {
		panic(fmt.Sprintln(a...))
	}
	return cond
}

// Assertf panics if cond is false.
func Assertf(cond bool, format string, a ...interface{}) bool {
	if !cond {
		panic(fmt.Sprintf(format, a...))
	}
	return cond
}
