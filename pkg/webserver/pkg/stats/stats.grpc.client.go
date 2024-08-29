// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stats

import (
	"context"

	"google.golang.org/grpc/stats"
)

var _ stats.Handler = (*ClientHandler)(nil)

// ClientHandler implements a gRPC stats.Handler for recording gRPC stats.
// Use with gRPC clients only.
type ClientHandler struct{}

// TagRPC implements per-RPC context management.
func (c ClientHandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	// no-op
	return ctx
}

// HandleRPC implements per-RPC tracing and stats instrumentation.
func (c ClientHandler) HandleRPC(ctx context.Context, rpcStats stats.RPCStats) {
	statsHandleRPC(ctx, rpcStats)
}

// TagConn exists to satisfy gRPC stats.Handler.
func (c ClientHandler) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	// no-op
	return ctx
}

// HandleConn exists to satisfy gRPC stats.Handler.
func (c ClientHandler) HandleConn(ctx context.Context, connStats stats.ConnStats) {
	// no-op
}
