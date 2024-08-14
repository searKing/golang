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

// FieldsFromContextWithForward fill "X-Forwarded-For" and "X-Forwarded-Host" to record http callers
func FieldsFromContextWithForward(ctx context.Context) logging.Fields {
	const xForwardedFor = "X-Forwarded-For"
	const xForwardedHost = "X-Forwarded-Host"

	fields := logging.Fields{}
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		for _, key := range []string{strings.ToLower(xForwardedFor), strings.ToLower(xForwardedHost)} {
			if fwd := md.Get(key); len(fwd) > 0 {
				fields = fields.AppendUnique(logging.Fields{"grpc." + strings.ToLower(key), fwd})
			}
		}
	}
	return fields
}
