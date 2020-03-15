// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package generator

import "context"

// runtimeGenerator is an implement of Generator's behavior actually.
type runtimeGenerator struct {
	// fired func, as callback when supplierC is consumed successfully
	// arg for msg receiver
	// msg for msg to be delivered
	f   func(ctx context.Context, arg interface{}, msg interface{})
	arg interface{}

	ctx    context.Context
	cancel context.CancelFunc

	// data src
	supplierC <-chan interface{}
	// data dst
	consumerC chan interface{}

	// guard channels below
}

func (g *runtimeGenerator) start() {
	go func() {
		for {
			select {
			case <-g.ctx.Done():
				return
			case s, ok := <-g.supplierC:
				if !ok {
					close(g.consumerC)
					g.stop()
					return
				}
				g.f(g.ctx, g.arg, s)
			}
		}
	}()
}

func (g *runtimeGenerator) stop() bool {
	select {
	case <-g.ctx.Done():
		return false
	default:
		g.cancel()
		return true
	}
}

func (g *runtimeGenerator) stopped() bool {
	select {
	case <-g.ctx.Done():
		return true
	default:
		return false
	}
}
