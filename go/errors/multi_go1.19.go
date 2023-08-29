// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !go1.20

package errors

import (
	"errors"
	"strings"
)

var _ error = (*multiError)(nil) // verify that Error implements error

// Multi returns an error that wraps the given errors.
// Any nil error values are discarded.
// Join returns nil if every value in errs is nil.
// The error formats as the concatenation of the strings obtained
// by calling the Error method of each element of errs, with a '|'
// between each string.
// Deprecated: Use errors.Join instead since go1.20.
func Multi(errs ...error) error {
	n := 0
	for _, err := range errs {
		if err != nil {
			n++
		}
	}
	if n == 0 {
		return nil
	}
	me := make(multiError, 0, n)
	for _, err := range errs {
		if err != nil {
			me = append(me, err)
		}
	}
	return me
}

type multiError []error

func (e multiError) Error() string {
	var b strings.Builder
	for i, err := range e {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString(err.Error())
	}
	return b.String()
}

// Unwrap returns the error in e, if there is exactly one. If there is more than one
// error, Unwrap returns first non-nil error , since there is no way to determine which should be
// returned.
// Deprecated: Unwrap() []error supported instead since go1.20.
func (e multiError) Unwrap() error {
	for _, err := range e {
		if err != nil {
			return err
		}
	}
	// Return nil when e is nil, or has more than one error.
	// When there are multiple errors, it doesn't make sense to return any of them.
	return nil
}

// Is reports whether any error in multiError matches target.
// Deprecated: Unwrap() []error supported instead since go1.20.
func (e multiError) Is(target error) bool {
	if target == nil {
		for _, err := range e {
			if err != nil {
				return true
			}
		}
		return false
	}
	for _, err := range e {
		if err == nil {
			continue
		}
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

// As finds the first error in err's chain that matches target, and if one is found, sets
// target to that error value and returns true. Otherwise, it returns false.
// Deprecated: Unwrap() []error supported instead since go1.20.
func (e multiError) As(target any) bool {
	for _, err := range e {
		if errors.As(err, target) {
			return true
		}
	}
	return false
}
