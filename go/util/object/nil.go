// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package object

import (
	"errors"
	"reflect"
	"strings"
)

type ErrorNilPointer error
type ErrorMissMatch error

var (
	errorNilPointer = errors.New("nil pointer")
)

// IsNil returns {@code true} if the provided reference is {@code nil} otherwise
// returns {@code false}.
func IsNil(obj interface{}) bool {
	if obj == nil {
		return true
	}
	if IsNilable(obj) {
		return reflect.ValueOf(obj).IsNil()
	}
	return false
}

// IsNil returns {@code true} if the provided reference is non-{@code nil} otherwise
// returns {@code false}.
func NoneNil(obj interface{}) bool {
	return !IsNil(obj)
}

// IsNil returns {@code true} if the provided reference can be assigned {@code nil} otherwise
// returns {@code false}.
func IsNilable(obj interface{}) (canBeNil bool) {
	defer func() {
		// As we can not access v.flag&reflect.flagMethod&v.ptr
		// So we use recover() instead
		if r := recover(); r != nil {
			canBeNil = false
		}
	}()
	reflect.ValueOf(obj).IsNil()

	canBeNil = true
	return
}

// RequireNonNil checks that the specified object reference is not {@code nil}. This
// method is designed primarily for doing parameter validation in methods
// and constructors
func RequireNonNil(obj interface{}, msg ...string) interface{} {
	if msg == nil {
		msg = []string{"nil pointer"}
	}
	if IsNil(obj) {
		panic(ErrorNilPointer(errors.New(strings.Join(msg, ""))))
	}
	return obj
}

// grammer surgar for RequireNonNil
func RequireNonNull(obj interface{}, msg ...string) interface{} {
	return RequireNonNil(obj, msg...)
}

// RequireNonNullElse returns the first argument if it is non-{@code nil} and
// otherwise returns the non-{@code nil} second argument.
func RequireNonNullElse(obj, defaultObj interface{}) interface{} {
	if NoneNil(obj) {
		return obj
	}
	return RequireNonNil(defaultObj, "defaultObj")
}

// RequireNonNullElseGet returns the first argument if it is non-{@code nil} and
// returns the non-{@code nil} value of {@code supplier.Get()}.
func RequireNonNullElseGet(obj interface{}, sup interface{ Get() interface{} }) interface{} {
	if NoneNil(obj) {
		return obj
	}
	return RequireNonNil(RequireNonNil(sup, "supplier").(interface{ Get() interface{} }).Get(), "supplier.Get()")
}

// RequireNonNil checks that the specified object reference is not {@code nil}. This
// method is designed primarily for doing parameter validation in methods
// and constructors
func RequireEqual(actual, expected interface{}, msg ...string) interface{} {
	if msg == nil {
		msg = []string{"miss match"}
	}
	if !Equals(actual, expected) {
		panic(ErrorMissMatch(errors.New(strings.Join(msg, ""))))
	}
	return actual
}
