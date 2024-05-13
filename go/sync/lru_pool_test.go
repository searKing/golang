// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync

import (
	"context"
	"runtime"
	"runtime/debug"
	"testing"
)

func TestLruPool(t *testing.T) {
	// disable GC so we can control when it happens.
	defer debug.SetGCPercent(debug.SetGCPercent(-1))
	var p LruPool
	g, put := p.Get(context.Background(), "")
	defer put()
	if g != nil {
		t.Fatal("expected empty")
	}

	// Make sure that the goroutine doesn't migrate to another P
	// between Put and Get calls.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	g, put = p.Get(context.Background(), "")
	defer put()
	if g != nil {
		t.Fatal("expected empty")
	}
}

func TestLruPool_New(t *testing.T) {
	// disable GC so we can control when it happens.
	defer debug.SetGCPercent(debug.SetGCPercent(-1))

	i := 0
	p := LruPool{
		New: func(ctx context.Context, req any) (resp any, err error) {
			i++
			return i, nil
		},
	}
	a, puta := p.Get(context.Background(), "")
	if a != 1 {
		t.Fatalf("got %v; want 1", a)
	}
	b, putb := p.Get(context.Background(), "")
	if b != 2 {
		t.Fatalf("got %v; want 2", b)
	}

	// Make sure that the goroutine doesn't migrate to another P
	// between Put and Get calls.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	puta()
	a, puta = p.Get(context.Background(), "")
	defer puta()
	if a != 1 {
		t.Fatalf("got %v; want 1", a)
	}
	c, putc := p.Get(context.Background(), "")
	defer putc()
	if c != 3 {
		t.Fatalf("got %v; want 3", a)
	}
	putb()
	b, putb = p.Get(context.Background(), "")
	defer putb()

}
