// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package otelgrpc

import (
	errors_ "github.com/searKing/golang/go/errors"
	"go.opentelemetry.io/otel/metric/instrument"
)

var (
	// DefaultClientMetrics is the default instance of ClientMetrics. It is
	// intended to be used in conjunction the default Prometheus metrics
	// registry.
	DefaultClientMetrics = NewClientMetrics()

	// UnaryClientMetricInterceptor is a gRPC client-side interceptor that provides Metric monitoring for Unary RPCs.
	UnaryClientMetricInterceptor = DefaultClientMetrics.UnaryClientInterceptor()

	// StreamClientMetricInterceptor is a gRPC client-side interceptor that provides Metric monitoring for Streaming RPCs.
	StreamClientMetricInterceptor = DefaultClientMetrics.StreamClientInterceptor()
)

// ResetClientCounter recreate recording of all counters of RPCs.
// This function acts on the DefaultClientMetrics variable.
func ResetClientCounter(opts ...instrument.Option) {
	DefaultClientMetrics.ResetCounter(opts...)
}

// EnableClientHandledTimeHistogram turns on recording of handling time of
// RPCs. Histogram metrics can be very expensive for Prometheus to retain and
// query. This function acts on the DefaultClientMetrics variable.
func EnableClientHandledTimeHistogram(opts ...instrument.Option) {
	errors_.Must(DefaultClientMetrics.EnableClientHandledTimeHistogram(opts...))
}

// EnableClientStreamReceiveTimeHistogram turns on recording of
// single message receive time of streaming RPCs.
// This function acts on the DefaultClientMetrics variable.
func EnableClientStreamReceiveTimeHistogram(opts ...instrument.Option) {
	errors_.Must(DefaultClientMetrics.EnableClientStreamReceiveTimeHistogram(opts...))
}

// EnableClientStreamReceiveSizeHistogram turns on recording of
// single message receive size of streaming RPCs.
// This function acts on the DefaultClientMetrics variable
func EnableClientStreamReceiveSizeHistogram(opts ...instrument.Option) {
	errors_.Must(DefaultClientMetrics.EnableClientStreamReceiveSizeHistogram(opts...))
}

// EnableClientStreamSendTimeHistogram turns on recording of
// single message send time of streaming RPCs.
// This function acts on the DefaultClientMetrics variable.
func EnableClientStreamSendTimeHistogram(opts ...instrument.Option) {
	errors_.Must(DefaultClientMetrics.EnableClientStreamSendTimeHistogram(opts...))
}

// EnableClientStreamSendSizeHistogram turns on recording of
// single message receive size of streaming RPCs.
// This function acts on the DefaultClientMetrics variable
func EnableClientStreamSendSizeHistogram(opts ...instrument.Option) {
	errors_.Must(DefaultClientMetrics.EnableClientStreamSendSizeHistogram(opts...))
}
