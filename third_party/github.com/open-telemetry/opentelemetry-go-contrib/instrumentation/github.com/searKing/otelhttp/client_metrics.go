// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package otelgrpc

import (
	"net"
	"net/http"

	net_ "github.com/searKing/golang/go/net"
	http_ "github.com/searKing/golang/go/net/http"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/unit"
)

// ClientMetrics represents a collection of metrics to be registered on a
// Prometheus metrics registry for a gRPC client.
type ClientMetrics struct {
	ClientHostport string

	clientStartedCounter        metric.Int64Counter
	clientHandledCounter        metric.Int64Counter
	clientStreamRequestReceived metric.Int64Counter
	clientStreamRequestSent     metric.Int64Counter

	// "grpc_type", "grpc_service", "grpc_method"
	clientHandledTimeHistogramEnabled bool
	clientHandledTimeHistogram        metric.Float64Histogram

	clientStreamReceiveSizeHistogramEnabled bool
	clientStreamReceiveSizeHistogram        metric.Int64Histogram

	clientStreamSendSizeHistogramEnabled bool
	clientStreamSendSizeHistogram        metric.Int64Histogram
}

func Meter() metric.Meter {
	return global.Meter(InstrumentationName, metric.WithInstrumentationVersion(InstrumentationVersion))
}

// NewClientMetrics returns a ClientMetrics object. Use a new instance of
// ClientMetrics when not using the default Prometheus metrics registry, for
// example when wanting to control which metrics are added to a registry as
// opposed to automatically adding metrics via init functions.
func NewClientMetrics(opts ...metric.InstrumentOption) *ClientMetrics {
	m := &ClientMetrics{}
	m.ResetCounter(opts...)
	return m
}

// ResetCounter recreate recording of all counters of RPCs.
func (m *ClientMetrics) ResetCounter(opts ...metric.InstrumentOption) {
	var addr string
	if ip, err := net_.ListenIP(); err == nil {
		addr = ip.String()
	}
	m.ClientHostport = net.JoinHostPort(addr, "0")
	// "grpc_type", "grpc_service", "grpc_method"
	m.clientStartedCounter = metric.Must(Meter()).NewInt64Counter(
		"grpc_client_started_total",
		func() []metric.InstrumentOption {
			var options []metric.InstrumentOption
			options = append(options, metric.WithDescription("Total number of RPCs started on the client."))
			options = append(options, opts...)
			return options
		}()...)
	// "grpc_type", "grpc_service", "grpc_method", "grpc_code"
	m.clientHandledCounter = metric.Must(Meter()).NewInt64Counter(
		"grpc_client_handled_total",
		func() []metric.InstrumentOption {
			var options []metric.InstrumentOption
			options = append(options, metric.WithDescription("Total number of RPCs completed by the client, regardless of success or failure."))
			options = append(options, opts...)
			return options
		}()...)
	// "grpc_type", "grpc_service", "grpc_method"
	m.clientStreamRequestReceived = metric.Must(Meter()).NewInt64Counter(
		"grpc_client_msg_received_total",
		func() []metric.InstrumentOption {
			var options []metric.InstrumentOption
			options = append(options, metric.WithDescription("Total number of RPC stream messages received by the client."))
			options = append(options, opts...)
			return options
		}()...)
	// "grpc_type", "grpc_service", "grpc_method"
	m.clientStreamRequestSent = metric.Must(Meter()).NewInt64Counter(
		"grpc_client_msg_sent_total",
		func() []metric.InstrumentOption {
			var options []metric.InstrumentOption
			options = append(options, metric.WithDescription("Total number of gRPC stream messages sent by the client."))
			options = append(options, opts...)
			return options
		}()...)
}

// EnableClientHandledTimeHistogram turns on recording of handling time of RPCs.
// Histogram metrics can be very expensive for Prometheus to retain and query.
func (m *ClientMetrics) EnableClientHandledTimeHistogram(opts ...metric.InstrumentOption) {
	var options []metric.InstrumentOption
	options = append(options,
		metric.WithDescription("Histogram of response latency (seconds) of the gRPC until it is finished by the application."),
		metric.WithUnit("s"))
	options = append(options, opts...)
	if !m.clientHandledTimeHistogramEnabled {
		// https://github.com/open-telemetry/opentelemetry-go/issues/1280
		m.clientHandledTimeHistogram = metric.Must(Meter()).NewFloat64Histogram("grpc_client_handling_seconds", options...)
	}
	m.clientHandledTimeHistogramEnabled = true
}

// EnableClientStreamReceiveSizeHistogram turns on recording of single message receive size of streaming RPCs.
// Histogram metrics can be very expensive for Prometheus to retain and query.
func (m *ClientMetrics) EnableClientStreamReceiveSizeHistogram(opts ...metric.InstrumentOption) {
	var options []metric.InstrumentOption
	options = append(options,
		metric.WithDescription("Histogram of message size (bytes) of the gRPC single message receive."),
		metric.WithUnit(unit.Bytes))
	options = append(options, opts...)
	if !m.clientStreamReceiveSizeHistogramEnabled {
		// https://github.com/open-telemetry/opentelemetry-go/issues/1280
		m.clientStreamReceiveSizeHistogram = metric.Must(Meter()).NewInt64Histogram("grpc_client_msg_recv_handling_bytes", options...)
	}
	m.clientStreamReceiveSizeHistogramEnabled = true
}

// EnableClientStreamSendSizeHistogram turns on recording of single message send size of streaming RPCs.
// Histogram metrics can be very expensive for Prometheus to retain and query.
func (m *ClientMetrics) EnableClientStreamSendSizeHistogram(opts ...metric.InstrumentOption) {
	var options []metric.InstrumentOption
	options = append(options,
		metric.WithDescription("Histogram of message size (bytes) of the gRPC single message send."),
		metric.WithUnit(unit.Bytes))
	options = append(options, opts...)
	if !m.clientStreamSendSizeHistogramEnabled {
		// https://github.com/open-telemetry/opentelemetry-go/issues/1280
		m.clientStreamSendSizeHistogram = metric.Must(Meter()).NewInt64Histogram("grpc_client_msg_send_handling_bytes", options...)
	}
	m.clientStreamSendSizeHistogramEnabled = true
}

// UnaryClientInterceptor is a gRPC client-side interceptor that provides Prometheus monitoring for Unary RPCs.
func (m *ClientMetrics) UnaryClientInterceptor() func(req *http.Request, retry int, invoker http_.ClientInvoker, opts ...http_.DoWithBackoffOption) (resp *http.Response, err error) {
	return func(req *http.Request, retry int, invoker http_.ClientInvoker, opts ...http_.DoWithBackoffOption) (resp *http.Response, err error) {
		ctx := req.Context()
		monitor := newClientReporter(ctx, m, req)
		monitor.SentRequest(ctx, req)
		resp, err = invoker(req, retry)
		if err == nil {
			monitor.ReceivedResponse(ctx, resp)
			monitor.Handled(ctx, resp)
		} else {
			monitor.Handled(ctx, resp)
		}
		return resp, err
	}
}
