package cause

import (
	"fmt"
	"io"
)

// WithError annotates cause error with a new error.
// If err is nil, WithError returns new error.
func WithError(cause error, err error) error {
	if cause == nil {
		return err
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
