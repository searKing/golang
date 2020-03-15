// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stack

import "container/list"

// take advantage of list
type (
	Stack struct {
		list list.List
	}
	// Element is an element of a linked list.
	Element list.Element
)

// Create a new stack
func New() *Stack {
	return &Stack{}
}

// Return the number of items in the stack
func (s *Stack) Len() int {
	return s.list.Len()
}

// View the top item on the stack
func (s *Stack) Peek() *Element {
	return (*Element)(s.list.Back())
}

// Pop the top item of the stack and return it
func (s *Stack) Pop() *Element {
	ele := s.list.Back()
	if ele == nil {
		return nil
	}
	s.list.Remove(ele)
	return (*Element)(ele)
}

// Push a value onto the top of the stack
func (s *Stack) Push(value interface{}) *Element {
	return (*Element)(s.list.PushBack(value))
}
