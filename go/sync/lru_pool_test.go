// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync

import (
	"context"
	"errors"
	"runtime"
	"runtime/debug"
	"testing"
	"time"
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

func TestLruPool_ReqPassedToNew(t *testing.T) {
	p := LruPool{
		New: func(ctx context.Context, req any) (any, error) { return req, nil },
	}
	v, put, err := p.GetOrError(context.Background(), "a")
	if err != nil || v != "a" {
		t.Fatalf("GetOrError = %v, %v; want a, nil", v, err)
	}
	put()
	v, put = p.GetByKey(context.Background(), "k", "b")
	if v != "b" {
		t.Fatalf("GetByKey = %v, want b", v)
	}
	put()
}

func TestLruPool_NewError(t *testing.T) {
	errNew := errors.New("new")
	p := LruPool{
		New: func(ctx context.Context, req any) (any, error) { return nil, errNew },
	}
	v, put, err := p.GetOrError(context.Background(), "a")
	if !errors.Is(err, errNew) || v != nil {
		t.Fatalf("GetOrError = %v, %v; want nil, %v", v, err, errNew)
	}
	put() // must not panic
	v, put = p.Get(context.Background(), "a")
	if v != nil {
		t.Fatalf("Get = %v, want nil", v)
	}
	put()
}

func TestLruPool_MaxIdleResourcesPerBucket(t *testing.T) {
	for _, tt := range []struct {
		maxIdle int
		reused  int // of 3 released resources, how many are reused
	}{
		{0, DefaultMaxIdleResourcesPerBucket},
		{1, 1},
		{-1, 0},
	} {
		i := 0
		p := LruPool{
			New: func(ctx context.Context, req any) (any, error) {
				i++
				return i, nil
			},
			MaxIdleResourcesPerBucket: tt.maxIdle,
		}
		var puts []context.CancelFunc
		for range 3 {
			_, put := p.Get(context.Background(), "")
			puts = append(puts, put)
		}
		for _, put := range puts {
			put()
		}
		seen := map[any]bool{}
		for range 3 {
			v, put := p.Get(context.Background(), "")
			seen[v] = true
			defer put()
		}
		reused := 0
		for v := range seen {
			if v.(int) <= 3 {
				reused++
			}
		}
		if reused != tt.reused {
			t.Errorf("MaxIdleResourcesPerBucket=%d: reused %d, want %d", tt.maxIdle, reused, tt.reused)
		}
	}
}

func TestLruPool_MaxResourcesPerBucket(t *testing.T) {
	p := LruPool{
		New:                   func(ctx context.Context, req any) (any, error) { return new(int), nil },
		MaxResourcesPerBucket: 1,
	}
	v, put := p.Get(context.Background(), "a")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	if _, _, err := p.GetOrError(ctx, "a"); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("GetOrError over limit: err = %v, want DeadlineExceeded", err)
	}
	// Other buckets are not limited.
	_, putB := p.Get(context.Background(), "b")
	putB()
	put()
	v2, put2 := p.Get(context.Background(), "a")
	defer put2()
	if v2 != v {
		t.Error("released resource not reused")
	}
}

func TestLruPool_CloseIdleResources(t *testing.T) {
	i := 0
	p := LruPool{
		New: func(ctx context.Context, req any) (any, error) {
			i++
			return i, nil
		},
	}
	_, put := p.Get(context.Background(), "")
	put()
	p.CloseIdleResources()
	v, put := p.Get(context.Background(), "")
	defer put()
	if v != 2 {
		t.Errorf("Get after CloseIdleResources = %v, want a new resource 2", v)
	}
}

func TestLruPool_DisableKeepAlives(t *testing.T) {
	i := 0
	p := LruPool{
		New: func(ctx context.Context, req any) (any, error) {
			i++
			return i, nil
		},
		DisableKeepAlives: true,
	}
	_, put := p.Get(context.Background(), "")
	put()
	v, put := p.Get(context.Background(), "")
	defer put()
	if v != 2 {
		t.Errorf("Get = %v, want a new resource 2", v)
	}
}

func TestLruPool_RemoveIdleResource(t *testing.T) {
	i := 0
	p := LruPool{
		New: func(ctx context.Context, req any) (any, error) {
			i++
			return i, nil
		},
	}
	pc, err := p.GetByKeyOrError(context.Background(), "", "")
	if err != nil {
		t.Fatal(err)
	}
	if p.RemoveIdleResource(pc) {
		t.Error("RemoveIdleResource of a leased resource = true, want false")
	}
	pc.Put()
	if !p.RemoveIdleResource(pc) {
		t.Error("RemoveIdleResource of an idle resource = false, want true")
	}
	if p.RemoveIdleResource(pc) {
		t.Error("second RemoveIdleResource = true, want false")
	}
	v, put := p.Get(context.Background(), "")
	defer put()
	if v != 2 {
		t.Errorf("Get = %v, want a new resource 2", v)
	}
}

func TestLruPool_ZeroPersistResource(t *testing.T) {
	var pc PersistResource
	if v := pc.Get(); v != nil {
		t.Errorf("Get = %v, want nil", v)
	}
	pc.Put()
	var p LruPool
	if p.RemoveIdleResource(&pc) {
		t.Error("RemoveIdleResource = true, want false")
	}
}

func TestLruPool_CloseIdleResourcesThenPut(t *testing.T) {
	i := 0
	p := LruPool{
		New: func(ctx context.Context, req any) (any, error) {
			i++
			return i, nil
		},
	}
	_, put := p.Get(context.Background(), "")
	p.CloseIdleResources()
	put() // in use during CloseIdleResources: not kept
	v, put := p.Get(context.Background(), "")
	if v != 2 {
		t.Errorf("Get = %v, want a new resource 2", v)
	}
	put() // Get makes the pool keep idle resources again
	v, put = p.Get(context.Background(), "")
	defer put()
	if v != 2 {
		t.Errorf("Get = %v, want the idle resource 2", v)
	}
}

func TestLruPool_RemoveIdleResourceOfOtherPool(t *testing.T) {
	var a, b LruPool
	pc, err := a.GetByKeyOrError(context.Background(), "", "")
	if err != nil {
		t.Fatal(err)
	}
	pc.Put()
	if b.RemoveIdleResource(pc) {
		t.Error("RemoveIdleResource of a resource of another pool = true, want false")
	}
	if !a.RemoveIdleResource(pc) {
		t.Error("RemoveIdleResource = false, want true")
	}
}

func TestLruPool_NegativeMaxIdleResources(t *testing.T) {
	i := 0
	p := LruPool{
		New: func(ctx context.Context, req any) (any, error) {
			i++
			return i, nil
		},
		MaxIdleResources: -1,
	}
	_, put := p.Get(context.Background(), "")
	put()
	v, put := p.Get(context.Background(), "")
	defer put()
	if v != 2 {
		t.Errorf("Get = %v, want a new resource 2", v)
	}
}
