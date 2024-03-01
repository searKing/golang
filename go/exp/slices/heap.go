// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices

import (
	"cmp"
)

// An MinHeap is a min-heap of slices.
type MinHeap[E cmp.Ordered] []E

func (h MinHeap[E]) Len() int           { return len(h) }
func (h MinHeap[E]) Less(i, j int) bool { return h[i] < h[j] }
func (h MinHeap[E]) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MinHeap[E]) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(E))
}

func (h *MinHeap[E]) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// An MaxHeap is a max-heap of slices.
type MaxHeap[E cmp.Ordered] []E

func (h MaxHeap[E]) Len() int           { return len(h) }
func (h MaxHeap[E]) Less(i, j int) bool { return h[i] > h[j] }
func (h MaxHeap[E]) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MaxHeap[E]) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(E))
}

func (h *MaxHeap[E]) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type Heap[S ~[]E, E any] struct {
	S S

	Comparator func(v1 E, v2 E) int
}

func NewHeapMin[S ~[]E, E cmp.Ordered](s S) *Heap[S, E] {
	return &Heap[S, E]{
		S:          s,
		Comparator: cmp.Compare[E],
	}
}

func NewHeapMax[S ~[]E, E cmp.Ordered](s S) *Heap[S, E] {
	return &Heap[S, E]{
		S: s,
		Comparator: func(v1 E, v2 E) int {
			return cmp.Compare[E](v2, v1)
		},
	}
}

func NewHeapFunc[S ~[]E, E any](s S, cmp func(a E, b E) int) *Heap[S, E] {
	return &Heap[S, E]{
		S:          s,
		Comparator: cmp,
	}
}

func (h Heap[S, E]) Len() int { return len(h.S) }

func (h Heap[S, E]) Less(i, j int) bool {
	if h.Comparator == nil { // nop, don't sort
		return true
	}
	return h.Comparator(h.S[i], h.S[j]) < 0
}

func (h Heap[S, E]) Swap(i, j int) { h.S[i], h.S[j] = h.S[j], h.S[i] }

func (h *Heap[S, E]) Push(x any) { // add x as element Len()
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	h.S = append(h.S, x.(E))
}

func (h *Heap[S, E]) Pop() any { // remove and return element Len() - 1.
	old := h.S
	n := len(old)
	x := old[n-1]
	h.S = old[0 : n-1]
	return x
}
