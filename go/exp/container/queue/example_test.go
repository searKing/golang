// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package queue_test

import (
	"fmt"

	"github.com/searKing/golang/go/exp/container/queue"
)

func Example() {
	// Create a new list and put some numbers in it.
	var l queue.Queue[int]
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)
	l.PushBack(4)

	var l2 queue.Queue[int]
	l2.PushBack(5)
	l2.PushBack(6)
	l2.PushBack(7)
	l2.PushBack(8)

	l.PushBackSeq(l2.Values())

	l2.TrimFrontFunc(func(e int) bool {
		return e%2 == 1
	})
	// Iterate through list and print its contents.
	for e := range l.Values() {
		fmt.Println(e)
	}
	for e := range l2.Values() {
		fmt.Println(e)
	}

	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
	// 6
	// 7
	// 8
	// 6
	// 7
	// 8
}
