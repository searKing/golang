// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stack

import (
	"reflect"
	"testing"
)

func TestStack_Push(t *testing.T) {
	s := New()
	var input int = 0
	ele := s.Push(input)
	e, ok := ele.Value.(int)
	if !ok {
		t.Errorf("type must be %s", reflect.TypeOf(input).String())
	}
	if e != input {
		t.Errorf("value must be %v", input)
	}
	if s.Len() != 1 {
		t.Errorf("len must be %v", 1)
	}
}

func TestStack_Pop(t *testing.T) {
	s := New()
	inputs := []int{2, 1, 0}
	for _, input := range inputs {
		s.Push(input)
	}
	for idx := 0; idx < len(inputs); idx++ {
		ele := s.Pop()
		e, ok := ele.Value.(int)
		if !ok {
			t.Errorf("type must be %s", reflect.TypeOf(inputs[len(inputs)-idx-1]).String())
		}
		if e != idx {
			t.Errorf("value must be %v", idx)
		}
		if s.Len() != len(inputs)-idx-1 {
			t.Errorf("len is %v must be %v", s.Len(), len(inputs)-idx)
		}
	}

}
func TestStack_Peek(t *testing.T) {
	s := New()
	var input int = 0
	s.Push(input)
	ele := s.Peek()
	e, ok := ele.Value.(int)
	if !ok {
		t.Errorf("return type must be %s", reflect.TypeOf(input).String())
	}
	if e != input {
		t.Errorf("value must be %v", input)
	}
	if s.Len() != 1 {
		t.Errorf("len must be %v", 1)
	}
}
