// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reflect

import (
	"reflect"
)

// Walk walks down v
func Walk(t reflect.Type, visitedOnce bool, do func(s reflect.Type, sf reflect.StructField) (stop bool)) {
	// Anonymous fields to explore at the current level and the next.
	var current []reflect.Type
	next := []reflect.Type{t}

	// Count of queued names for current level and the next.
	currentCount := map[reflect.Type]int{}
	nextCount := map[reflect.Type]int{}

	// Types already visited at an earlier level.
	// FIXME I havenot seen any case which can trigger visited
	visited := map[reflect.Type]bool{}
	for len(next) > 0 {
		current, next = next, current[:0]
		currentCount, nextCount = nextCount, map[reflect.Type]int{}

		for _, typ := range current {

			if typ.Kind() == reflect.Ptr {
				// Follow pointer.
				typ = typ.Elem()
			}
			if visitedOnce {
				if visited[typ] {
					continue
				}
				visited[typ] = true
			}

			if typ.Kind() != reflect.Struct {
				if do(typ, reflect.StructField{}) {
					return
				}
				continue
			}
			// Scan typ for fields to include.
			for i := 0; i < typ.NumField(); i++ {
				sf := typ.Field(i)
				if do(typ, sf) {
					continue
				}

				ft := sf.Type
				if ft.Name() == "" && ft.Kind() == reflect.Ptr {
					// Follow pointer.
					ft = ft.Elem()
				}

				// Record found field and index sequence.
				if ft.Name() != "" || !sf.Anonymous || ft.Kind() != reflect.Struct {
					if currentCount[typ] > 1 {
					}
					//continue
				}
				// Record new anonymous struct to explore in next round.
				nextCount[ft]++
				if !visitedOnce || nextCount[ft] == 1 {
					next = append(next, ft)
				}
			}
		}
	}

}
