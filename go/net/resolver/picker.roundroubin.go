// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"sync"
)

type rrPicker struct {
	mu sync.Mutex
	// Start at a random index, as the same RR balancer rebuilds a new
	// picker when SubConn states change, and we don't want to apply excess
	// load to the first server in the list.
	next int
}

func NewRoundRobinPicker(next int) *rrPicker {
	return &rrPicker{
		next: next,
	}
}

func (p *rrPicker) Pick(ctx context.Context, addrs []Address, opts ...PickOption) (Address, error) {
	if len(addrs) == 0 {
		return Address{}, ErrNoAddrAvailable
	}
	p.mu.Lock()
	addr := addrs[p.next]

	p.next = (p.next + 1) % len(addrs)
	p.mu.Unlock()
	return addr, nil
}
