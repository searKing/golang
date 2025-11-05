// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package buffer

import (
	"bytes"
	"encoding"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
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

// AppendTextMarshal writes string represents v to the buffer.
func (b *Buffer) AppendTextMarshal(v any) error {
	if tm, ok := v.(encoding.TextMarshaler); ok {
		data, err := tm.MarshalText()
		if err != nil {
			return err
		}
		b.Write(data)
		return nil
	}
	if bs, ok := byteSlice(v); ok {
		b.WriteString(strconv.Quote(string(bs)))
		return nil
	}
	b.WriteString(fmt.Sprintf("%+v", v))
	return nil
}

// byteSlice returns its argument as a []byte if the argument's
// underlying type is []byte, along with a second return value of true.
// Otherwise it returns nil, false.
func byteSlice(a any) ([]byte, bool) {
	if bs, ok := a.([]byte); ok {
		return bs, true
	}
	// Like Printf's %s, we allow both the slice type and the byte element type to be named.
	t := reflect.TypeOf(a)
	if t != nil && t.Kind() == reflect.Slice && t.Elem().Kind() == reflect.Uint8 {
		return reflect.ValueOf(a).Bytes(), true
	}
	return nil, false
}
