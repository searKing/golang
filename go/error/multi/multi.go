package multi

import (
	"fmt"
	"io"
)

// New returns an error with the supplied errors.
func New(errs ...error) error {
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
