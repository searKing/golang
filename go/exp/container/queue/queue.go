// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package queue

import (
	"iter"
	"slices"
)

// A Queue is a queue of FIFO, not a deque.
type Queue[E any] struct {
	// This is a queue, not a deque.
	// It is split into two stages - head[headPos:] and tail.
	// PopFront is trivial (headPos++) on the first stage, and
	// PushBack is trivial (append) on the second stage.
	// If the first stage is empty, PopFront can swap the
	// first and second stages to remedy the situation.
	//
	// This two-stage split is analogous to the use of two lists
	// in Okasaki's purely functional queue but without the
	// overhead of reversing the list when swapping stages.
	head    []E
	headPos int
	tail    []E
}

// Len returns the number of items in the queue.
func (q *Queue[E]) Len() int {
	return len(q.head) - q.headPos + len(q.tail)
}

// PushBack adds w to the back of the queue.
func (q *Queue[E]) PushBack(w E) {
	q.tail = append(q.tail, w)
}

// PopFront removes and returns the item at the front of the queue.
func (q *Queue[E]) PopFront() E {
	v, _ := q.popFront()
	return v
}

// Front returns the item at the front of the queue without removing it.
func (q *Queue[E]) Front() E {
	v, _ := q.front()
	return v
}

// Next reports whether there are more iterations to execute.
// Every call to PopFront, even the first one, must be preceded by a call to Next.
func (q *Queue[E]) Next() bool {
	return q.Len() > 0
}

func (q *Queue[E]) popFront() (E, bool) {
	if q.headPos >= len(q.head) {
		if len(q.tail) == 0 {
			var zeroE E
			return zeroE, false
		}
		// Pick up tail as new head, clear tail.
		q.head, q.headPos, q.tail = q.tail, 0, q.head[:0]
	}
	w := q.head[q.headPos]
	var zeroE E
	q.head[q.headPos] = zeroE
	q.headPos++
	return w, true
}

func (q *Queue[E]) front() (E, bool) {
	if q.headPos < len(q.head) {
		return q.head[q.headPos], true
	}
	if len(q.tail) > 0 {
		return q.tail[0], true
	}
	var zeroE E
	return zeroE, false
}

// PushBackSeq appends the values from seq to the queue.
func (q *Queue[E]) PushBackSeq(seq iter.Seq[E]) {
	for e := range seq {
		q.PushBack(e)
	}
}

// All returns an iterator over index-value pairs in the queue
// in the usual order.
func (q *Queue[E]) All() iter.Seq2[int, E] {
	return func(yield func(int, E) bool) {
		var i int
		for v := range q.Values() {
			if !yield(i, v) {
				return
			}
			i++
		}
	}
}

// Values returns an iterator that yields the queue elements in order.
func (q *Queue[E]) Values() iter.Seq[E] {
	return q.Range
}

// Range calls f sequentially for each value present in the queue[E] in forward order.
// If f returns false, range stops the iteration.
func (q *Queue[E]) Range(f func(e E) bool) {
	for i := q.headPos; i < len(q.head); i++ {
		if !f(q.head[i]) {
			return
		}
	}
	for i := 0; i < len(q.tail); i++ {
		if !f(q.tail[i]) {
			return
		}
	}
}

// ShrinkToFit requests the removal of unused capacity.
func (q *Queue[E]) ShrinkToFit() {
	if q.headPos >= len(q.head) { // empty head
		if len(q.tail) == 0 { // empty tail
			q.head = nil
			q.headPos = 0
			return
		}
		// Pick up tail as new head, clear tail.
		q.head, q.headPos, q.tail = q.tail, 0, nil
		return
	}
	// shrink slice head actually
	if len(q.head) > 1 {
		q.head = append([]E{}, q.head[q.headPos:]...)
		q.headPos = 0
	}
	return
}

// TrimFrontFunc pops all leading elem that satisfying f(c) from the head of the
// queue, reporting whether any were popped.
func (q *Queue[E]) TrimFrontFunc(f func(e E) bool) (cleaned bool) {
	for q.Next() {
		if !f(q.Front()) {
			return cleaned
		}
		q.PopFront()
		cleaned = true
	}
	return cleaned
}

// New returns an initialized queue.
// Deprecated: Use var q Queue[E] instead.
func New[E any](e ...E) *Queue[E] {
	var q Queue[E]
	q.PushBackSeq(slices.Values(e))
	return &q
}

// PushBackQueue inserts a copy of another queue at the back of queue l.
// Deprecated: Use [PushBackSeq] instead.
func (q *Queue[E]) PushBackQueue(other *Queue[E]) {
	q.PushBackSeq(other.Values())
}

// Do calls function f on each element of the queue without removing it, in forward order.
// The behavior of Do is undefined if f changes *q.
// Deprecated: Use [Range] or [All] instead.
func (q *Queue[E]) Do(f func(E)) {
	q.Range(func(e E) bool {
		f(e)
		return true
	})
}
