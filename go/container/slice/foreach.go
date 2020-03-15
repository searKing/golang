// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slice

import (
	"sync"

	"github.com/searKing/golang/go/util/object"
)

// ForEachFunc Performs an action for each element of this slice.
// <p>The behavior of this operation is explicitly nondeterministic.
// For parallel slice pipelines, this operation does <em>not</em>
// guarantee to respect the encounter order of the slice, as doing so
// would sacrifice the benefit of parallelism.  For any given element, the
// action may be performed at whatever time and in whatever thread the
// library chooses.  If the action accesses shared state, it is
// responsible for providing the required synchronization.
func ForEachFunc(s interface{}, f func(interface{})) {
	forEachFunc(Of(s), f)
}

// forEachFunc is the same as ForEachFunc
func forEachFunc(s []interface{}, f func(interface{})) {
	object.RequireNonNil(s, "forEachFunc called on nil slice")
	object.RequireNonNil(f, "forEachFunc called on nil callfn")
	var wg sync.WaitGroup
	for _, r := range s {
		wg.Add(1)
		go func(rr interface{}) {
			defer wg.Done()
			f(rr)
		}(r)
	}
	wg.Wait()
	return
}
