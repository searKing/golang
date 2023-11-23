// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices

// SwigVector represents Go wrapper with std::vector<T>
//
// %include <std_vector.i>
// %template(vector_float) std::vector<float>;
// See: https://www.swig.org/Doc4.1/Go.html
type SwigVector[E any] interface {
	Size() int64
	Capacity() int64
	Reserve(int64)
	IsEmpty() bool
	Clear()
	Add(e E)
	Get(i int) E
	Set(i int, e E)
}

// FromSwigVector returns a slice mapped by v.Get() within all c in the SwigVector[E].
// FromSwigVector does not modify the contents of the SwigVector[E] v; it creates a new slice.
func FromSwigVector[S ~[]E, E any](v SwigVector[E]) S {
	if v == nil {
		return nil
	}

	var s = make(S, 0, v.Size())
	if v.IsEmpty() || v.Size() == 0 {
		return s
	}
	for i := 0; i < int(v.Size()); i++ {
		s = append(s, v.Get(i))
	}
	return s
}

// FromSwigVectorFunc returns a slice mapped by mapped by f(v.Get()) within all c in the SwigVector[E].
// FromSwigVector does not modify the contents of the SwigVector[E] v; it creates a new slice.
func FromSwigVectorFunc[S ~[]R, E any, R any](v SwigVector[E], f func(E) R) S {
	if v == nil {
		return nil
	}

	var s = make(S, 0, v.Size())
	if v.IsEmpty() || v.Size() == 0 {
		return s
	}
	for i := 0; i < int(v.Size()); i++ {
		s = append(s, f(v.Get(i)))
	}
	return s
}

// ToSwigVector returns a SwigVector[E] mapped by v.Get() within all c in the slice.
// ToSwigVector does not modify the contents of the slice s; it modifies the SwigVector[E] v.
func ToSwigVector[S ~[]E, E any](s S, v SwigVector[E]) {
	if len(s) == 0 {
		return
	}
	v.Clear()
	v.Reserve(int64(len(s)))
	for i := range s {
		v.Add(s[i])
	}
}

// ToSwigVectorFunc returns a SwigVector[E] mapped by f(c) within all c in the slice.
// ToSwigVectorFunc does not modify the contents of the slice s; it modifies the SwigVector[E] v.
func ToSwigVectorFunc[S ~[]E, E any, R any](s S, v SwigVector[R], f func(E) R) {
	if len(s) == 0 {
		return
	}
	v.Clear()
	v.Reserve(int64(len(s)))
	for i := range s {
		v.Add(f(s[i]))
	}
}
