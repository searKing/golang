// Copyright 2025 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync_test

import (
	"runtime"
	"runtime/debug"
	"testing"

	sync_ "github.com/searKing/golang/go/exp/sync"
)

func TestPool(t *testing.T) {
	// disable GC so we can control when it happens.
	defer debug.SetGCPercent(debug.SetGCPercent(-1))
	var p sync_.Pool[any]
	if p.Get() != nil {
		t.Fatal("expected empty")
	}

	// Make sure that the goroutine doesn't migrate to another P
	// between Put and Get calls.
	runtime.LockOSThread()
	p.Put("a")
	p.Put("b")
	if g := p.Get(); g != "a" {
		t.Fatalf("got %#v; want a", g)
	}
	if g := p.Get(); g != "b" {
		t.Fatalf("got %#v; want b", g)
	}
	if g := p.Get(); g != nil {
		t.Fatalf("got %#v; want nil", g)
	}
	runtime.UnlockOSThread()

	// Put in a large number of objects so they spill into
	// stealable space.
	for i := 0; i < 100; i++ {
		p.Put("c")
	}
	// After one GC, the victim cache should keep them alive.
	runtime.GC()
	if g := p.Get(); g != "c" {
		t.Fatalf("got %#v; want c after GC", g)
	}
	// A second GC should drop the victim cache.
	runtime.GC()
	if g := p.Get(); g != nil {
		t.Fatalf("got %#v; want nil after second GC", g)
	}
}

func TestPoolNew(t *testing.T) {
	// disable GC so we can control when it happens.
	defer debug.SetGCPercent(debug.SetGCPercent(-1))

	i := 0
	p := sync_.Pool[int]{
		New: func() int {
			i++
			return i
		},
	}
	if v := p.Get(); v != 1 {
		t.Fatalf("got %v; want 1", v)
	}
	if v := p.Get(); v != 2 {
		t.Fatalf("got %v; want 2", v)
	}

	// Make sure that the goroutine doesn't migrate to another P
	// between Put and Get calls.
	runtime.LockOSThread()
	p.Put(42)
	if v := p.Get(); v != 42 {
		t.Fatalf("got %v; want 42", v)
	}
	runtime.UnlockOSThread()

	if v := p.Get(); v != 3 {
		t.Fatalf("got %v; want 3", v)
	}
}
