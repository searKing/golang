// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cause

import (
	"fmt"
	"io"
)

// WithError annotates cause error with a new error.
// If cause is nil, WithError returns new error.
// If err is nil, WithError returns nil.
func WithError(cause error, err error) error {
	if cause == nil {
		return err
	}
	if err == nil {
		return nil
	}
	return &withError{
		cause: cause,
		err:   err,
	}
}

type withError struct {
	cause error
	err   error
}

func (w *withError) Error() string { return w.err.Error() + ": " + w.cause.Error() }
func (w *withError) Cause() error  { return w.cause }

func (w *withError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "%+v\n", w.Cause())
			_, _ = fmt.Fprintf(s, "%+v", w.err)
			return
		}
		fallthrough
	case 's', 'q':
		_, _ = io.WriteString(s, w.Error())
	}
}
