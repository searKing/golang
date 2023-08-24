// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package unsafe

import (
	"bytes"
	"math/rand"
	"testing"

	rand_ "github.com/searKing/golang/go/crypto/rand"
)

var testString = "Albert Einstein: Logic will get you from A to B. Imagination will take you everywhere."
var testBytes = []byte(testString)

func rawBytesToString(b []byte) string {
	return string(b)
}

func rawStringToBytes(s string) []byte {
	return []byte(s)
}

func TestUnsafeConversions(t *testing.T) {
	t.Parallel()

	// needs to be large to force allocations so we pick a random value between [1024, 2048]
	size := 1024 + rand.Intn(1024+1)

	t.Run("StringToBytes semantics", func(t *testing.T) {
		t.Parallel()

		s := rand_.String(size)
		b := StringToBytes(s)
		if len(b) != size {
			t.Errorf("unexpected length: %d", len(b))
		}
		if cap(b) != size {
			t.Errorf("unexpected capacity: %d", cap(b))
		}
		if !bytes.Equal(b, []byte(s)) {
			t.Errorf("unexpected equality failure: %#v", b)
		}
	})

	t.Run("StringToBytes allocations", func(t *testing.T) {
		t.Parallel()

		s := rand_.String(size)
		f := func() {
			b := StringToBytes(s)
			if len(b) != size {
				t.Errorf("invalid length: %d", len(b))
			}
		}
		allocs := testing.AllocsPerRun(100, f)
		if allocs > 0 {
			t.Errorf("expected zero allocations, got %v", allocs)
		}
	})

	t.Run("BytesToString semantics", func(t *testing.T) {
		t.Parallel()

		b := make([]byte, size)
		if _, err := rand.Read(b); err != nil {
			t.Fatal(err)
		}
		s := BytesToString(b)
		if len(s) != size {
			t.Errorf("unexpected length: %d", len(s))
		}
		if s != string(b) {
			t.Errorf("unexpected equality failure: %#v", s)
		}
	})

	t.Run("BytesToString allocations", func(t *testing.T) {
		t.Parallel()

		b := make([]byte, size)
		if _, err := rand.Read(b); err != nil {
			t.Fatal(err)
		}
		f := func() {
			s := BytesToString(b)
			if len(s) != size {
				t.Errorf("invalid length: %d", len(s))
			}
		}
		allocs := testing.AllocsPerRun(100, f)
		if allocs > 0 {
			t.Errorf("expected zero allocations, got %v", allocs)
		}
	})
}

func BenchmarkBytesToStringRaw(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rawBytesToString(testBytes)
	}
}

func BenchmarkBytesToString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BytesToString(testBytes)
	}
}

func BenchmarkStringToBytesRaw(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rawStringToBytes(testString)
	}
}

func BenchmarkStringToBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		StringToBytes(testString)
	}
}
