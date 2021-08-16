// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package url

import (
	"context"
	"net/url"

	"github.com/searKing/golang/go/net/resolver"
)

// ResolveWithTarget reset Host in url.Url by resolver.Target
// ResolveWithTarget always returns a new URL instance,
// even if the returned URL is identical to either the
// base or reference.
func ResolveWithTarget(ctx context.Context, u *url.URL, target string) (*url.URL, error) {
	if u == nil {
		return nil, nil
	}
	u2 := *u
	if target == "" {
		return &u2, nil
	}
	address, err := resolver.ResolveOneAddr(ctx, target)
	if err != nil {
		return nil, err
	}
	u2.Host = address.Addr
	return &u2, nil
}
