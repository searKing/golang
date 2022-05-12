// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	"errors"
	"fmt"
	"io"
)

var _ error = multiError{} // verify that Error implements error

// Multi returns an error with the supplied errors.
// If no error contained, return nil.
func Multi(errs ...error) error {
	me := multiError(errs).clean()
	if me == nil || len(me) == 0 {
		return nil
	}
	return me
}

type multiError []error

func (e multiError) Error() string {
	errs := e.clean()
	if errs == nil || len(errs) == 0 {
		return ""
	}
	message := errs[0].Error()
	for _, err := range errs[1:] {
		message += "|" + err.Error()
	}

	return message
}

func (e multiError) Format(s fmt.State, verb rune) {
	errs := e.clean()
	if errs == nil || len(errs) == 0 {
		return
	}
	switch verb {
	case 'v':
		if s.Flag('+') {
			if len(errs) == 1 {
				_, _ = fmt.Fprintf(s, "%+v", errs[0])
				return
			}
			_, _ = io.WriteString(s, "Multiple errors occurred:\n")

			_, _ = fmt.Fprintf(s, "|\t%+v", errs[0])
			for _, err := range errs[1:] {
				_, _ = fmt.Fprintf(s, "\n|\t%+v", err)
			}
			return
		}
		fallthrough
	case 's', 'q':
		_, _ = io.WriteString(s, errs.Error())
	}
}

// clean removes all none nil elem in all the errors
func (e multiError) clean() multiError {
	var errs []error
	for _, err := range e {
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

// Is reports whether any error in multiError matches target.
func (e multiError) Is(target error) bool {
	if target == nil {
		errs := e.clean()
		if errs == nil || len(errs) == 0 {
			return true
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

// Unwrap returns the error in e, if there is exactly one. If there is more than one
// error, Unwrap returns nil, since there is no way to determine which should be
// returned.
func (e multiError) Unwrap() error {
	if len(e) == 1 {
		return e[0]
	}
	// Return nil when e is nil, or has more than one error.
	// When there are multiple errors, it doesn't make sense to return any of them.
	return nil
}

// As finds the first error in err's chain that matches target, and if one is found, sets
// target to that error value and returns true. Otherwise, it returns false.
func (e multiError) As(target any) bool {
	for _, err := range e {
		if errors.As(err, target) {
			return true
		}
	}
	return false
}
