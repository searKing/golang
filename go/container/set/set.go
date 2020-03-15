// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package set

import (
	"github.com/searKing/golang/go/container/slice"
	"github.com/searKing/golang/go/util/object"
)

// take advantage of map
type Set struct {
	elems map[interface{}]struct{}
}

// Create a new set
func New() *Set {
	return &Set{}
}

func (s *Set) Init() *Set {
	s.elems = make(map[interface{}]struct{})
	return s
}

// lazyInit lazily initializes a zero List value.
func (s *Set) lazyInit() *Set {
	if s.elems == nil {
		s.Init()
	}
	return s
}

// Count returns the number of elements in this set (its cardinality).
func (s *Set) Count() int {
	if s.elems == nil {
		return 0
	}
	return len(s.elems)
}

// IsEmpty returns {@code true} if this set contains no elements.
func (s *Set) IsEmpty() bool {
	if s.Count() == 0 {
		return true
	}
	return false
}

// Contains returns {@code true} if this set contains the specified element.
func (s *Set) Contains(e interface{}) bool {
	if s.IsEmpty() {
		return false
	}
	if _, ok := s.elems[e]; ok {
		return true
	}
	return false
}

// ToSLice returns a slice containing all of the elements in this set.
func (s *Set) ToSlice() []interface{} {
	es := make([]interface{}, 0, s.Count())
	if s.IsEmpty() {
		return es
	}
	for e, _ := range s.elems {
		es = append(es, e)
	}
	return es
}

// Add adds the specified element to this set if it is not already present.
func (s *Set) Add(e interface{}) *Set {
	s.lazyInit()
	s.elems[e] = struct{}{}
	return s
}

// Remove removes the specified element from this set if it is present.
func (s *Set) Remove(e interface{}) *Set {
	if !s.Contains(e) {
		return s
	}
	delete(s.elems, e)
	return s
}

// ContainsAll returns {@code true} {@code true} if this set contains all of the elements of the specified collection.
func (s *Set) ContainsAll(stream *slice.Stream) bool {
	object.RequireNonNil(stream)
	if s.IsEmpty() {
		return false
	}
	for k, _ := range s.elems {
		if stream.NoneMatch(func(e interface{}) bool {
			if e == k {
				return true
			}
			return false
		}) {
			return false
		}
	}
	return true
}

// AddAll adds all of the elements in the specified collection to this set if
// they're not already present (optional operation).
func (s *Set) AddAll(stream *slice.Stream) *Set {
	object.RequireNonNil(stream)
	stream.ForEach(func(e interface{}) {
		s.Add(e)
	})
	return s
}

// AddAll adds all of the elements in the specified collection to this set if
// they're not already present (optional operation).
// <p>This operation processes the elements one at a time, in encounter
// order if one exists.  Performing the action for one element
// performing the action for subsequent elements, but for any given element,
// the action may be performed in whatever thread the library chooses.
func (s *Set) AddAllOrdered(stream *slice.Stream) *Set {
	object.RequireNonNil(stream)
	stream.ForEachOrdered(func(e interface{}) {
		s.Add(e)
	})
	return s
}

// RetainAll retains only the elements in this set that are contained in the
// specified collection (optional operation).  In other words, removes
// from this set all of its elements that are not contained in the
// specified collection.
func (s *Set) RetainAll(stream *slice.Stream) *Set {
	object.RequireNonNil(stream)
	if s.IsEmpty() {
		return s
	}
	// stream is too big
	retained := stream.Filter(func(e interface{}) bool {
		_, ok := s.elems[e]
		return ok
	})
	s.Clear()
	s.AddAll(retained)
	return s
}

// RemoveAll removes from this set all of its elements that are contained in the
// specified collection (optional operation).  If the specified
// collection is also a set, this operation effectively modifies this
// set so that its value is the <i>asymmetric set difference</i> of
// the two sets.
func (s *Set) RemoveAll(stream *slice.Stream) *Set {
	object.RequireNonNil(stream)
	if s.IsEmpty() {
		return s
	}
	stream.ForEachOrdered(func(e interface{}) {
		if s.Contains(e) {
			s.Remove(e)
		}
	})
	return s
}

// Clear removes all of the elements from this set (optional operation).
// The set will be empty after this call returns.
func (s *Set) Clear() *Set {
	s.elems = nil
	return s
}

// Compares the specified object with this set for equality.  Returns
// {@code true} if the specified object is also a set, the two sets
// have the same size, and every member of the specified set is
// contained in this set (or equivalently, every member of this set is
// contained in the specified set).
func (s *Set) Equals(other *Set) bool {
	object.RequireNonNil(other)
	if s == other {
		return true
	}

	if s.IsEmpty() != other.IsEmpty() {
		return false
	}

	if s.Count() != other.Count() {
		return false
	}
	for k, _ := range s.elems {
		if !other.Contains(k) {
			return false
		}
	}
	return true
}

// Clone returns a deepcopy of it's self
func (s *Set) Clone() *Set {
	return (object.DeepClone(s)).(*Set)
}

func (s *Set) ToStream() *slice.Stream {
	return slice.NewStream().WithSlice(s.ToSlice())
}

// Of returns an unmodifiable set containing the input element(s).
func Of(es ...interface{}) *Set {
	s := New()
	slice.NewStream().WithSlice(es).ForEachOrdered(
		func(e interface{}) {
			s.Add(e)
		})
	return s
}

//grammar surger for count
func (s *Set) Size() int {
	return s.Count()
}

func (s *Set) Length() int {
	return s.Count()
}
