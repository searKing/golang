// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices

import "reflect"

// Filter returns a slice satisfying c != zero within all c in the slice.
// Filter modifies the contents of the slice s; it does not create a new slice.
func Filter[S ~[]E, E comparable](s S) S {
	if len(s) == 0 {
		return s
	}
	i := 0
	for _, v := range s {
		var zeroE E
		if v != zeroE {
			s[i] = v
			i++
		}
	}
	return s[:i]
}

// FilterFunc returns a slice satisfying f(c) within all c in the slice.
// FilterFunc modifies the contents of the slice s; it does not create a new slice.
func FilterFunc[S ~[]E, E any](s S, f func(E) bool) S {
	if len(s) == 0 {
		return s
	}
	i := 0
	for _, v := range s {
		if f(v) {
			s[i] = v
			i++
		}
	}
	return s[:i]
}

// TypeAssertFilter returns a slice satisfying r, ok := any(c).(R); ok == true within all r in the slice.
// TypeAssertFilter does not modify the contents of the slice s; it creates a new slice.
func TypeAssertFilter[S ~[]E, M ~[]R, E any, R any](s S) M {
	if len(s) == 0 {
		if s == nil {
			return nil
		}
		var emptyM = M{}
		return emptyM
	}
	var m = M{}

	var zeroE E
	var zeroR R
	var nilable = any(zeroE) == nil

	var rt = reflect.TypeOf(zeroR)
	var convertible bool
	if !nilable || any(zeroE) != nil {
		var et = reflect.TypeOf(zeroE)
		convertible = et.ConvertibleTo(rt)
	}

	for _, v := range s {
		if !convertible && !nilable {
			continue
		}

		if any(v) == nil && any(zeroR) == nil {
			var zeroR R
			m = append(m, zeroR)
			continue
		}

		if r, ok := any(v).(R); ok {
			m = append(m, r)
			continue
		}
		if convertible {
			if r, ok := reflect.ValueOf(v).Convert(rt).Interface().(R); ok {
				m = append(m, r)
				continue
			}
		}
	}
	return m
}

// TypeAssertFilterFunc returns a slice satisfying f(c) within all c in the slice.
// TypeAssertFilterFunc does not modify the contents of the slice s; it creates a new slice.
func TypeAssertFilterFunc[S ~[]E, M ~[]R, E any, R any](s S, f func(E) (R, bool)) M {
	if len(s) == 0 {
		if s == nil {
			return nil
		}
		var emptyM = M{}
		return emptyM
	}

	var m = M{}
	for _, v := range s {
		if r, ok := f(v); ok {
			m = append(m, r)
		}
	}
	return m
}
