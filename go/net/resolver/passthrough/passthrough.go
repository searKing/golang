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

func (*passthroughBuilder) Build(ctx context.Context, target resolver.Target, opts ...resolver.BuildOption) (resolver.Resolver, error) {
	var opt resolver.Build
	opt.ApplyOptions(opts...)
	r := &passthroughResolver{
		target: target,
		cc:     opt.ClientConn,
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

func (r *passthroughResolver) ResolveOneAddr(ctx context.Context, opts ...resolver.ResolveOneAddrOption) (resolver.Address, error) {
	return resolver.Address{Addr: r.target.Endpoint}, nil
}

func (r *passthroughResolver) ResolveAddr(ctx context.Context, opts ...resolver.ResolveAddrOption) ([]resolver.Address, error) {
	return []resolver.Address{{Addr: r.target.Endpoint}}, nil
}

func (r *passthroughResolver) ResolveNow(ctx context.Context, opts ...resolver.ResolveNowOption) {}

func (r *passthroughResolver) start() {
	if r.cc != nil {
		_ = r.cc.UpdateState(resolver.State{Addresses: []resolver.Address{{Addr: r.target.Endpoint}}})
	}
}

func (*passthroughResolver) Close() {}
