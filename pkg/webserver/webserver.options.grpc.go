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
	if f.fc.MaxReceiveMessageSizeInBytes > 0 {
		opts = append(opts, grpc.MaxRecvMsgSize(f.fc.MaxReceiveMessageSizeInBytes))
	} else {
		opts = append(opts, grpc.MaxRecvMsgSize(defaultMaxReceiveMessageSize))
	}
	if f.fc.StatsHandling {
		// log for the related stats handling (e.g., RPCs, connections).
		opts = append(opts, grpc.StatsHandler(&stats.ServerHandler{}))
	}
	if f.fc.EnableOpenTelemetry {
		opts = append(opts, otel.ServerOptions()...)
	}
	return opts
}

func (f *Factory) DialOptions(opts ...grpc.DialOption) []grpc.DialOption {
	if f.fc.NoGrpcProxy {
		opts = append(opts, grpc.WithNoProxy())
	}
	if !f.fc.ForceDisableTls {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	if f.fc.MaxReceiveMessageSizeInBytes > 0 {
		opts = append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(f.fc.MaxReceiveMessageSizeInBytes), grpc.MaxCallSendMsgSize(f.fc.MaxReceiveMessageSizeInBytes)))
	} else {
		opts = append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(defaultMaxReceiveMessageSize), grpc.MaxCallSendMsgSize(defaultMaxSendMessageSize)))
	}
	if f.fc.StatsHandling {
		// log for the related stats handling (e.g., RPCs, connections).
		opts = append(opts, grpc.WithStatsHandler(&stats.ClientHandler{}))
	}
	if f.fc.EnableOpenTelemetry {
		opts = append(opts, otel.DialOptions()...)
	}
	opts = append(opts, grpc.WithChainUnaryInterceptor(f.UnaryClientInterceptors()...))
	opts = append(opts, grpc.WithChainStreamInterceptor(f.StreamClientInterceptors()...))
	return opts
}
