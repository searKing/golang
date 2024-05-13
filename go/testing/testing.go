// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testing

import (
	"fmt"
	"reflect"

	reflect_ "github.com/searKing/golang/go/reflect"
)

// Nil asserts that the specified object is nil.
func Nil(v any, a ...any) (bool, string) {
	return reflect_.IsNil(v), fmt.Sprintf(fmt.Sprintf("expected nil, got: %#v; ", v), a...)
}

func Nilf(v any, format string, a ...any) (bool, string) {
	return reflect_.IsNil(v), fmt.Sprintf(fmt.Sprintf("expected nil, got: %#v; %s", v, format), a...)
}

// NonNil asserts that the specified object is none nil.
func NonNil(v any, a ...any) (bool, string) {
	return !reflect_.IsNil(v), fmt.Sprintf(fmt.Sprintf("expected non-nil, got: %#v; ", v), a...)
}

func NonNilf(v any, format string, a ...any) (bool, string) {
	return !reflect_.IsNil(v), fmt.Sprintf(fmt.Sprintf("expected non-nil, got: %#v; %s", v, format), a...)
}

// Zero asserts that the specified object is zero.  I.e. nil, "", false, 0 or either
// a slice or a channel with len == 0.
func Zero(v any, a ...any) (bool, string) {
	return reflect_.IsZeroValue(reflect.ValueOf(v)), fmt.Sprintf(fmt.Sprintf("expected zero value, got: %#v; ", v), a...)
}

func Zerof(v any, format string, a ...any) (bool, string) {
	return reflect_.IsZeroValue(reflect.ValueOf(v)), fmt.Sprintf(fmt.Sprintf("expected zero value, got: %#v; %s", v, format), a...)
}

// NonZero asserts that the specified object is none zero.
func NonZero(v any, a ...any) (bool, string) {
	return !reflect_.IsZeroValue(reflect.ValueOf(v)), fmt.Sprintf(fmt.Sprintf("expected non-zero value, got: %#v; ", v), a...)
}

func NonZerof(v any, format string, a ...any) (bool, string) {
	return !reflect_.IsZeroValue(reflect.ValueOf(v)), fmt.Sprintf(fmt.Sprintf("expected non-zero value, got: %#v; %s", v, format), a...)
}

// Error asserts that a function returned an error (i.e. not `nil`).
func Error(v any, a ...any) (bool, string) {
	return NonNil(v, a...)
}

func Errorf(v any, format string, a ...any) (bool, string) {
	return NonNilf(v, format, a...)
}

// NonError asserts that a function returned a none error (i.e. `nil`).
func NonError(v any, a ...any) (bool, string) {
	return Nil(v, a...)
}

func NonErrorf(v any, format string, a ...any) (bool, string) {
	return Nilf(v, format, a...)
}

// EqualError asserts that a function returned an error (i.e. not `nil`)
// and that it is equal to the provided error.
func EqualError(actual error, expected error, a ...any) (bool, string) {
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

func EqualErrorf(actual error, expected error, format string, a ...any) (bool, string) {
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
