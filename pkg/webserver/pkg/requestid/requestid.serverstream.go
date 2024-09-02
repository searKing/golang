// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package requestid

import (
	"github.com/searKing/golang/third_party/google.golang.org/grpc/interceptors"
	"google.golang.org/grpc"
)

// requestIdServerStream wraps grpc.ServerStream allowing each Sent/Recv of message to get/set request-id.
type requestIdServerStream struct {
	grpc.ServerStream
}

func (s *requestIdServerStream) SendMsg(reply any) error {
	newCtx, _ := tagLoggingRequestId(s.Context(), reply)
	wrapped := interceptors.WrapServerStream(s.ServerStream)
	wrapped.WrappedContext = newCtx
	return wrapped.SendMsg(reply)
}

func (s *requestIdServerStream) RecvMsg(req any) error {
	newCtx, _ := tagLoggingRequestId(s.Context(), req)
	wrapped := interceptors.WrapServerStream(s.ServerStream)
	wrapped.WrappedContext = newCtx
	return wrapped.RecvMsg(req)
}
