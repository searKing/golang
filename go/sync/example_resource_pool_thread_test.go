// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/searKing/golang/go/sync"
)

// model stands for a handle of a non-Go library, such as a CUDA context,
// that must be created, used and freed on the OS thread that created it.
type model struct {
	gpu   int
	freed bool
}

func newModel(gpu int) (*model, error) { return &model{gpu: gpu}, nil }

func (m *model) infer(input string) (string, error) {
	return fmt.Sprintf("%s@gpu%d", input, m.gpu), nil
}

func (m *model) free() error {
	m.freed = true
	return nil
}

// boundModel is a model bound to the Thread that created it.
type boundModel struct {
	thread *sync.Thread
	m      *model
}

// do calls f with the model on its Thread.
func (b *boundModel) do(ctx context.Context, f func(m *model) error) error {
	var err error
	if terr := b.thread.Do(ctx, func(context.Context) { err = f(b.m) }); terr != nil {
		return terr
	}
	return err
}

// newModelPool returns a pool of at most 2 models per GPU,
// each created on, used on and freed on a Thread of its own.
func newModelPool() *sync.ResourcePool[int, *boundModel] {
	return &sync.ResourcePool[int, *boundModel]{
		New: func(ctx context.Context, gpu int) (*boundModel, error) {
			b := &boundModel{thread: new(sync.Thread)}
			var err error
			if terr := b.thread.Do(ctx, func(context.Context) { b.m, err = newModel(gpu) }); terr != nil {
				err = terr
			}
			if err != nil {
				b.thread.Close()
				return nil, err
			}
			return b, nil
		},
		Close: func(gpu int, b *boundModel) error {
			defer b.thread.Close()
			return b.do(context.Background(), (*model).free)
		},
		MaxResourcesPerKey: 2,
	}
}

// This example pools handles of a non-Go library that must stay on the OS thread
// that created them: New creates each handle on a Thread of its own,
// every use goes through that Thread, and Close frees it there.
func ExampleResourcePool_thread() {
	p := newModelPool()
	defer p.Shutdown(context.Background())

	// Stateless: one call on a model of GPU 0 or 1.
	err := p.Do(context.Background(), func(ctx context.Context, b *boundModel) error {
		return b.do(ctx, func(m *model) error {
			out, err := m.infer("cat.jpg")
			fmt.Println(out)
			return err
		})
	}, 0, 1)
	if errors.Is(err, sync.ErrThreadClosed) {
		fmt.Println("thread closed")
	}

	// Stateful: several calls on the same model.
	l, err := p.Get(context.Background(), 1)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Release()
	for _, in := range []string{"frame1", "frame2"} {
		_ = l.Value().do(context.Background(), func(m *model) error {
			out, err := m.infer(in)
			fmt.Println(out)
			return err
		})
	}

	// Output:
	// cat.jpg@gpu0
	// frame1@gpu1
	// frame2@gpu1
}
