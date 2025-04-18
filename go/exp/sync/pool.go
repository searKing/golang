// Copyright 2025 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync

import (
	"sync"

	"github.com/searKing/golang/go/pragma"
)

// Pool is generic format of Go sync.Pool, it is safe for concurrent use
// by multiple goroutines without additional locking or coordination.
// Loads, stores, and deletes run in amortized constant time.
type Pool[E any] struct {
	_ pragma.DoNotCopy
	p sync.Pool

	// New optionally specifies a function to generate
	// a value when Get would otherwise return zero.
	// It may not be changed concurrently with calls to Get.
	New func() E
}

// Put adds x to the pool.
func (p *Pool[E]) Put(x E) {
	p.p.Put(x)
}

// Get selects an arbitrary item from the [Pool], removes it from the
// Pool, and returns it to the caller.
// Get may choose to ignore the pool and treat it as empty.
// Callers should not assume any relation between values passed to [Pool.Put] and
// the values returned by Get.
//
// If Get would otherwise return zero and p.New is non-nil, Get returns
// the result of calling p.New.
func (p *Pool[E]) Get() E {
	x := p.p.Get()
	if x == nil {
		if p.New != nil {
			return p.New()
		}
	}
	if x == nil {
		var zeroE E
		return zeroE
	}
	return x.(E)
}
