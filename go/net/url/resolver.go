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
func ResolveWithTarget(ctx context.Context, u *url.URL, target string) {
	if u == nil {
		return
	}
	address, err := resolver.ResolveOneAddr(ctx, target)
	if err != nil {
		return
	}
	u.Host = address.Addr
	return
}
