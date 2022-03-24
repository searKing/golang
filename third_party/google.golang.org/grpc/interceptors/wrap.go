// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package interceptors

import (
	"context"

	"google.golang.org/grpc"
)

// WrappedServerStream is a thin wrapper around grpc.ServerStream that allows modifying context.
type WrappedServerStream struct {
	grpc.ServerStream
	// WrappedContext is the wrapper's own Context. You can assign it.
	WrappedContext context.Context
}

// Context returns the wrapper's WrappedContext, overwriting the nested grpc.ServerStream.Context()
func (w *WrappedServerStream) Context() context.Context {
	return w.WrappedContext
}

// WrapServerStream returns a ServerStream that has the ability to overwrite context.
func WrapServerStream(stream grpc.ServerStream) *WrappedServerStream {
	if s, ok := stream.(*WrappedServerStream); ok {
		return s
	}
	return &WrappedServerStream{ServerStream: stream, WrappedContext: stream.Context()}
}
