// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slice

import (
	"sync"

	"github.com/searKing/golang/go/util/object"
)

// FindAnyFunc returns an {@link interface{}} describing some element of the stream, or an
// empty {@code Optional} if the stream is empty.
func FindAnyIndexFunc(s interface{}, f func(interface{}) bool) int {
	return findAnyIndexFunc(Of(s), f, true)
}

// findAnyFunc is the same as FindAnyFunc.
func findAnyIndexFunc(s []interface{}, f func(interface{}) bool, truth bool) int {
	object.RequireNonNil(s, "findAnyIndexFunc called on nil slice")
	object.RequireNonNil(f, "findAnyIndexFunc called on nil callfn")
	var findc chan int
	findc = make(chan int)
	defer close(findc)
	var mu sync.Mutex
	var wg sync.WaitGroup
	var found bool
	hasFound := func() bool {
		mu.Lock()
		defer mu.Unlock()
		return found
	}
	for idx, r := range s {
		if hasFound() {
			break
		}

		wg.Add(1)
		go func(rr interface{}) {
			defer wg.Done()
			foundYet := func() bool {
				mu.Lock()
				defer mu.Unlock()
				return found
			}()
			if foundYet {
				return
			}
			if f(rr) == truth {
				mu.Lock()
				defer mu.Unlock()
				if found {
					return
				}
				found = true
				findc <- idx
				return
			}
		}(r)
	}
	go func() {
		defer close(findc)
		wg.Done()
	}()
	out, ok := <-findc
	if ok {
		return out
	}
	return -1
}
