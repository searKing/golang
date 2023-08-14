// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpc

import (
	"context"
	"strings"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc/metadata"
)

func FieldsFromContextWithForward(ctx context.Context) logging.Fields {
	const xForwardedFor = "X-Forwarded-For"
	const xForwardedHost = "X-Forwarded-Host"

	var fields logging.Fields
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		for _, key := range []string{strings.ToLower(xForwardedFor), strings.ToLower(xForwardedHost)} {
			fwd := md.Get(key)
			if len(fwd) > 0 {
				fields = fields.AppendUnique(logging.Fields{"grpc." + strings.ToLower(key), fwd})
			}
		}
	}
	return fields
}
