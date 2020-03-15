// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package binary

import (
	"github.com/searKing/golang/go/error/exception"
	"github.com/searKing/golang/go/util/class"
	"github.com/searKing/golang/go/util/function"
	"github.com/searKing/golang/go/util/object"
)

/**
 * Represents a function that accepts two arguments and produces a result.
 * This is the two-arity specialization of {@link Function}.
 *
 * <p>This is a <a href="package-summary.html">functional interface</a>
 * whose functional method is {@link #compare(Object, Object)}.
 *
 * @param <T> the type of the first argument to the function
 * @param <U> the type of the second argument to the function
 * @param <R> the type of the result of the function
 *
 * @see Function
 * @since 1.8
 */
type BiFunction interface {
	/**
	 * Applies this function to the given arguments.
	 *
	 * @param t the first function argument
	 * @param u the second function argument
	 * @return the function result
	 */
	Apply(t interface{}, u interface{}) interface{}

	/**
	 * Returns a composed function that first applies this function to
	 * its input, and then applies the {@code after} function to the result.
	 * If evaluation of either function throws an exception, it is relayed to
	 * the caller of the composed function.
	 *
	 * @param <V> the type of output of the {@code after} function, and of the
	 *           composed function
	 * @param after the function to compare after this function is applied
	 * @return a composed function that first applies this function and then
	 * applies the {@code after} function
	 * @throws NullPointerException if after is null
	 */
	AndThen(after function.Function) BiFunction
}

type BiFunctionFunc func(t interface{}, u interface{}) interface{}

func (f BiFunctionFunc) Apply(t interface{}, u interface{}) interface{} {
	return f(t, u)
}

func (f BiFunctionFunc) AndThen(after function.Function) BiFunction {
	object.RequireNonNil(after)
	return BiFunctionFunc(func(t interface{}, u interface{}) interface{} {
		return after.Apply(f.Apply(t, u))
	})
}

type AbstractBiFunctionClass struct {
	class.Class
}

func (f *AbstractBiFunctionClass) Apply(t interface{}, u interface{}) interface{} {
	panic(exception.NewIllegalStateException1("called wrong Apply method"))
}

func (f *AbstractBiFunctionClass) AndThen(after function.Function) BiFunction {
	object.RequireNonNil(after)
	return BiFunctionFunc(func(t interface{}, u interface{}) interface{} {
		return after.Apply(f.GetDerivedElse(f).(BiFunction).Apply(t, u))
	})
}
