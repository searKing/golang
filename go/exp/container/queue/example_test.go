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
	l := queue.New[int]()
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)
	l.PushBack(4)

	l2 := queue.New[int]()
	l2.PushBack(5)
	l2.PushBack(6)
	l2.PushBack(7)
	l2.PushBack(8)

	l.PushBackQueue(l2)

	l2.TrimFrontFunc(func(e int) bool {
		return e%2 == 1
	})
	// Iterate through list and print its contents.
	l.Do(func(e int) {
		fmt.Println(e)
	})
	l2.Do(func(e int) {
		fmt.Println(e)
	})

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
