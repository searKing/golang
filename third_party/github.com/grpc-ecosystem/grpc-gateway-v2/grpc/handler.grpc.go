// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpc

import "google.golang.org/grpc"

type GRPCHandler interface {
	Register(srv *grpc.Server)
}
type GRPCHandlerFunc func(srv *grpc.Server)

func (f GRPCHandlerFunc) Register(srv *grpc.Server) {
	f(srv)
}
