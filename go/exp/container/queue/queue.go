// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package queue

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

// New returns an initialized list.
func New[E any]() *Queue[E] { return new(Queue[E]) }

// Len returns the number of items in the queue.
func (q *Queue[E]) Len() int {
	n := 0
	if q != nil {
		n = len(q.head) - q.headPos + len(q.tail)
	}
	return n
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

// PushBackQueue inserts a copy of another queue at the back of queue l.
// The queues l and other may be the same. They must not be nil.
func (q *Queue[E]) PushBackQueue(other *Queue[E]) {
	other.Do(func(a E) {
		q.PushBack(a)
	})
}

// Do calls function f on each element of the queue without removing it, in forward order.
// The behavior of Do is undefined if f changes *q.
func (q *Queue[E]) Do(f func(E)) {
	if q != nil {
		for i := q.headPos; i < len(q.head); i++ {
			f(q.head[i])
		}
		for i := 0; i < len(q.tail); i++ {
			f(q.tail[i])
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
