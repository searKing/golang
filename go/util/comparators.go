// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package util

import "github.com/searKing/golang/go/util/object"

var NaturalOrderComparator = ComparableComparatorFunc(nil)
var ReverseOrderComparator = NaturalOrderComparator.Reversed()

type ComparableComparatorFunc func(a, b interface{}) int

func (f ComparableComparatorFunc) Compare(a, b interface{}) int {
	if ac, ok := a.(Comparable); ok {
		return ac.CompareTo(b)
	}
	if bc, ok := b.(Comparable); ok {
		return -1 * bc.CompareTo(a)
	}
	object.RequireNonNil(f)
	return f(a, b)
}

func (f ComparableComparatorFunc) Reversed() Comparator {
	return ComparatorFunc(func(a, b interface{}) int {
		return f.Compare(b, a)
	})
}

func (f ComparableComparatorFunc) ThenComparing(after Comparator) Comparator {
	object.RequireNonNil(after)
	return ComparatorFunc(func(a, b interface{}) int {
		res := f.Compare(a, b)
		if res != 0 {
			return res
		}
		return after.Compare(a, b)
	})
}

// Null-friendly comparators
type NullComparator struct {
	nilFirst bool
	// if null, non-null Ts are considered equal
	real Comparator
}

func NewNullComparator(nilFirst bool, real Comparator) *NullComparator {
	return &NullComparator{
		nilFirst: nilFirst,
		real:     real,
	}
}

func (n *NullComparator) Compare(a, b interface{}) int {
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		if n.nilFirst {
			return -1
		}
		return 1
	}

	if b == nil {
		if n.nilFirst {
			return 1
		}
		return -1
	}

	if n.real == nil {
		return 0
	}

	return n.real.Compare(a, b)
}

func (n *NullComparator) ThenComparing(other Comparator) Comparator {
	object.RequireNonNil(other)

	if n.real == nil {
		return other
	}
	return n.real.ThenComparing(other)
}

func (n *NullComparator) Reversed() Comparator {
	if n.real == nil {
		return &NullComparator{
			nilFirst: !n.nilFirst,
			real:     nil,
		}
	}
	return &NullComparator{
		nilFirst: !n.nilFirst,
		real:     n.real.Reversed(),
	}
}
