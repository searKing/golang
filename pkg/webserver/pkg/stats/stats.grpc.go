// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stats

import (
	"context"
	"fmt"
	"log/slog"

	"google.golang.org/grpc/stats"

	"github.com/searKing/golang/pkg/webserver/pkg/logging"
)

// statsHandleRPC processes the RPC events for conn infos.
func statsHandleRPC(ctx context.Context, s stats.RPCStats) {
	switch st := s.(type) {
	case *stats.InHeader:
		logger := slog.With(logging.Attrs[any](ctx)...)
		if st.Client {
			logger.Info(fmt.Sprintf("gRPC Client Got Conn in for %s from %s to %s",
				st.FullMethod, st.RemoteAddr, st.LocalAddr))
		} else {
			logger.Info(fmt.Sprintf("gRPC Server Got Conn in for %s from %s to %s",
				st.FullMethod, st.LocalAddr, st.RemoteAddr))
		}
	case *stats.OutHeader:
		logger := slog.With(logging.Attrs[any](ctx)...)
		if st.Client {
			logger.Info(fmt.Sprintf("gRPC Client Got Conn out for %s from %s to %s",
				st.FullMethod, st.LocalAddr, st.RemoteAddr))
		} else {
			logger.Info(fmt.Sprintf("gRPC Server Got Conn out for %s from %s to %s",
				st.FullMethod, st.RemoteAddr, st.LocalAddr))
		}
	case *stats.Begin, *stats.End, *stats.InTrailer, *stats.OutTrailer, *stats.OutPayload, *stats.InPayload, *stats.PickerUpdated:
		// do nothing for client
	default:
		logger := slog.With(logging.Attrs[any](ctx)...)
		logger.Info(fmt.Sprintf("unexpected grpc stats: %T", st))
	}
}
