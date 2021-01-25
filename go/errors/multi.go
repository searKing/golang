// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	"errors"
	"fmt"
	"io"
)

// New returns an error with the supplied errors.
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

// clean removes all none nil elem in all of the errors
func (e multiError) clean() multiError {
	var errs []error
	for _, err := range e {
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

// Is reports whether any error in multiError and it's chain chain matches target.
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
