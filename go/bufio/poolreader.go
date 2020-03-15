// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bufio

import (
	"bufio"
	"io"
	"sync"
)

type ReaderPool struct {
	pool sync.Pool
	size int
}

func NewReaderPool() *ReaderPool {
	return &ReaderPool{}
}

func NewReaderPoolSize(size int) *ReaderPool {
	return &ReaderPool{
		size: size,
	}
}

func (p *ReaderPool) Put(br *bufio.Reader) {
	br.Reset(nil)
	p.pool.Put(br)
}

func (p *ReaderPool) Get(r io.Reader) *bufio.Reader {
	if v := p.pool.Get(); v != nil {
		br := v.(*bufio.Reader)
		br.Reset(r)
		return br
	}
	// Note: if this reader size is ever changed, update
	// TestHandlerBodyClose's assumptions.
	return bufio.NewReaderSize(r, p.size)
}

func (p *ReaderPool) Clear() {
	// Get One elem in the pool
	// for p.pool.New is nil, so p.pool.New will return nil if empty
	if p.pool.Get() == nil {
		// The pool is empty
		return
	}
	p.Clear()
	return
}
