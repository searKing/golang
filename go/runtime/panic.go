package runtime

import (
	"fmt"
	"log"
	"net/http"
)

var (
	DefaultPanic     = Panic{}
	NeverPanic       = Panic{IgnoreCrash: true}
	LogPanic         = Panic{PanicHandlers: []func(any){logPanic}}
	NeverPanicButLog = Panic{IgnoreCrash: true, PanicHandlers: []func(any){logPanic}}
)

// Panic simply catches a panic and logs an error. Meant to be called via
// defer.  Additional context-specific handlers can be provided, and will be
// called in case of panic.
type Panic struct {
	// IgnoreCrash controls the behavior of Recover and now defaults false.
	// if false, crash immediately, rather than eating panics.
	IgnoreCrash bool

	// PanicHandlers for something like logging the panic message, shutting down go routines gracefully.
	PanicHandlers []func(any)
}

// Recover actually crashes if IgnoreCrash is false, after calling PanicHandlers.
func (p Panic) Recover() {
	if r := recover(); r != nil {
		for _, fn := range p.PanicHandlers {
			fn(r)
		}
		if p.IgnoreCrash {
			return
		}
		// Actually proceed to panic.
		panic(r)
	}
}

func (p *Panic) AppendHandler(handlers ...func(any)) *Panic {
	p.PanicHandlers = append(p.PanicHandlers, handlers...)
	return p
}

func HandlePanicWith(handlers ...func(any)) Panic {
	p := Panic{}
	p.AppendHandler(handlers...)
	return p
}

// RecoverFromPanic replaces the specified error with an error containing the
// original error, and the call tree when a panic occurs. This enables error
// handlers to handle errors and panics the same way.
func RecoverFromPanic(err error) error {
	if r := recover(); r != nil {
		const size = 64 << 10
		stacktrace := GetCallStack(size)
		if err == nil {
			return fmt.Errorf(
				"recovered from panic %q. Call stack:\n%s",
				r,
				stacktrace)
		}

		return fmt.Errorf(
			"recovered from panic %q. (err=%w) Call stack:\n%s",
			r,
			err,
			stacktrace)
	}
	return err
}

// logPanic logs the caller tree when a panic occurs (except in the special case of http.ErrAbortHandler).
func logPanic(r any) {
	if r == http.ErrAbortHandler {
		// honor the http.ErrAbortHandler sentinel panic value:
		//   ErrAbortHandler is a sentinel panic value to abort a handler.
		//   While any panic from ServeHTTP aborts the response to the client,
		//   panicking with ErrAbortHandler also suppresses logging of a stack trace to the server's error log.
		return
	}

	const size = 64 << 10
	stacktrace := GetCallStack(size)
	if _, ok := r.(string); ok {
		log.Printf("Observed a panic: %s\n%s", r, stacktrace)
	} else {
		log.Printf("Observed a panic: %#v (%v)\n%s", r, r, stacktrace)
	}
}
