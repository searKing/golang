// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package otelgrpc

import (
	"go.opentelemetry.io/otel/metric"
)

var (
	// DefaultClientMetrics is the default instance of ClientMetrics. It is
	// intended to be used in conjunction the default Prometheus metrics
	// registry.
	DefaultClientMetrics = NewClientMetrics()

	// UnaryClientMetricInterceptor is a gRPC client-side interceptor that provides Metric monitoring for Unary RPCs.
	UnaryClientMetricInterceptor = DefaultClientMetrics.UnaryClientInterceptor()
)

// ResetClientCounter recreate recording of all counters of RPCs.
// This function acts on the DefaultClientMetrics variable.
func ResetClientCounter(opts ...metric.InstrumentOption) error {
	return DefaultClientMetrics.ResetCounter(opts...)
}

// EnableClientHandledTimeHistogram turns on recording of handling time of
// RPCs. Histogram metrics can be very expensive for Prometheus to retain and
// query. This function acts on the DefaultClientMetrics variable.
func EnableClientHandledTimeHistogram(opts ...metric.Float64HistogramOption) {
	DefaultClientMetrics.EnableClientHandledTimeHistogram(opts...)
}

// EnableClientStreamReceiveSizeHistogram turns on recording of
// single message receive size of streaming RPCs.
// This function acts on the DefaultClientMetrics variable
func EnableClientStreamReceiveSizeHistogram(opts ...metric.Int64HistogramOption) {
	DefaultClientMetrics.EnableClientStreamReceiveSizeHistogram(opts...)
}

// EnableClientStreamSendSizeHistogram turns on recording of
// single message receive size of streaming RPCs.
// This function acts on the DefaultClientMetrics variable
func EnableClientStreamSendSizeHistogram(opts ...metric.Int64HistogramOption) {
	DefaultClientMetrics.EnableClientStreamSendSizeHistogram(opts...)
}
