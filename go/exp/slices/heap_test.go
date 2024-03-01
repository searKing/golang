// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices_test

import (
	"cmp"
	"container/heap"
	"fmt"

	"github.com/searKing/golang/go/exp/slices"
)

// This example inserts several ints into an MinHeap, checks the minimum,
// and removes them in order of priority.
func Example_minHeap() {
	h := &slices.MinHeap[int]{2, 1, 5}
	heap.Init(h)
	heap.Push(h, 3)
	fmt.Printf("minimum: %d\n", (*h)[0])
	for h.Len() > 0 {
		fmt.Printf("%d ", heap.Pop(h))
	}

	// Output:
	// minimum: 1
	// 1 2 3 5
}

// This example inserts several ints into an MaxHeap, checks the maximum,
// and removes them in order of priority.
func Example_maxHeap() {
	h := &slices.MaxHeap[int]{2, 1, 5}
	heap.Init(h)
	heap.Push(h, 3)
	fmt.Printf("maximum: %d\n", (*h)[0])
	for h.Len() > 0 {
		fmt.Printf("%d ", heap.Pop(h))
	}

	// Output:
	// maximum: 5
	// 5 3 2 1
}

// This example inserts several ints into an Min Heap, checks the minimum,
// and removes them in order of priority.
func Example_heapMin() {
	h := slices.NewHeapMin([]int{2, 1, 5})
	heap.Init(h)
	heap.Push(h, 3)
	fmt.Printf("minimum: %d\n", h.S[0])
	for h.Len() > 0 {
		fmt.Printf("%d ", heap.Pop(h))
	}

	// Output:
	// minimum: 1
	// 1 2 3 5
}

// This example inserts several ints into an Max Heap, checks the maximum,
// and removes them in order of priority.
func Example_heapMax() {
	h := slices.NewHeapMax([]int{2, 1, 5})
	heap.Init(h)
	heap.Push(h, 3)
	fmt.Printf("maximum: %d\n", h.S[0])
	for h.Len() > 0 {
		fmt.Printf("%d ", heap.Pop(h))
	}

	// Output:
	// maximum: 5
	// 5 3 2 1
}

// This example inserts several ints into an Min Heap, checks the minimum,
// and removes them in order of priority.
func Example_heapFunc() {
	h := slices.NewHeapFunc([]int{2, 1, 5}, cmp.Compare[int])
	heap.Init(h)
	heap.Push(h, 3)
	fmt.Printf("minimum: %d\n", h.S[0])
	for h.Len() > 0 {
		fmt.Printf("%d ", heap.Pop(h))
	}

	// Output:
	// minimum: 1
	// 1 2 3 5
}
