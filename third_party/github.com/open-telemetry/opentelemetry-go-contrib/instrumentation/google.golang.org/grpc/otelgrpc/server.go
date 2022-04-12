// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package otelgrpc

import (
	"context"

	errors_ "github.com/searKing/golang/go/errors"
	"go.opentelemetry.io/otel/metric/instrument"
	"google.golang.org/grpc"
)

var (
	// DefaultServerMetrics is the default instance of ServerMetrics. It is
	// intended to be used in conjunction the default Prometheus metrics
	// registry.
	DefaultServerMetrics = NewServerMetrics()

	// UnaryServerMetricInterceptor is a gRPC server-side interceptor that provides Metric monitoring for Unary RPCs.
	UnaryServerMetricInterceptor = DefaultServerMetrics.UnaryServerInterceptor()

	// StreamServerMetricInterceptor is a gRPC server-side interceptor that provides Metric monitoring for Streaming RPCs.
	StreamServerMetricInterceptor = DefaultServerMetrics.StreamServerInterceptor()
)

// ResetServerCounter recreate recording of all counters of RPCs.
// This function acts on the DefaultServerMetrics variable.
func ResetServerCounter(opts ...instrument.Option) {
	DefaultServerMetrics.ResetCounter(opts...)
}

// Register takes a gRPC server and pre-initializes all counters to 0. This
// allows for easier monitoring in Prometheus (no missing metrics), and should
// be called *after* all services have been registered with the server. This
// function acts on the DefaultServerMetrics variable.
func Register(ctx context.Context, server *grpc.Server) {
	DefaultServerMetrics.InitializeMetrics(ctx, server)
}

// EnableServerHandledTimeHistogram turns on recording of handling time of RPCs.
// Histogram metrics can be very expensive for Prometheus to retain and query.
func EnableServerHandledTimeHistogram(opts ...instrument.Option) {
	errors_.Must(DefaultServerMetrics.EnableServerHandledTimeHistogram(opts...))
}

// EnableServerStreamReceiveTimeHistogram turns on recording of single message receive time of streaming RPCs.
// Histogram metrics can be very expensive for Prometheus to retain and query.
func EnableServerStreamReceiveTimeHistogram(opts ...instrument.Option) {
	errors_.Must(DefaultServerMetrics.EnableServerStreamReceiveTimeHistogram(opts...))
}

// EnableServerStreamReceiveSizeHistogram turns on recording of single message receive size of streaming RPCs.
// Histogram metrics can be very expensive for Prometheus to retain and query.
func EnableServerStreamReceiveSizeHistogram(opts ...instrument.Option) {
	errors_.Must(DefaultServerMetrics.EnableServerStreamReceiveSizeHistogram(opts...))
}

// EnableServerStreamSendTimeHistogram turns on recording of single message send time of streaming RPCs.
// Histogram metrics can be very expensive for Prometheus to retain and query.
func EnableServerStreamSendTimeHistogram(opts ...instrument.Option) {
	errors_.Must(DefaultServerMetrics.EnableServerStreamSendTimeHistogram(opts...))
}

// EnableServerStreamSendSizeHistogram turns on recording of single message send size of streaming RPCs.
// Histogram metrics can be very expensive for Prometheus to retain and query.
func EnableServerStreamSendSizeHistogram(opts ...instrument.Option) {
	errors_.Must(DefaultServerMetrics.EnableServerStreamSendSizeHistogram(opts...))
}
