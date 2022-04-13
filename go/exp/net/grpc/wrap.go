// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpc

import "context"

func Wrap[REQ any, RESP any](next func(ctx context.Context, req REQ) (RESP, error),
	wrappers ...func(handler func(ctx context.Context, req REQ) (RESP, error)) func(ctx context.Context, req REQ) (RESP, error)) func(
	ctx context.Context, req REQ) (RESP, error) {
	for i := range wrappers {
		w := wrappers[len(wrappers)-1-i]
		if w != nil {
			next = w(next)
		}
	}
	return next
}
