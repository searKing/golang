package bufio

import (
	"bufio"
	"io"
	"sync"
)

type WriterPool struct {
	pool sync.Pool
	size int
}

func NewWriterPool() *WriterPool {
	return &WriterPool{}
}
func NewWriterPoolSize(size int) *WriterPool {
	return &WriterPool{
		size: size,
	}
}

func (p *WriterPool) Put(br *bufio.Writer) {
	br.Reset(nil)
	p.pool.Put(br)
}

func (p *WriterPool) Get(w io.Writer) *bufio.Writer {
	if v := p.pool.Get(); v != nil {
		bw := v.(*bufio.Writer)
		bw.Reset(w)
		return bw
	}
	// Note: if this reader size is ever changed, update
	// TestHandlerBodyClose's assumptions.
	return bufio.NewWriterSize(w, p.size)
}

func (p *WriterPool) Clear() {
	// Get One elem in the pool
	if p.pool.Get() == nil {
		// The pool is empty
		return
	}
	p.Clear()
	return
}
