// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package queue_test

import (
	"fmt"

	"github.com/searKing/golang/go/exp/container/queue"
)

func Example() {
	// Create a new queue and put some numbers in it.
	var l1 queue.Queue[int]
	l1.PushBack(1)
	l1.PushBack(2)
	l1.PushBack(3)
	l1.PushBack(4)

	var l2 queue.Queue[int]
	l2.PushBack(5)
	l2.PushBack(6)
	l2.PushBack(7)
	l2.PushBack(8)

	l1.PushBackSeq(l2.Values())

	var s []int
	if cleaned := l2.TrimFrontFunc(func(e int) bool {
		if e%2 == 1 {
			s = append(s, e)
			return true
		}
		return false
	}); cleaned {
		fmt.Printf("l2: clean leading: %v\n", s)
	}
	// Iterate through queue and print its contents.
	for i, e := range l1.All() {
		fmt.Printf("l1[%d]: %d\n", i, e)
	}
	var i int
	for e := range l2.Values() {
		fmt.Printf("l2[%d]: %d\n", i, e)
		i++
	}

	// Output:
	// l2: clean leading: [5]
	// l1[0]: 1
	// l1[1]: 2
	// l1[2]: 3
	// l1[3]: 4
	// l1[4]: 5
	// l1[5]: 6
	// l1[6]: 7
	// l1[7]: 8
	// l2[0]: 6
	// l2[1]: 7
	// l2[2]: 8
}
