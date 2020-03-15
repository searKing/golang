// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package binary

import (
	"github.com/searKing/golang/go/util"
	"github.com/searKing/golang/go/util/object"
)

/**
 * Returns a {@link BinaryOperator} which returns the lesser of two elements
 * according to the specified {@code Comparator}.
 *
 * @param <T> the type of the input arguments of the comparator
 * @param comparator a {@code Comparator} for comparing the two values
 * @return a {@code BinaryOperator} which returns the lesser of its operands,
 *         according to the supplied {@code Comparator}
 * @throws NullPointerException if the argument is null
 */
func MinBy(comparator util.Comparator) BiFunction {
	object.RequireNonNil(comparator)
	return BiFunctionFunc(func(t interface{}, u interface{}) interface{} {
		c := comparator.Compare(t, u)
		if c <= 0 {
			return t
		}

		return u
	})
}

/**
 * Returns a {@link BinaryOperator} which returns the greater of two elements
 * according to the specified {@code Comparator}.
 *
 * @param <T> the type of the input arguments of the comparator
 * @param comparator a {@code Comparator} for comparing the two values
 * @return a {@code BinaryOperator} which returns the greater of its operands,
 *         according to the supplied {@code Comparator}
 * @throws NullPointerException if the argument is null
 */
func MaxBy(comparator util.Comparator) BiFunction {
	object.RequireNonNil(comparator)
	return BiFunctionFunc(func(t interface{}, u interface{}) interface{} {
		c := comparator.Compare(t, u)
		if c >= 0 {
			return t
		}

		return u
	})
}
