// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package queue_test

import (
	"testing"

	"github.com/searKing/golang/go/exp/container/queue"
)

// For debugging - keep around.
func dump(t *testing.T, r *queue.Queue[int]) {
	if r == nil {
		t.Logf("empty")
		return
	}
	var i int
	r.Do(func(e int) {
		t.Logf("%4d: %d\n", i, e)
		i++
	})
	t.Logf("\n")
}

func verify(t *testing.T, r *queue.Queue[int], N int, sum int) {
	// Len
	n := r.Len()
	if n != N {
		t.Errorf("r.Len() == %d; expected %d", n, N)
	}

	// iteration
	n = 0
	s := 0
	r.Do(func(p int) {
		n++
		s += p
	})
	if n != N {
		t.Errorf("number of forward iterations == %d; expected %d", n, N)
	}
	if sum >= 0 && s != sum {
		t.Errorf("forward queue sum = %d; expected %d", s, sum)
	}

	if r == nil {
		return
	}
}

func TestCornerCases(t *testing.T) {
	var (
		r0 queue.Queue[int]
		r1 queue.Queue[int]
	)
	r0.PushBack(1)

	// Basics
	verify(t, &r0, 1, 1)
	verify(t, &r1, 0, 0)
	// Insert
	r1.PushBackQueue(&r0)
	r1.ShrinkToFit()
	verify(t, &r0, 1, 1)
	verify(t, &r1, 1, 1)
	// Insert
	r1.PushBackQueue(&r0)
	r1.ShrinkToFit()
	verify(t, &r0, 1, 1)
	verify(t, &r1, 2, 2)
	// Remove
	r1.PopFront()
	r1.ShrinkToFit()
	verify(t, &r0, 1, 1)
	verify(t, &r1, 1, 1)
	// Remove
	r1.PopFront()
	r1.ShrinkToFit()
	verify(t, &r0, 1, 1)
	verify(t, &r1, 0, 0)
}

func makeN(n int) *queue.Queue[int] {
	r := queue.New[int]()
	for i := 1; i <= n; i++ {
		r.PushBack(i)
	}
	return r
}

func sumN(n int) int { return (n*n + n) / 2 }

func TestNew(t *testing.T) {
	for i := 0; i < 10; i++ {
		r := queue.New[int]()
		verify(t, r, 0, -1)
	}
	for i := 0; i < 10; i++ {
		r := makeN(i)
		verify(t, r, i, sumN(i))
	}
}
