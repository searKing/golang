// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package buffer_test

import (
	"testing"

	"github.com/searKing/golang/go/log/slog/internal/buffer"
)

func TestBuffer(t *testing.T) {
	b := buffer.New()
	defer b.Free()
	b.WriteString("hello")
	b.WriteByte(',')
	b.Write([]byte(" world"))
	b.WritePosIntWidth(17, 4)

	got := b.String()
	want := "hello, world0017"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestAlloc(t *testing.T) {
	got := int(testing.AllocsPerRun(5, func() {
		b := buffer.New()
		defer b.Free()
		b.WriteString("not 1K worth of bytes")
	}))
	if got != 0 {
		t.Errorf("got %d allocs, want 0", got)
	}
}
