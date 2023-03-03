// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package otelgrpc

import (
	"context"
	"io"
	"net"
	"time"

	errors_ "github.com/searKing/golang/go/errors"
	slices_ "github.com/searKing/golang/go/exp/slices"
	net_ "github.com/searKing/golang/go/net"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// ClientMetrics represents a collection of metrics to be registered on a
// Prometheus metrics registry for a gRPC client.
type ClientMetrics struct {
	ClientHostport string

	clientStartedCounter    instrument.Int64Counter
	clientHandledCounter    instrument.Int64Counter
	clientStreamMsgReceived instrument.Int64Counter
	clientStreamMsgSent     instrument.Int64Counter

	// "grpc_type", "grpc_service", "grpc_method"
	clientHandledTimeHistogramEnabled bool
	clientHandledTimeHistogram        instrument.Float64Histogram

	clientStreamReceiveTimeHistogramEnabled bool
	clientStreamReceiveTimeHistogram        instrument.Float64Histogram

	clientStreamReceiveSizeHistogramEnabled bool
	clientStreamReceiveSizeHistogram        instrument.Int64Histogram

	clientStreamSendTimeHistogramEnabled bool
	clientStreamSendTimeHistogram        instrument.Float64Histogram

	clientStreamSendSizeHistogramEnabled bool
	clientStreamSendSizeHistogram        instrument.Int64Histogram
}

func Meter() metric.Meter {
	return global.MeterProvider().Meter(InstrumentationName, metric.WithInstrumentationVersion(InstrumentationVersion))
}

// NewClientMetrics returns a ClientMetrics object. Use a new instance of
// ClientMetrics when not using the default Prometheus metrics registry, for
// example when wanting to control which metrics are added to a registry as
// opposed to automatically adding metrics via init functions.
func NewClientMetrics(opts ...instrument.Option) *ClientMetrics {
	m := &ClientMetrics{}
	errors_.Must(m.ResetCounter(opts...))
	return m
}

// ResetCounter recreate recording of all counters of RPCs.
func (m *ClientMetrics) ResetCounter(opts ...instrument.Option) (err error) {
	int64Opts := slices_.TypeAssertFilterFunc[[]instrument.Option, []instrument.Int64Option](opts,
		func(opt instrument.Option) (instrument.Int64Option, bool) {
			o, ok := opt.(instrument.Int64Option)
			return o, ok
		})
	var addr string
	if ip, err := net_.ListenIP(); err == nil {
		addr = ip.String()
	}
	m.ClientHostport = net.JoinHostPort(addr, "0")
	// "grpc_type", "grpc_service", "grpc_method"
	m.clientStartedCounter, err = Meter().Int64Counter(
		"grpc_client_started_total",
		func() []instrument.Int64Option {
			var options []instrument.Int64Option
			options = append(options, instrument.WithDescription("Total number of RPCs started on the client."))
			options = append(options, int64Opts...)
			return options
		}()...)
	if err != nil {
		return err
	}
	// "grpc_type", "grpc_service", "grpc_method", "grpc_code"
	m.clientHandledCounter, err = Meter().Int64Counter(
		"grpc_client_handled_total",
		func() []instrument.Int64Option {
			var options []instrument.Int64Option
			options = append(options, instrument.WithDescription("Total number of RPCs completed by the client, regardless of success or failure."))
			options = append(options, int64Opts...)
			return options
		}()...)
	if err != nil {
		return err
	}
	// "grpc_type", "grpc_service", "grpc_method"
	m.clientStreamMsgReceived, err = Meter().Int64Counter(
		"grpc_client_msg_received_total",
		func() []instrument.Int64Option {
			var options []instrument.Int64Option
			options = append(options, instrument.WithDescription("Total number of RPC stream messages received by the client."))
			options = append(options, int64Opts...)
			return options
		}()...)
	if err != nil {
		return err
	}
	// "grpc_type", "grpc_service", "grpc_method"
	m.clientStreamMsgSent, err = Meter().Int64Counter(
		"grpc_client_msg_sent_total",
		func() []instrument.Int64Option {
			var options []instrument.Int64Option
			options = append(options, instrument.WithDescription("Total number of gRPC stream messages sent by the client."))
			options = append(options, int64Opts...)
			return options
		}()...)
	return err
}

// EnableClientHandledTimeHistogram turns on recording of handling time of RPCs.
// Histogram metrics can be very expensive for Prometheus to retain and query.
func (m *ClientMetrics) EnableClientHandledTimeHistogram(opts ...instrument.Float64Option) (err error) {
	var options []instrument.Float64Option
	options = append(options,
		instrument.WithDescription("Histogram of response latency (seconds) of the gRPC until it is finished by the application."),
		instrument.WithUnit("s"))
	options = append(options, opts...)
	if !m.clientHandledTimeHistogramEnabled {
		// https://github.com/open-telemetry/opentelemetry-go/issues/1280
		m.clientHandledTimeHistogram, err = Meter().Float64Histogram("grpc_client_handling_seconds", options...)
		if err != nil {
			return err
		}
	}
	m.clientHandledTimeHistogramEnabled = true
	return nil
}

// EnableClientStreamReceiveTimeHistogram turns on recording of single message receive time of streaming RPCs.
// Histogram metrics can be very expensive for Prometheus to retain and query.
func (m *ClientMetrics) EnableClientStreamReceiveTimeHistogram(opts ...instrument.Float64Option) (err error) {
	var options []instrument.Float64Option
	options = append(options,
		instrument.WithDescription("Histogram of response latency (seconds) of the gRPC single message receive."),
		instrument.WithUnit("s"))
	options = append(options, opts...)
	if !m.clientStreamReceiveTimeHistogramEnabled {
		// https://github.com/open-telemetry/opentelemetry-go/issues/1280
		m.clientStreamReceiveTimeHistogram, err = Meter().Float64Histogram("grpc_client_msg_recv_handling_seconds", options...)
		if err != nil {
			return err
		}
	}
	m.clientStreamReceiveTimeHistogramEnabled = true
	return nil
}

// EnableClientStreamReceiveSizeHistogram turns on recording of single message receive size of streaming RPCs.
// Histogram metrics can be very expensive for Prometheus to retain and query.
func (m *ClientMetrics) EnableClientStreamReceiveSizeHistogram(opts ...instrument.Int64Option) (err error) {
	var options []instrument.Int64Option
	options = append(options,
		instrument.WithDescription("Histogram of message size (bytes) of the gRPC single message receive."),
		instrument.WithUnit(uBytes))
	options = append(options, opts...)
	if !m.clientStreamReceiveSizeHistogramEnabled {
		// https://github.com/open-telemetry/opentelemetry-go/issues/1280
		m.clientStreamReceiveSizeHistogram, err = Meter().Int64Histogram("grpc_client_msg_recv_handling_bytes", options...)
		if err != nil {
			return err
		}
	}
	m.clientStreamReceiveSizeHistogramEnabled = true
	return nil
}

// EnableClientStreamSendTimeHistogram turns on recording of single message send time of streaming RPCs.
// Histogram metrics can be very expensive for Prometheus to retain and query.
func (m *ClientMetrics) EnableClientStreamSendTimeHistogram(opts ...instrument.Float64Option) (err error) {
	var options []instrument.Float64Option
	options = append(options,
		instrument.WithDescription("Histogram of response latency (seconds) of the gRPC single message send."),
		instrument.WithUnit("s"))
	options = append(options, opts...)
	if !m.clientStreamSendTimeHistogramEnabled {
		// https://github.com/open-telemetry/opentelemetry-go/issues/1280
		m.clientStreamSendTimeHistogram, err = Meter().Float64Histogram("grpc_client_msg_send_handling_seconds", options...)
		if err != nil {
			return err
		}
	}
	m.clientStreamSendTimeHistogramEnabled = true
	return nil
}

// EnableClientStreamSendSizeHistogram turns on recording of single message send size of streaming RPCs.
// Histogram metrics can be very expensive for Prometheus to retain and query.
func (m *ClientMetrics) EnableClientStreamSendSizeHistogram(opts ...instrument.Int64Option) (err error) {
	var options []instrument.Int64Option
	options = append(options,
		instrument.WithDescription("Histogram of message size (bytes) of the gRPC single message send."),
		instrument.WithUnit(uBytes))
	options = append(options, opts...)
	if !m.clientStreamSendSizeHistogramEnabled {
		// https://github.com/open-telemetry/opentelemetry-go/issues/1280
		m.clientStreamSendSizeHistogram, err = Meter().Int64Histogram("grpc_client_msg_send_handling_bytes", options...)
		if err != nil {
			return err
		}
	}
	m.clientStreamSendSizeHistogramEnabled = true
	return nil
}

// UnaryClientInterceptor is a gRPC client-side interceptor that provides Prometheus monitoring for Unary RPCs.
func (m *ClientMetrics) UnaryClientInterceptor() func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		monitor := newClientReporter(ctx, m, Unary, method, cc.Target(), m.ClientHostport)
		monitor.SentMessage(ctx, req)
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err == nil {
			monitor.ReceivedMessage(ctx, reply)
		}
		st, _ := FromError(err)
		monitor.Handled(ctx, st.Code())
		return err
	}
}

// StreamClientInterceptor is a gRPC client-side interceptor that provides Prometheus monitoring for Streaming RPCs.
func (m *ClientMetrics) StreamClientInterceptor() func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		monitor := newClientReporter(ctx, m, clientStreamType(desc), method, cc.Target(), m.ClientHostport)
		clientStream, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			st, _ := FromError(err)
			monitor.Handled(ctx, st.Code())
			return nil, err
		}
		return &monitoredClientStream{clientStream, monitor}, nil
	}
}

func clientStreamType(desc *grpc.StreamDesc) grpcType {
	if desc.ClientStreams && !desc.ServerStreams {
		return ClientStream
	} else if !desc.ClientStreams && desc.ServerStreams {
		return ServerStream
	}
	return BidiStream
}

// monitoredClientStream wraps grpc.ClientStream allowing each Sent/Recv of message to increment counters.
type monitoredClientStream struct {
	grpc.ClientStream
	monitor *clientReporter
}

func (s *monitoredClientStream) SendMsg(m interface{}) error {
	now := time.Now()
	err := s.ClientStream.SendMsg(m)
	s.monitor.SendMessageTimer(context.Background(), now)
	if err == nil {
		s.monitor.SentMessage(context.Background(), m)
	}
	return err
}

func (s *monitoredClientStream) RecvMsg(m interface{}) error {
	now := time.Now()
	err := s.ClientStream.RecvMsg(m)
	s.monitor.ReceiveMessageTimer(context.Background(), now)

	if err == nil {
		s.monitor.ReceivedMessage(context.Background(), m)
	} else if err == io.EOF {
		s.monitor.Handled(context.Background(), codes.OK)
	} else {
		st, _ := FromError(err)
		s.monitor.Handled(context.Background(), st.Code())
	}
	return err
}
