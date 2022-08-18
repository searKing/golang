// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"fmt"

	sync_ "github.com/searKing/golang/go/sync"
)

var resolverPool = sync_.LruPool{
	New: func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		target, ok := req.(string)
		if !ok {
			return nil, err
		}
		return NewResolver(ctx, target)
	},
	DisableKeepAlives:         false,
	MaxIdleResources:          sync_.DefaultLruPool.MaxIdleResources,
	MaxIdleResourcesPerBucket: 1,
	MaxResourcesPerBucket:     sync_.DefaultLruPool.MaxResourcesPerBucket,
	IdleResourceTimeout:       sync_.DefaultLruPool.IdleResourceTimeout,
}

func ResolveOneAddr(ctx context.Context, target string, opts ...ResolveOneAddrOption) (Address, error) {
	_, err := GetBuilderOrDefault(target)
	if err != nil {
		return Address{}, err
	}
	resolver, put, err := resolverPool.GetOrError(ctx, target)
	if err != nil {
		return Address{}, err
	}
	defer put()
	if resolver, ok := resolver.(Resolver); ok {
		return resolver.ResolveOneAddr(ctx, opts...)
	}
	return Address{}, fmt.Errorf("could not get resolver for target: %q", target)
}

func ResolveAddr(ctx context.Context, target string, opts ...ResolveAddrOption) ([]Address, error) {
	_, err := GetBuilderOrDefault(target)
	if err != nil {
		return nil, err
	}
	resolver, put, err := resolverPool.GetOrError(ctx, target)
	if err != nil {
		return nil, err
	}
	defer put()
	if resolver, ok := resolver.(Resolver); ok {
		return resolver.ResolveAddr(ctx, opts...)
	}
	return nil, fmt.Errorf("could not get resolver for target: %q", target)
}

func ResolveNow(ctx context.Context, target string, opts ...ResolveNowOption) error {
	_, err := GetBuilderOrDefault(target)
	if err != nil {
		return err
	}

	resolver, put, err := resolverPool.GetOrError(ctx, target)
	if err != nil {
		return err
	}
	defer put()
	if resolver, ok := resolver.(Resolver); ok {
		resolver.ResolveNow(ctx, opts...)
		return nil
	}
	return fmt.Errorf("could not get resolver for target: %q", target)
}

func GetBuilderOrDefault(target string) (Builder, error) {
	tgt := ParseTarget(target)
	builder := Get(tgt.Scheme)
	if builder == nil {
		// If resolver builder is still nil, the parsed target's scheme is
		// not registered. Fallback to default resolver and set Endpoint to
		// the original target.
		tgt.Scheme = GetDefaultScheme()
		builder = Get(tgt.Scheme)
		if builder == nil {
			return nil, fmt.Errorf("could not get resolver for default scheme: %q", tgt.Scheme)
		}
	}
	return builder, nil
}

func NewResolver(ctx context.Context, target string, opts ...BuildOption) (Resolver, error) {
	tgt := ParseTarget(target)
	builder := Get(tgt.Scheme)
	if builder == nil {
		// If resolver builder is still nil, the parsed target's scheme is
		// not registered. Fallback to default resolver and set Endpoint to
		// the original target.
		tgt.Scheme = GetDefaultScheme()
		builder = Get(tgt.Scheme)
		if builder == nil {
			return nil, fmt.Errorf("could not get resolver for default scheme: %q", tgt.Scheme)
		}
	}
	return builder.Build(ctx, tgt, opts...)
}
