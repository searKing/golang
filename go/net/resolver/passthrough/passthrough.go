// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package passthrough

import (
	"context"

	"github.com/searKing/golang/go/net/resolver"
)

const scheme = "passthrough"

func init() {
	resolver.Register(&passthroughBuilder{})
}

type passthroughBuilder struct{}

func (*passthroughBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts ...resolver.BuildOption) (resolver.Resolver, error) {
	r := &passthroughResolver{
		target: target,
		cc:     cc,
	}
	r.start()
	return r, nil
}

func (*passthroughBuilder) Scheme() string {
	return scheme
}

type passthroughResolver struct {
	target resolver.Target
	cc     resolver.ClientConn
}

func (r *passthroughResolver) ResolveAddr(ctx context.Context, opts ...resolver.ResolveAddrOption) ([]resolver.Address, error) {
	return []resolver.Address{{Addr: r.target.Endpoint}}, nil
}

func (r *passthroughResolver) ResolveNow(ctx context.Context, opts ...resolver.ResolveNowOption) {}

func (r *passthroughResolver) start() {
	_ = r.cc.UpdateState(resolver.State{Addresses: []resolver.Address{{Addr: r.target.Endpoint}}})
}

func (*passthroughResolver) Close() {}
