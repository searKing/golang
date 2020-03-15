// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package util

import (
	"github.com/searKing/golang/go/util/object"
)

type Comparator interface {
	/**
	 * Compares its two arguments for order.  Returns a negative integer,
	 * zero, or a positive integer as the first argument is less than, equal
	 * to, or greater than the second.<p>
	 *
	 * The implementor must ensure that {@code sgn(compare(x, y)) ==
	 * -sgn(compare(y, x))} for all {@code x} and {@code y}.  (This
	 * implies that {@code compare(x, y)} must throw an exception if and only
	 * if {@code compare(y, x)} throws an exception.)<p>
	 *
	 * The implementor must also ensure that the relation is transitive:
	 * {@code ((compare(x, y)>0) && (compare(y, z)>0))} implies
	 * {@code compare(x, z)>0}.<p>
	 *
	 * Finally, the implementor must ensure that {@code compare(x, y)==0}
	 * implies that {@code sgn(compare(x, z))==sgn(compare(y, z))} for all
	 * {@code z}.<p>
	 *
	 * It is generally the case, but <i>not</i> strictly required that
	 * {@code (compare(x, y)==0) == (x.equals(y))}.  Generally speaking,
	 * any comparator that violates this condition should clearly indicate
	 * this fact.  The recommended language is "Note: this comparator
	 * imposes orderings that are inconsistent with equals."<p>
	 *
	 * In the foregoing description, the notation
	 * {@code sgn(}<i>expression</i>{@code )} designates the mathematical
	 * <i>signum</i> function, which is defined to return one of {@code -1},
	 * {@code 0}, or {@code 1} according to whether the value of
	 * <i>expression</i> is negative, zero, or positive, respectively.
	 *
	 * @param o1 the first object to be compared.
	 * @param o2 the second object to be compared.
	 * @return a negative integer, zero, or a positive integer as the
	 *         first argument is less than, equal to, or greater than the
	 *         second.
	 * @throws NullPointerException if an argument is null and this
	 *         comparator does not permit null arguments
	 * @throws ClassCastException if the arguments' types prevent them from
	 *         being compared by this comparator.
	 */
	Compare(a, b interface{}) int
	/**
	 * Returns a lexicographic-order comparator with another comparator.
	 * If this {@code Comparator} considers two elements equal, i.e.
	 * {@code compare(a, b) == 0}, {@code other} is used to determine the order.
	 *
	 * <p>The returned comparator is serializable if the specified comparator
	 * is also serializable.
	 *
	 * @apiNote
	 * For example, to sort a collection of {@code String} based on the length
	 * and then case-insensitive natural ordering, the comparator can be
	 * composed using following code,
	 *
	 * <pre>{@code
	 *     Comparator<String> cmp = Comparator.comparingInt(String::length)
	 *             .thenComparing(String.CASE_INSENSITIVE_ORDER);
	 * }</pre>
	 *
	 * @param  other the other comparator to be used when this comparator
	 *         compares two objects that are equal.
	 * @return a lexicographic-order comparator composed of this and then the
	 *         other comparator
	 * @throws NullPointerException if the argument is null.
	 * @since 1.8
	 */
	ThenComparing(after Comparator) Comparator

	Reversed() Comparator
}

// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as HTTP handlers. If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler that calls f.
type ComparatorFunc func(a, b interface{}) int

// Compare calls f(a,b).
func (f ComparatorFunc) Compare(a, b interface{}) int {
	return f(a, b)
}

func (f ComparatorFunc) Reversed() Comparator {
	return ComparatorFunc(func(a, b interface{}) int {
		return f.Compare(b, a)
	})
}

func (f ComparatorFunc) ThenComparing(after Comparator) Comparator {
	object.RequireNonNil(after)
	return ComparatorFunc(func(a, b interface{}) int {
		res := f.Compare(a, b)
		if res != 0 {
			return res
		}
		return after.Compare(a, b)
	})
}

/**
 * Returns a null-friendly comparator that considers {@code null} to be
 * less than non-null. When both are {@code null}, they are considered
 * equal. If both are non-null, the specified {@code Comparator} is used
 * to determine the order. If the specified comparator is {@code null},
 * then the returned comparator considers all non-null values to be equal.
 *
 * <p>The returned comparator is serializable if the specified comparator
 * is serializable.
 *
 * @param  <T> the type of the elements to be compared
 * @param  comparator a {@code Comparator} for comparing non-null values
 * @return a comparator that considers {@code null} to be less than
 *         non-null, and compares non-null objects with the supplied
 *         {@code Comparator}.
 * @since 1.8
 */
func NullFirst(cmp Comparator, nilFirst bool) Comparator {
	return NewNullComparator(nilFirst, cmp)
}

// DefaultComparator returns 0 if a == b, returns -1 otherwise
func DefaultComparator() Comparator {
	return ComparatorFunc(func(a, b interface{}) int {
		if a == b {
			return 0
		}
		return -1
	})
}

// AlwaysEqualComparator returns 0 always
func AlwaysEqualComparator() Comparator {
	return ComparatorFunc(func(_, _ interface{}) int {
		return 0
	})
}

// NeverEqualComparator returns -1 always
func NeverEqualComparator() Comparator {
	return ComparatorFunc(func(_, _ interface{}) int {
		return -1
	})
}
