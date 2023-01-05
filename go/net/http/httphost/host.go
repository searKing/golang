// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httphost

import (
	"context"
	"fmt"
	"net/url"

	"github.com/searKing/golang/go/net/resolver"
)

// unique type to prevent assignment.
type hostContextKey struct{}

// Host specifies TargetUrl and HostTarget to return a dynamic host.
type Host struct {
	// HostTarget is as like gRPC Naming for host service discovery, with Host in TargetUrl replaced if not empty.
	HostTarget string
	// resolve HostTarget to host and replace host if HostTarget resolved
	ReplaceHostInRequest bool

	// HostTargetAddrResolved is the host's addr resolved and picked from resolver.
	HostTargetAddrResolved resolver.Address
}

// ContextHost returns the Host associated with the
// provided context. If none, it returns nil.
func ContextHost(ctx context.Context) *Host {
	host, _ := ctx.Value(hostContextKey{}).(*Host)
	return host
}

// WithHost returns a new context based on the provided parent
// ctx. HTTP client requests made with the returned context will use
// the provided host hooks
func WithHost(ctx context.Context, host *Host) context.Context {
	if host == nil {
		panic("nil host")
	}
	return context.WithValue(ctx, hostContextKey{}, host)
}

func ParseTargetUrl(host string) (*url.URL, error) {
	if host == "" {
		return nil, nil
	}

	hostURL, err := url.Parse(host)
	if err != nil {
		return nil, fmt.Errorf("invalid host address %q: %v", host, err)
	}
	return hostURL, nil
}
