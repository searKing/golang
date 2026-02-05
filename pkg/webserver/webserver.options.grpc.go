// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webserver

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/searKing/golang/pkg/webserver/pkg/otel"
	"github.com/searKing/golang/pkg/webserver/pkg/stats"
)

func (f *Factory) ServerOptions(opts ...grpc.ServerOption) []grpc.ServerOption {
	var s []grpc.ServerOption
	if f.fc.MaxReceiveMessageSizeInBytes > 0 {
		s = append(s, grpc.MaxRecvMsgSize(f.fc.MaxReceiveMessageSizeInBytes))
	} else {
		s = append(s, grpc.MaxRecvMsgSize(defaultMaxReceiveMessageSize))
	}
	if f.fc.StatsHandling {
		// log for the related stats handling (e.g., RPCs, connections).
		s = append(s, grpc.StatsHandler(&stats.ServerHandler{}))
	}
	if f.fc.OtelHandling {
		s = append(s, otel.ServerOptions(f.fc.OtelGrpcOptions...)...)
	}
	return append(s, opts...)
}

func (f *Factory) DialOptions(opts ...grpc.DialOption) []grpc.DialOption {
	var s []grpc.DialOption
	if f.fc.NoGrpcProxy {
		s = append(s, grpc.WithNoProxy())
	}
	if !f.fc.ForceDisableTls {
		s = append(s, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	if f.fc.MaxReceiveMessageSizeInBytes > 0 {
		s = append(s, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(f.fc.MaxReceiveMessageSizeInBytes), grpc.MaxCallSendMsgSize(f.fc.MaxReceiveMessageSizeInBytes)))
	} else {
		s = append(s, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(defaultMaxReceiveMessageSize), grpc.MaxCallSendMsgSize(defaultMaxSendMessageSize)))
	}
	if f.fc.StatsHandling {
		// log for the related stats handling (e.g., RPCs, connections).
		s = append(s, grpc.WithStatsHandler(&stats.ClientHandler{}))
	}
	if f.fc.OtelHandling {
		s = append(s, otel.DialOptions(f.fc.OtelGrpcOptions...)...)
	}
	s = append(s, grpc.WithChainUnaryInterceptor(f.UnaryClientInterceptors()...))
	s = append(s, grpc.WithChainStreamInterceptor(f.StreamClientInterceptors()...))
	return append(s, opts...)
}
