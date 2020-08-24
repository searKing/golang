// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"context"
	"net/http"

	"github.com/golang/protobuf/proto"
	structpb "github.com/golang/protobuf/ptypes/struct"
	struct_ "github.com/searKing/golang/third_party/github.com/golang/protobuf/ptypes/struct"
	"google.golang.org/grpc/grpclog"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type StatusHandler struct {
	Handler func(s *status.Status, errStructpb *structpb.Struct) proto.Message
}

func (h *StatusHandler) HTTPError(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	s, ok := status.FromError(err)
	if !ok {
		if err == nil {
			s = status.New(codes.OK, codes.OK.String())
		} else {
			s = status.New(codes.Unknown, err.Error())
		}
	}

	if err := s.Err(); err != nil {
		grpclog.Errorf("Failed to handle http request %s %s: %v", r.Method, r.RequestURI, err)
	}

	w.Header().Del("Trailer")
	cause := errors.Cause(err)
	var errStructpb *structpb.Struct
	if cause != err {
		errStructpb, _ = struct_.ToProtoStruct(cause.Error())
	}
	var body proto.Message
	if h.Handler != nil {
		body = h.Handler(s, errStructpb)
	} else {
		body = s.Proto()
	}

	// 200 OK forever
	runtime.ForwardResponseMessage(ctx, mux, marshaler, w, r, body)
}
