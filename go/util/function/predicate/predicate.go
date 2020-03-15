// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package predicate

import (
	"github.com/searKing/golang/go/error/exception"
	"github.com/searKing/golang/go/util/class"
	"github.com/searKing/golang/go/util/object"
)

/**
 * Represents a predicate (boolean-valued function) of one argument.
 */
type Predicater interface {
	/**
	 * Evaluates this predicate on the given argument.
	 *
	 * @param t the input argument
	 * @return {@code true} if the input argument matches the predicate,
	 * otherwise {@code false}
	 */
	Test(value interface{}) bool

	/**
	 * Returns a composed predicate that represents a short-circuiting logical
	 * AND of this predicate and another.  When evaluating the composed
	 * predicate, if this predicate is {@code false}, then the {@code other}
	 * predicate is not evaluated.
	 *
	 * <p>Any exceptions thrown during evaluation of either predicate are relayed
	 * to the caller; if evaluation of this predicate throws an exception, the
	 * {@code other} predicate will not be evaluated.
	 *
	 * @param other a predicate that will be logically-ANDed with this
	 *              predicate
	 * @return a composed predicate that represents the short-circuiting logical
	 * AND of this predicate and the {@code other} predicate
	 * @throws NullPointerException if other is null
	 */
	And(other Predicater) Predicater

	/**
	 * Returns a predicate that represents the logical negation of this
	 * predicate.
	 *
	 * @return a predicate that represents the logical negation of this
	 * predicate
	 */
	Negate() Predicater

	/**
	 * Returns a composed predicate that represents a short-circuiting logical
	 * OR of this predicate and another.  When evaluating the composed
	 * predicate, if this predicate is {@code true}, then the {@code other}
	 * predicate is not evaluated.
	 *
	 * <p>Any exceptions thrown during evaluation of either predicate are relayed
	 * to the caller; if evaluation of this predicate throws an exception, the
	 * {@code other} predicate will not be evaluated.
	 *
	 * @param other a predicate that will be logically-ORed with this
	 *              predicate
	 * @return a composed predicate that represents the short-circuiting logical
	 * OR of this predicate and the {@code other} predicate
	 * @throws NullPointerException if other is null
	 */
	Or(other Predicater) Predicater

	/**
	 * Returns a predicate that tests if two arguments are equal according
	 * to {@link Objects#equals(Object, Object)}.
	 *
	 * @param <T> the type of arguments to the predicate
	 * @param targetRef the object reference with which to compare for equality,
	 *               which may be {@code null}
	 * @return a predicate that tests if two arguments are equal according
	 * to {@link Objects#equals(Object, Object)}
	 */
	IsEqual(targetRef interface{}) Predicater
}

type PredicaterFunc func(value interface{}) bool

func (f PredicaterFunc) Test(value interface{}) bool {
	return f(value)
}

func (f PredicaterFunc) And(other Predicater) Predicater {
	object.RequireNonNil(other)
	return PredicaterFunc(func(value interface{}) bool {
		return f.Test(value) && other.Test(value)
	})
}

func (f PredicaterFunc) Negate() Predicater {
	return PredicaterFunc(func(value interface{}) bool {
		return !f.Test(value)
	})
}

func (f PredicaterFunc) Or(other Predicater) Predicater {
	object.RequireNonNil(other)
	return PredicaterFunc(func(value interface{}) bool {
		return f.Test(value) || other.Test(value)
	})
}

func (f PredicaterFunc) IsEqual(targetRef interface{}) Predicater {
	if targetRef == nil {
		return PredicaterFunc(object.IsNil)
	}
	return PredicaterFunc(func(value interface{}) bool {
		return value == targetRef
	})
}

type AbstractPredicaterClass struct {
	class.Class
}

func (pred *AbstractPredicaterClass) Test(value interface{}) bool {
	panic(exception.NewIllegalStateException1("called wrong Test method"))
}

func (pred *AbstractPredicaterClass) And(other Predicater) Predicater {
	return PredicaterFunc(func(value interface{}) bool {
		return pred.GetDerivedElse(pred).(Predicater).Test(value) && other.Test(value)
	})
}

func (pred *AbstractPredicaterClass) Negate() Predicater {
	return PredicaterFunc(func(value interface{}) bool {
		return !pred.GetDerivedElse(pred).(Predicater).Test(value)
	})
}

func (pred *AbstractPredicaterClass) Or(other Predicater) Predicater {
	object.RequireNonNil(other)
	return PredicaterFunc(func(value interface{}) bool {
		return pred.GetDerivedElse(pred).(Predicater).Test(value) || other.Test(value)
	})
}

func (pred *AbstractPredicaterClass) IsEqual(targetRef interface{}) Predicater {
	if targetRef == nil {
		return PredicaterFunc(object.IsNil)
	}
	return PredicaterFunc(func(value interface{}) bool {
		return value == targetRef
	})
}
