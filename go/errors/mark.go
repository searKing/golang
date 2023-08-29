// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	"errors"
	"fmt"
	"io"
)

var _ error = (*markError)(nil) // verify that Error implements error

// Mark returns an error with the supplied errors as marks.
// If err is nil, return nil.
// marks take effects only when Is and '%v' in fmt.
// Is returns true if err or any marks match the target.
func Mark(err error, marks ...error) error {
	if err == nil {
		return nil
	}
	n := 0
	for _, mark := range marks {
		if mark != nil {
			// avoid repeat marks
			if errors.Is(err, mark) {
				return err
			}
			n++
		}
	}
	if n == 0 {
		return err
	}

	me := markError{
		err:   err,
		marks: make([]error, 0, n),
	}
	for _, mark := range marks {
		if mark != nil {
			me.marks = append(me.marks, mark)
		}
	}

	return me
}

type markError struct {
	err   error   // visual error
	marks []error // hidden errors as marks, take effects only when Is and '%v' in fmt.
}

func (e markError) Error() string {
	if e.err == nil {
		return ""
	}
	return e.err.Error()
}

func (e markError) Format(s fmt.State, verb rune) {
	if e.err == nil {
		return
	}
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "Marked errors occurred:%+v", e.err)
			for i, mark := range e.marks {
				_, _ = fmt.Fprintf(s, "\nM[%d]/[%d]\t%+v", i, len(e.marks), mark)
			}
			return
		}
		fallthrough
	case 's', 'q':
		_, _ = io.WriteString(s, e.Error())
	}
}

// Is reports whether any error in markError or it's mark errors matches target.
func (e markError) Is(target error) bool {
	if errors.Is(e.err, target) {
		return true
	}
	for _, err := range e.marks {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

// Unwrap returns the error in e, if there is exactly one. If there is more than one
// error, Unwrap returns nil, since there is no way to determine which should be
// returned.
func (e markError) Unwrap() error {
	return e.err
}

// As finds the first error in err's chain that matches target, and if one is found, sets
// target to that error value and returns true. Otherwise, it returns false.
func (e markError) As(target any) bool {
	if errors.As(e.err, target) {
		return true
	}
	for _, err := range e.marks {
		if errors.As(err, target) {
			return true
		}
	}
	return false
}
