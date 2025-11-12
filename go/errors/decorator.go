// Copyright 2025 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	"slices"
)

// HandlerFunc defines an error processing function that takes an error
// and returns a processed error (or nil).
type HandlerFunc func(err error) error

// HandlerDecorator is an interface for decorating errors.
// It provides a WrapHandler method that takes a HandlerFunc and returns a new HandlerFunc.
// This allows the decorator to wrap the error processing function.
type HandlerDecorator interface {
	// WrapHandler returns a new error that wraps the provided error.
	WrapHandler(h HandlerFunc) HandlerFunc
}

// HandlerDecoratorFunc is a function type that implements the HandlerDecorator interface.
// This allows a simple function to be used as a HandlerDecorator.
type HandlerDecoratorFunc func(HandlerFunc) HandlerFunc

// WrapHandler implements the HandlerDecorator interface by calling the underlying function.
func (f HandlerDecoratorFunc) WrapHandler(h HandlerFunc) HandlerFunc { return f(h) }

// HandlerDecorators is a composite decorator that applies a slice of decorators
// to an error. The decorators are applied in reverse order to simulate
// a typical wrapping chain (e.g., the last added decorator is the outermost).
type HandlerDecorators []HandlerDecorator

// WrapHandler applies all decorators in the slice, in reverse order, to the error.
func (hds HandlerDecorators) WrapHandler(next HandlerFunc) HandlerFunc {
	for _, decorator := range slices.Backward(hds) {
		next = decorator.WrapHandler(next)
	}
	return next
}

// WrapError applies all decorators in the slice, in reverse order, to the error.
func (hds HandlerDecorators) WrapError(err error) error {
	return hds.WrapHandler(func(err error) error { return err })(err)
}

// ErrorOr applies all decorators in the slice, in reverse order, to the error.
// If the error is not handled, returns the default error.
func (hds HandlerDecorators) ErrorOr(err error, defaultErr error) error {
	return hds.WrapHandler(func(err error) error {
		if err == nil {
			return nil
		}
		return defaultErr
	})(err)
}
