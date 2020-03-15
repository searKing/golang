// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testing

import (
	"fmt"

	"github.com/searKing/golang/go/util/object"
)

// Nil asserts that the specified object is nil.
func Nil(v interface{}, a ...interface{}) (bool, string) {
	return object.IsNil(v), fmt.Sprintf(fmt.Sprintf("Expected nil, but got: %#v; ", v), a...)
}

func Nilf(v interface{}, format string, a ...interface{}) (bool, string) {
	return object.IsNil(v), fmt.Sprintf(fmt.Sprintf("Expected nil, but got: %#v; %s", v, format), a...)
}

// NonNil asserts that the specified object is none nil.
func NonNil(v interface{}, a ...interface{}) (bool, string) {
	return !object.IsNil(v), fmt.Sprintf(fmt.Sprintf("Expected non-nil, but got: %#v; ", v), a...)
}

func NonNilf(v interface{}, format string, a ...interface{}) (bool, string) {
	return !object.IsNil(v), fmt.Sprintf(fmt.Sprintf("Expected non-nil, but got: %#v; %s", v, format), a...)
}

// Zero asserts that the specified object is zero.  I.e. nil, "", false, 0 or either
// a slice or a channel with len == 0.
func Zero(v interface{}, a ...interface{}) (bool, string) {
	return object.IsZero(v), fmt.Sprintf(fmt.Sprintf("Expected zero value, but got: %#v; ", v), a...)
}

func Zerof(v interface{}, format string, a ...interface{}) (bool, string) {
	return object.IsZero(v), fmt.Sprintf(fmt.Sprintf("Expected zero value, but got: %#v; %s", v, format), a...)
}

// NonZero asserts that the specified object is none zero.
func NonZero(v interface{}, a ...interface{}) (bool, string) {
	return !object.IsZero(v), fmt.Sprintf(fmt.Sprintf("Expected non-zero value, but got: %#v; ", v), a...)
}

func NonZerof(v interface{}, format string, a ...interface{}) (bool, string) {
	return !object.IsZero(v), fmt.Sprintf(fmt.Sprintf("Expected non-zero value, but got: %#v; %s", v, format), a...)
}

// Error asserts that a function returned an error (i.e. not `nil`).
func Error(v interface{}, a ...interface{}) (bool, string) {
	return NonNil(v, a...)
}

func Errorf(v interface{}, format string, a ...interface{}) (bool, string) {
	return NonNilf(v, format, a...)
}

// NonError asserts that a function returned a none error (i.e. `nil`).
func NonError(v interface{}, a ...interface{}) (bool, string) {
	return Nil(v, a...)
}

func NonErrorf(v interface{}, format string, a ...interface{}) (bool, string) {
	return Nilf(v, format, a...)
}

// EqualError asserts that a function returned an error (i.e. not `nil`)
// and that it is equal to the provided error.
func EqualError(actual error, expected error, a ...interface{}) (bool, string) {
	if actual == nil && expected == nil {
		return true, ""
	}
	if actual != nil && expected != nil {
		return actual.Error() == expected.Error(), fmt.Sprintf(fmt.Sprintf("Error message not equal:\n"+
			"expected: %q\n"+
			"actual  : %q", expected, actual), a...)
	}
	return false, fmt.Sprintf(fmt.Sprintf("Error message not equal:\n"+
		"expected: %q\n"+
		"actual  : %q", expected, actual), a...)
}

func EqualErrorf(actual error, expected error, format string, a ...interface{}) (bool, string) {
	if actual == nil && expected == nil {
		return true, ""
	}
	if actual != nil && expected != nil {
		return actual.Error() == expected.Error(), fmt.Sprintf(fmt.Sprintf("Error message not equal:\n"+
			"expected: %q\n"+
			"actual  : %q: %s", expected, actual, format), a...)
	}

	return false, fmt.Sprintf(fmt.Sprintf("Error message not equal:\n"+
		"expected: %q\n"+
		"actual  : %q: %s", expected, actual, format), a...)
}
