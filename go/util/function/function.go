// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package function

import (
	"github.com/searKing/golang/go/error/exception"
	"github.com/searKing/golang/go/util/class"
	"github.com/searKing/golang/go/util/object"
)

/**
 * Represents a function that accepts one argument and produces a result.
 *
 * <p>This is a <a href="package-summary.html">functional interface</a>
 * whose functional method is {@link #apply(Object)}.
 *
 * @param <T> the type of the input to the function
 * @param <R> the type of the result of the function
 *
 * @since 1.8
 */
type Function interface {
	/**
	 * Applies this function to the given argument.
	 *
	 * @param t the function argument
	 * @return the function result
	 */
	Apply(t interface{}) interface{}

	/**
	 * Returns a composed function that first applies the {@code before}
	 * function to its input, and then applies this function to the result.
	 * If evaluation of either function throws an exception, it is relayed to
	 * the caller of the composed function.
	 *
	 * @param <V> the type of input to the {@code before} function, and to the
	 *           composed function
	 * @param before the function to apply before this function is applied
	 * @return a composed function that first applies the {@code before}
	 * function and then applies this function
	 * @throws NullPointerException if before is null
	 *
	 * @see #andThen(Function)
	 */
	Compose(before Function) Function

	/**
	 * Returns a composed function that first applies this function to
	 * its input, and then applies the {@code after} function to the result.
	 * If evaluation of either function throws an exception, it is relayed to
	 * the caller of the composed function.
	 *
	 * @param <V> the type of output of the {@code after} function, and of the
	 *           composed function
	 * @param after the function to apply after this function is applied
	 * @return a composed function that first applies this function and then
	 * applies the {@code after} function
	 * @throws NullPointerException if after is null
	 *
	 * @see #compose(Function)
	 */
	AndThen(before Function) Function
}

/**
 * Returns a function that always returns its input argument.
 *
 * @param <T> the type of the input and output objects to the function
 * @return a function that always returns its input argument
 */
func Identity() Function {
	return FunctionFunc(func(t interface{}) interface{} {
		return t
	})
}

type FunctionFunc func(t interface{}) interface{}

// Apply calls f(t).
func (f FunctionFunc) Apply(t interface{}) interface{} {
	return f(t)
}

func (f FunctionFunc) Compose(before Function) Function {
	object.RequireNonNil(before)
	return FunctionFunc(func(t interface{}) interface{} {
		return f.Apply(before.Apply(t))
	})
}

func (f FunctionFunc) AndThen(after Function) Function {
	object.RequireNonNil(after)
	return FunctionFunc(func(t interface{}) interface{} {
		return after.Apply(f.Apply(t))
	})
}

type AbstractFunction struct {
	class.Class
}

func (f *AbstractFunction) Apply(t interface{}) interface{} {
	panic(exception.NewIllegalStateException1("called wrong Apply method"))
}

func (f *AbstractFunction) Compose(before Function) Function {
	object.RequireNonNil(before)
	return FunctionFunc(func(t interface{}) interface{} {
		return f.GetDerivedElse(f).(Function).Apply(before.Apply(t))
	})
}

func (f *AbstractFunction) AndThen(after Function) Function {
	object.RequireNonNil(after)
	return FunctionFunc(func(t interface{}) interface{} {
		return after.Apply(f.GetDerivedElse(f).(Function).Apply(t))
	})
}
