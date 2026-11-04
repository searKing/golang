// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync_test

import (
	"context"
	"errors"
	"testing"

	sync_ "github.com/searKing/golang/go/sync"
)

// TestResourcePool_Thread checks the pool of ExampleResourcePool_thread
// uses each model on a single thread, and frees it and its Thread on Shutdown.
func TestResourcePool_Thread(t *testing.T) {
	p := newModelPool()
	threadOf := func(ctx context.Context, b *boundModel) uint64 {
		var id uint64
		if err := b.do(ctx, func(*model) error { id = sync_.CallerID(); return nil }); err != nil {
			t.Fatal(err)
		}
		return id
	}

	var first *boundModel
	var id uint64
	if err := p.Do(context.Background(), func(ctx context.Context, b *boundModel) error {
		first, id = b, threadOf(ctx, b)
		return nil
	}, 0); err != nil {
		t.Fatal(err)
	}
	if id == 0 || id == sync_.CallerID() {
		t.Fatalf("model used on thread %d, want a Thread of its own", id)
	}

	l, err := p.Get(context.Background(), 0)
	if err != nil {
		t.Fatal(err)
	}
	if l.Value() != first {
		t.Fatal("idle model not reused")
	}
	for range 3 {
		if got := threadOf(context.Background(), l.Value()); got != id {
			t.Errorf("model used on thread %d, then %d", id, got)
		}
	}
	l.Release()

	if err := p.Shutdown(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !first.m.freed {
		t.Error("model not freed on Shutdown")
	}
	if err := first.do(context.Background(), func(*model) error { return nil }); !errors.Is(err, sync_.ErrThreadClosed) {
		t.Errorf("Thread not closed after Shutdown: err = %v", err)
	}
}
