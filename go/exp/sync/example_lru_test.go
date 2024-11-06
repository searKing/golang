// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync_test

import (
	"fmt"

	sync_ "github.com/searKing/golang/go/exp/sync"
)

func ExampleNewLRU() {
	l := sync_.NewLRU[int, int](2)
	var evictCounter int
	l.SetEvictCallback(func(k int, v int) {
		evictCounter++
		fmt.Printf("{%d,%d} evicted as oldest: %d\n", k, v, evictCounter)
	})
	printAll := func() {
		fmt.Printf("lru: ")
		first := true
		for k, v := range l.All() {
			if !first {
				fmt.Printf(", ")
			}
			first = false
			fmt.Printf("{%d,%d}", k, v)
		}
		fmt.Println()
	}
	l.Add(1, 1)
	l.Add(2, 2)
	l.Add(3, 3) // evict 1, that's the oldest
	printAll()

	fmt.Println("try to refresh {2,22}")
	l.Add(2, 22) // refresh 2; Now 3 is the oldest
	printAll()

	fmt.Println("try to remove oldest")
	l.RemoveOldest() // evict 2, that's the oldest
	printAll()

	fmt.Println("try to purge all elements")
	l.Purge()
	printAll()

	// Output:
	// {1,1} evicted as oldest: 1
	// lru: {2,2}, {3,3}
	// try to refresh {2,22}
	// lru: {3,3}, {2,22}
	// try to remove oldest
	// {3,3} evicted as oldest: 2
	// lru: {2,22}
	// try to purge all elements
	// {2,22} evicted as oldest: 3
	// lru:
}
