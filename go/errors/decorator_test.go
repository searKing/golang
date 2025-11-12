// Copyright 2025 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors_test

import (
	"errors"
	"fmt"
	"testing"

	errors_ "github.com/searKing/golang/go/errors"
)

func baseHandler(err error) error {
	return errors.New("base_start: " + err.Error() + " :base_end")

}

func prefixDecorator(next errors_.HandlerFunc) errors_.HandlerFunc {
	return func(err error) error {
		return errors.New(next(errors.New("prefix_start: "+err.Error())).Error() + " :prefix_end")
	}
}

func suffixDecorator(next errors_.HandlerFunc) errors_.HandlerFunc {
	return func(err error) error {
		return errors.New(next(errors.New("suffix_start: "+err.Error())).Error() + " :suffix_end")
	}
}

func interceptorDecorator(_ errors_.HandlerFunc) errors_.HandlerFunc {
	return func(err error) error {
		return errors.New("intercepted_start: " + err.Error() + " :intercepted_end")
	}
}

func TestHandlerDecoratorFunc(t *testing.T) {
	originalErr := errors.New("test error")
	decorator := errors_.HandlerDecoratorFunc(prefixDecorator)
	wrappedHandler := decorator.WrapHandler(baseHandler)
	resultErr := wrappedHandler(originalErr)
	expectedErr := "base_start: prefix_start: test error :base_end :prefix_end"
	if resultErr.Error() != expectedErr {
		t.Errorf("Expected error '%s', but got '%s'", expectedErr, resultErr.Error())
	}
}

func TestHandlerDecorators(t *testing.T) {
	originalErr := errors.New("test error")

	decorators := errors_.HandlerDecorators{
		errors_.HandlerDecoratorFunc(prefixDecorator),
		errors_.HandlerDecoratorFunc(suffixDecorator),
	}

	wrappedHandler := decorators.WrapHandler(baseHandler)
	resultErr := wrappedHandler(originalErr)
	expectedErr := "base_start: suffix_start: prefix_start: test error :base_end :suffix_end :prefix_end"
	if resultErr.Error() != expectedErr {
		t.Errorf("Expected error '%s', but got '%s'", expectedErr, resultErr.Error())
	}
}

func TestHandlerDecorators_Order(t *testing.T) {
	var executionOrder []string
	originalErr := errors.New("test error")
	var decorators errors_.HandlerDecorators
	for i := range 3 {
		decorators = append(decorators, errors_.HandlerDecoratorFunc(func(next errors_.HandlerFunc) errors_.HandlerFunc {
			return func(err error) error {
				executionOrder = append(executionOrder, fmt.Sprintf("m%d-before", i+1))
				if err != nil {
					err = fmt.Errorf("m%d: %w", i+1, err)
				}
				result := next(err)
				executionOrder = append(executionOrder, fmt.Sprintf("m%d-after", i+1))
				return result
			}
		}))
	}

	wrappedHandler := decorators.WrapHandler(func(err error) error {
		executionOrder = append(executionOrder, "handler")
		return err
	})
	resultErr := wrappedHandler(originalErr)
	// Verify execution order
	expected := []string{
		"m1-before", "m2-before", "m3-before", "handler",
		"m3-after", "m2-after", "m1-after",
	}

	if len(executionOrder) != len(expected) {
		t.Fatalf("execution order length mismatch. expected %d, got %d",
			len(expected), len(executionOrder))
	}

	for i, step := range expected {
		if executionOrder[i] != step {
			t.Errorf("execution order mismatch at step %d. expected %q, got %q",
				i, step, executionOrder[i])
		}
	}

	// Verify error transformation
	// Expected: m3: m2: m1: original (each middleware wraps the error before passing down)
	expectedErr := "m3: m2: m1: test error"
	if resultErr.Error() != expectedErr {
		t.Errorf("Expected error '%s', but got '%s'", expectedErr, resultErr.Error())
	}
}

func TestHandlerDecorators_Interceptor(t *testing.T) {
	originalErr := errors.New("test error")

	decorators := errors_.HandlerDecorators{
		errors_.HandlerDecoratorFunc(prefixDecorator),
		errors_.HandlerDecoratorFunc(interceptorDecorator),
		errors_.HandlerDecoratorFunc(suffixDecorator),
	}

	wrappedHandler := decorators.WrapHandler(baseHandler)
	resultErr := wrappedHandler(originalErr)

	expectedErr := "intercepted_start: prefix_start: test error :intercepted_end :prefix_end"
	if resultErr.Error() != expectedErr {
		t.Errorf("Expected error '%s', but got '%s'", expectedErr, resultErr.Error())
	}
}

// A simple decorator for testing purposes.

type testDecorator struct {
	name string
}

func (td *testDecorator) WrapHandler(h errors_.HandlerFunc) errors_.HandlerFunc {
	return func(err error) error {
		// Wrap the error before calling the next handler in the chain.
		return h(fmt.Errorf("%s -> %w", td.name, err))
	}
}

func TestHandlerDecorators_WrapError(t *testing.T) {
	// The original error to be wrapped.
	originalErr := errors.New("original error")

	// Create a couple of test decorators.
	decoratorA := &testDecorator{name: "decorator A"}
	decoratorB := &testDecorator{name: "decorator B"}

	// Put decorators into the HandlerDecorators slice.
	// The expected application order is: decoratorA -> decoratorB.
	// Since `WrapError` applies them in reverse, decoratorB will be the outermost wrapper.
	decorators := errors_.HandlerDecorators{decoratorA, decoratorB}

	// Call the WrapError method.
	wrappedErr := decorators.WrapError(originalErr)

	// Verify the final error string.
	expectedErrStr := "decorator B -> decorator A -> original error"
	if wrappedErr.Error() != expectedErrStr {
		t.Errorf("Wrapped error string is incorrect.\nExpected: %s\nGot:      %s", expectedErrStr, wrappedErr.Error())
	}

	// Verify that the original error is still in the error chain.
	if !errors.Is(wrappedErr, originalErr) {
		t.Errorf("Wrapped error does not contain the original error.")
	}
}

func TestHandlerDecorators_ErrorOr(t *testing.T) {
	originalErr := errors.New("test error")

	decorators := errors_.HandlerDecorators{}
	resultErr := decorators.ErrorOr(originalErr, errors.New("default error"))
	expectedErr := "default error"
	if resultErr.Error() != expectedErr {
		t.Errorf("Expected error '%s', but got '%s'", expectedErr, resultErr.Error())
	}
}
