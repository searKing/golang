// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package unix implements a resolver for unix targets.
package unix

import (
	"context"
	"fmt"

	"github.com/searKing/golang/go/net/resolver"
)

const unixScheme = "unix"
const unixAbstractScheme = "unix-abstract"

type builder struct {
	scheme string
}

func (b *builder) Build(target resolver.Target, cc resolver.ClientConn, opts ...resolver.BuildOption) (resolver.Resolver, error) {
	if target.Authority != "" {
		return nil, fmt.Errorf("invalid (non-empty) authority: %v", target.Authority)
	}
	addr := target.Endpoint
	if b.scheme == unixAbstractScheme {
		// prepend "\x00" to address for unix-abstract
		addr = "\x00" + addr
	}
	_ = cc.UpdateState(resolver.State{Addresses: []resolver.Address{{Addr: addr}}})
	return &nopResolver{Addresses: []resolver.Address{{Addr: addr}}}, nil
}

func (b *builder) Scheme() string {
	return b.scheme
}

type nopResolver struct {
	// Addresses is the latest set of resolved addresses for the target.
	Addresses []resolver.Address
}

func (r *nopResolver) ResolveAddr(ctx context.Context, opts ...resolver.ResolveAddrOption) ([]resolver.Address, error) {
	return r.Addresses, nil
}

func (*nopResolver) ResolveNow(ctx context.Context, opts ...resolver.ResolveNowOption) {}

func (*nopResolver) Close() {}

func init() {
	resolver.Register(&builder{scheme: unixScheme})
	resolver.Register(&builder{scheme: unixAbstractScheme})
}
