// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package buffer

import (
	"bytes"
	"encoding/json"
	"sync"
)

// Having an initial size gives a dramatic speedup.
var bufPool = sync.Pool{
	New: func() any { return new(Buffer) },
}

type Buffer struct {
	bytes.Buffer
}

func New() *Buffer {
	return bufPool.Get().(*Buffer)
}

func (b *Buffer) Free() {
	// To reduce peak allocation, return only smaller buffers to the pool.
	const maxBufferSize = 16 << 10
	if b.Cap() <= maxBufferSize {
		b.Buffer.Reset()
		bufPool.Put(b)
	}
}

func (b *Buffer) WritePosInt(i int) {
	b.WritePosIntWidth(i, 0)
}

// WritePosIntWidth writes non-negative integer i to the buffer, padded on the left
// by zeroes to the given width. Use a width of 0 to omit padding.
func (b *Buffer) WritePosIntWidth(i, width int) {
	// Cheap integer to fixed-width decimal ASCII.
	// Copied from log/log.go.

	if i < 0 {
		panic("negative int")
	}

	// Assemble decimal in reverse order.
	var bb [20]byte
	bp := len(bb) - 1
	for i >= 10 || width > 1 {
		width--
		q := i / 10
		bb[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	bb[bp] = byte('0' + i)
	b.Write(bb[bp:])
}

// AppendJSONMarshal writes string represents v to the buffer.
func (b *Buffer) AppendJSONMarshal(v any) error {
	// Use a json.Encoder to avoid escaping HTML.
	var bb bytes.Buffer
	enc := json.NewEncoder(&bb)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return err
	}
	bs := bb.Bytes()
	b.Write(bs[:len(bs)-1]) // remove final newline
	return nil
}
