// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slice

import (
	"reflect"
)

// Of returns a slice consisting of the elements.
// obj: Accept Array、Slice、String(as []byte if ifStringAsRune else []rune)
func Of(obj interface{}) []interface{} {
	return of(obj)
}

type MapPair struct {
	Key   interface{}
	Value interface{}
}

//of is the same as Of
func of(obj interface{}) []interface{} {
	switch kind := reflect.ValueOf(obj).Kind(); kind {
	default:
		// element as a slice of one element
		out := []interface{}{}
		out = append(out, obj)
	case reflect.Array, reflect.Slice:
	case reflect.Map:
		out := []interface{}{}
		v := reflect.ValueOf(obj)
		keys := v.MapKeys()
		for _, k := range keys {
			e := v.MapIndex(k)
			pair := MapPair{
				Key:   reflect.Indirect(k).Interface(),
				Value: reflect.Indirect(e).Interface(),
			}
			out = append(out, pair)
		}
		return out
	}

	out := []interface{}{}
	v := reflect.ValueOf(obj)
	for i := 0; i < v.Len(); i++ {
		out = append(out, v.Slice(i, i+1).Index(0).Interface())
	}
	return out
}
