// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package multi

import (
	"fmt"
	"io"
)

// New returns an error with the supplied errors.
// If no error contained, return nil.
// Deprecated: Use errors.Multi instead.
func New(errs ...error) error {
	var has bool
	for _, err := range errs {
		if err != nil {
			has = true
			break
		}
	}
	if !has {
		return nil
	}
	return &multiError{
		errs: errs,
	}
}

type multiError struct {
	errs []error
}

func (w *multiError) Error() string {
	if w == nil || len(w.errs) == 0 {
		return ""
	}
	message := w.errs[0].Error()
	for _, err := range w.errs[1:] {
		message += "|" + err.Error()
	}

	return message
}

func (w *multiError) Format(s fmt.State, verb rune) {
	if w == nil {
		return
	}
	switch verb {
	case 'v':
		if s.Flag('+') {
			if len(w.errs) == 0 {
				return
			}

			_, _ = io.WriteString(s, "Multiple errors occurred:\n\t")

			_, _ = io.WriteString(s, w.errs[0].Error())

			for _, err := range w.errs[1:] {
				_, _ = fmt.Fprintf(s, "|%+v", err)
			}
			return
		}
		fallthrough
	case 's', 'q':
		_, _ = io.WriteString(s, w.Error())
	}
}
