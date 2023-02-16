// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package maps

import "reflect"

// Filter returns a map satisfying c != zero within all c in the map.
// Filter modifies the contents of the map s; it does not create a new map.
func Filter[M ~map[K]V, K comparable, V comparable](m M) M {
	if len(m) == 0 {
		return m
	}
	for k, v := range m {
		var zeroV V
		if v == zeroV {
			delete(m, k)
		}
	}
	return m
}

// FilterFunc returns a map satisfying f(c) within all c in the map.
// FilterFunc modifies the contents of the map s; it does not create a new map.
func FilterFunc[M ~map[K]V, K comparable, V any](m M, f func(K, V) bool) M {
	if len(m) == 0 {
		return m
	}
	for k, v := range m {
		if !f(k, v) {
			delete(m, k)
		}
	}
	return m
}

// TypeAssertFilter returns a map satisfying r, ok := any(c).(R); ok == true within all r in the map.
// TypeAssertFilter does not modify the contents of the map m; it creates a new map.
func TypeAssertFilter[M ~map[K]V, M2 ~map[K2]V2, K comparable, V comparable, K2 comparable, V2 comparable](m M) M2 {
	if len(m) == 0 {
		if m == nil {
			return nil
		}
		var emptyM2 = M2{}
		return emptyM2
	}

	var m2 = M2{}

	var zeroK K
	var zeroV V
	var zeroK2 K2
	var zeroV2 V2
	var nilableK = any(zeroK) == nil
	var nilableV = any(zeroV) == nil

	var k2t = reflect.TypeOf(zeroK2)
	var v2t = reflect.TypeOf(zeroV2)
	var convertibleK bool
	var convertibleV bool
	if !nilableK || any(zeroK) != nil {
		var et = reflect.TypeOf(zeroK)
		convertibleK = et.ConvertibleTo(k2t)
	}
	if !nilableV || any(zeroV) != nil {
		var et = reflect.TypeOf(zeroV)
		convertibleV = et.ConvertibleTo(v2t)
	}

	for k, v := range m {
		if !convertibleK && !nilableK {
			continue
		}
		if !convertibleV && !nilableV {
			continue
		}

		if (any(k) == nil && any(zeroK2) == nil) && (any(v) == nil && any(zeroV2) == nil) {
			var zeroK2 K2
			var zeroV2 V2
			m2[zeroK2] = zeroV2
			continue
		}

		if k2, ok := any(k).(K2); ok {
			if v2, ok := any(v).(V2); ok {
				m2[k2] = v2
				continue
			}
		}
		if convertibleK && convertibleV {
			if k2, ok := reflect.ValueOf(k).Convert(k2t).Interface().(K2); ok {
				if v2, ok := reflect.ValueOf(v).Convert(v2t).Interface().(V2); ok {
					m2[k2] = v2
					continue
				}
			}
		}
	}
	return m2
}

// TypeAssertFilterFunc returns a map satisfying f(c) within all c in the map.
// TypeAssertFilterFunc does not modify the contents of the map m; it creates a new map.
func TypeAssertFilterFunc[M ~map[K]V, M2 ~map[K2]V2, K comparable, V any, K2 comparable, V2 any](m M, f func(K, V) (K2, V2, bool)) M2 {
	if len(m) == 0 {
		if m == nil {
			return nil
		}
		var emptyM2 = M2{}
		return emptyM2
	}

	var m2 = M2{}
	for k, v := range m {
		if k, v, ok := f(k, v); ok {
			m2[k] = v
		}
	}
	return m2
}
