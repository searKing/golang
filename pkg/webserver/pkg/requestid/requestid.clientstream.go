// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package requestid

import (
	"github.com/searKing/golang/third_party/google.golang.org/grpc/interceptors"
	"google.golang.org/grpc"
)

// requestIdClientStream wraps grpc.ClientStream allowing each Sent/Recv of message to get/set request-id.
type requestIdClientStream struct {
	grpc.ClientStream
}

func (s *requestIdClientStream) SendMsg(reply any) error {
	newCtx, _ := tagLoggingRequestId(s.Context(), reply)
	wrapped := interceptors.WrapClientStream(s.ClientStream)
	wrapped.WrappedContext = newCtx
	return wrapped.SendMsg(reply)
}

func (s *requestIdClientStream) RecvMsg(req any) error {
	newCtx, _ := tagLoggingRequestId(s.Context(), req)
	wrapped := interceptors.WrapClientStream(s.ClientStream)
	wrapped.WrappedContext = newCtx
	return wrapped.RecvMsg(req)
}
