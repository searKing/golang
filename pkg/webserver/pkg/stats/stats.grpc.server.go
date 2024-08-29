// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stats

import (
	"context"

	"google.golang.org/grpc/stats"
)

var _ stats.Handler = (*ServerHandler)(nil)

// ServerHandler implements a gRPC stats.Handler for recording gRPC stats.
// Use with gRPC servers only.
type ServerHandler struct{}

// TagRPC implements per-RPC context management.
func (s ServerHandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	// no-op
	return ctx
}

// HandleRPC implements per-RPC tracing and stats instrumentation.
func (s ServerHandler) HandleRPC(ctx context.Context, rpcStats stats.RPCStats) {
	statsHandleRPC(ctx, rpcStats)
}

// TagConn exists to satisfy gRPC stats.Handler.
func (s ServerHandler) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	// no-op
	return ctx
}

// HandleConn exists to satisfy gRPC stats.Handler.
func (s ServerHandler) HandleConn(ctx context.Context, connStats stats.ConnStats) {
	return
}
