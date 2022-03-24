// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package otelgrpc

import (
	"context"
	"time"

	net_ "github.com/searKing/golang/go/net"
	otelgrpc_ "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/unit"
	"google.golang.org/grpc"
)

// ServerMetrics represents a collection of metrics to be registered on a
// Prometheus metrics registry for a gRPC server.
type ServerMetrics struct {
	ServerHostport string

	serverStartedCounter    metric.Int64Counter
	serverHandledCounter    metric.Int64Counter
	serverStreamMsgReceived metric.Int64Counter
	serverStreamMsgSent     metric.Int64Counter

	serverHandledTimeHistogramEnabled bool
	serverHandledTimeHistogram        metric.Float64Histogram

	serverStreamReceiveTimeHistogramEnabled bool
	serverStreamReceiveTimeHistogram        metric.Float64Histogram

	serverStreamReceiveSizeHistogramEnabled bool
	serverStreamReceiveSizeHistogram        metric.Int64Histogram

	serverStreamSendTimeHistogramEnabled bool
	serverStreamSendTimeHistogram        metric.Float64Histogram

	serverStreamSendSizeHistogramEnabled bool
	serverStreamSendSizeHistogram        metric.Int64Histogram
}

// NewServerMetrics returns a ServerMetrics object. Use a new instance of
// ServerMetrics when not using the default Prometheus metrics registry, for
// example when wanting to control which metrics are added to a registry as
// opposed to automatically adding metrics via init functions.
func NewServerMetrics(opts ...metric.InstrumentOption) *ServerMetrics {
	m := &ServerMetrics{}
	m.ResetCounter(opts...)
	return m
}

// ResetCounter recreate recording of all counters of RPCs.
func (m *ServerMetrics) ResetCounter(opts ...metric.InstrumentOption) {
	var addr string
	if ip, err := net_.ListenIP(); err == nil {
		addr = ip.String()
	}

	m.ServerHostport = addr
	// "grpc_type", "grpc_service", "grpc_method"
	m.serverStartedCounter = metric.Must(Meter()).NewInt64Counter(
		"grpc_server_started_total",
		func() []metric.InstrumentOption {
			var options []metric.InstrumentOption
			options = append(options, metric.WithDescription("Total number of RPCs started on the server."))
			options = append(options, opts...)
			return options
		}()...)
	// "grpc_type", "grpc_service", "grpc_method", "grpc_code"
	m.serverHandledCounter = metric.Must(Meter()).NewInt64Counter(
		"grpc_server_handled_total",
		func() []metric.InstrumentOption {
			var options []metric.InstrumentOption
			options = append(options, metric.WithDescription("Total number of RPCs completed by the server, regardless of success or failure."))
			options = append(options, opts...)
			return options
		}()...)
	// "grpc_type", "grpc_service", "grpc_method"
	m.serverStreamMsgReceived = metric.Must(Meter()).NewInt64Counter(
		"grpc_server_msg_received_total",
		func() []metric.InstrumentOption {
			var options []metric.InstrumentOption
			options = append(options, metric.WithDescription("Total number of RPC stream messages received by the server."))
			options = append(options, opts...)
			return options
		}()...)
	// "grpc_type", "grpc_service", "grpc_method"
	m.serverStreamMsgSent = metric.Must(Meter()).NewInt64Counter(
		"grpc_server_msg_sent_total",
		func() []metric.InstrumentOption {
			var options []metric.InstrumentOption
			options = append(options, metric.WithDescription("Total number of gRPC stream messages sent by the server."))
			options = append(options, opts...)
			return options
		}()...)
}

// EnableServerHandledTimeHistogram turns on recording of handling time of RPCs.
// Histogram metrics can be very expensive for Prometheus to retain and query.
func (m *ServerMetrics) EnableServerHandledTimeHistogram(opts ...metric.InstrumentOption) {
	var options []metric.InstrumentOption
	options = append(options,
		metric.WithDescription("Histogram of response latency (seconds) of the gRPC until it is finished by the server."),
		metric.WithUnit("s"))
	options = append(options, opts...)
	if !m.serverHandledTimeHistogramEnabled {
		// https://github.com/open-telemetry/opentelemetry-go/issues/1280
		m.serverHandledTimeHistogram = metric.Must(Meter()).NewFloat64Histogram("grpc_server_handling_seconds", options...)
	}
	m.serverHandledTimeHistogramEnabled = true
}

// EnableServerStreamReceiveTimeHistogram turns on recording of single message receive time of streaming RPCs.
// Histogram metrics can be very expensive for Prometheus to retain and query.
func (m *ServerMetrics) EnableServerStreamReceiveTimeHistogram(opts ...metric.InstrumentOption) {
	var options []metric.InstrumentOption
	options = append(options,
		metric.WithDescription("Histogram of response latency (seconds) of the gRPC single message receive."),
		metric.WithUnit("s"))
	options = append(options, opts...)
	if !m.serverStreamReceiveTimeHistogramEnabled {
		// https://github.com/open-telemetry/opentelemetry-go/issues/1280
		m.serverStreamReceiveTimeHistogram = metric.Must(Meter()).NewFloat64Histogram("grpc_server_msg_recv_handling_seconds", options...)
	}
	m.serverStreamReceiveTimeHistogramEnabled = true
}

// EnableServerStreamReceiveSizeHistogram turns on recording of single message receive size of streaming RPCs.
// Histogram metrics can be very expensive for Prometheus to retain and query.
func (m *ServerMetrics) EnableServerStreamReceiveSizeHistogram(opts ...metric.InstrumentOption) {
	var options []metric.InstrumentOption
	options = append(options,
		metric.WithDescription("Histogram of message size (bytes) of the gRPC single message receive."),
		metric.WithUnit(unit.Bytes))
	options = append(options, opts...)
	if !m.serverStreamReceiveSizeHistogramEnabled {
		// https://github.com/open-telemetry/opentelemetry-go/issues/1280
		m.serverStreamReceiveSizeHistogram = metric.Must(Meter()).NewInt64Histogram("grpc_server_msg_recv_handling_bytes", options...)
	}
	m.serverStreamReceiveSizeHistogramEnabled = true
}

// EnableServerStreamSendTimeHistogram turns on recording of single message send time of streaming RPCs.
// Histogram metrics can be very expensive for Prometheus to retain and query.
func (m *ServerMetrics) EnableServerStreamSendTimeHistogram(opts ...metric.InstrumentOption) {
	var options []metric.InstrumentOption
	options = append(options,
		metric.WithDescription("Histogram of response latency (seconds) of the gRPC single message send."),
		metric.WithUnit("s"))
	options = append(options, opts...)
	if !m.serverStreamSendTimeHistogramEnabled {
		// https://github.com/open-telemetry/opentelemetry-go/issues/1280
		m.serverStreamSendTimeHistogram = metric.Must(Meter()).NewFloat64Histogram("grpc_server_msg_send_handling_seconds", options...)
	}
	m.serverStreamSendTimeHistogramEnabled = true
}

// EnableServerStreamSendSizeHistogram turns on recording of single message send size of streaming RPCs.
// Histogram metrics can be very expensive for Prometheus to retain and query.
func (m *ServerMetrics) EnableServerStreamSendSizeHistogram(opts ...metric.InstrumentOption) {
	var options []metric.InstrumentOption
	options = append(options,
		metric.WithDescription("Histogram of message size (bytes) of the gRPC single message send."),
		metric.WithUnit(unit.Bytes))
	options = append(options, opts...)
	if !m.serverStreamSendSizeHistogramEnabled {
		// https://github.com/open-telemetry/opentelemetry-go/issues/1280
		m.serverStreamSendSizeHistogram = metric.Must(Meter()).NewInt64Histogram("grpc_server_msg_send_handling_bytes", options...)
	}
	m.serverStreamSendSizeHistogramEnabled = true
}

// UnaryServerInterceptor is a gRPC server-side interceptor that provides Prometheus monitoring for Unary RPCs.
func (m *ServerMetrics) UnaryServerInterceptor() func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		monitor := newServerReporter(ctx, m, Unary, info.FullMethod, peerFromCtx(ctx), m.ServerHostport)
		monitor.ReceivedMessage(ctx, req)
		resp, err := handler(ctx, req)
		st, _ := FromError(err)
		monitor.Handled(ctx, st.Code())
		if err == nil {
			monitor.SentMessage(ctx, resp)
		}
		return resp, err
	}
}

// StreamServerInterceptor is a gRPC server-side interceptor that provides Prometheus monitoring for Streaming RPCs.
func (m *ServerMetrics) StreamServerInterceptor() func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		monitor := newServerReporter(ss.Context(), m, streamRPCType(info), info.FullMethod, peerFromCtx(ss.Context()), m.ServerHostport)
		err := handler(srv, &monitoredServerStream{ss, monitor})
		st, _ := FromError(err)
		monitor.Handled(ss.Context(), st.Code())
		return err
	}
}

// InitializeMetrics initializes all metrics, with their appropriate null
// value, for all gRPC methods registered on a gRPC server. This is useful, to
// ensure that all metrics exist when collecting and querying.
func (m *ServerMetrics) InitializeMetrics(ctx context.Context, server *grpc.Server) {
	serviceInfo := server.GetServiceInfo()
	for serviceName, info := range serviceInfo {
		for _, mInfo := range info.Methods {
			preRegisterMethod(ctx, m, serviceName, &mInfo)
		}
	}
}

func streamRPCType(info *grpc.StreamServerInfo) grpcType {
	if info.IsServerStream && !info.IsClientStream {
		return ServerStream
	} else if info.IsClientStream && !info.IsServerStream {
		return ClientStream
	}
	return BidiStream
}

// monitoredStream wraps grpc.ServerStream allowing each Sent/Recv of message to increment counters.
type monitoredServerStream struct {
	grpc.ServerStream
	monitor *serverReporter
}

func (s *monitoredServerStream) SendMsg(m interface{}) error {
	now := time.Now()
	err := s.ServerStream.SendMsg(m)
	s.monitor.SendMessageTimer(context.Background(), now)
	if err == nil {
		s.monitor.SentMessage(context.Background(), m)
	}
	return err
}

func (s *monitoredServerStream) RecvMsg(m interface{}) error {
	now := time.Now()
	err := s.ServerStream.RecvMsg(m)
	s.monitor.ReceiveMessageTimer(context.Background(), now)
	if err == nil {
		s.monitor.ReceivedMessage(context.Background(), m)
	}
	return err
}

// preRegisterMethod is invoked on Register of a Server, allowing all gRPC services labels to be pre-populated.
func preRegisterMethod(ctx context.Context, metrics *ServerMetrics, serviceName string, mInfo *grpc.MethodInfo) {
	filter := AttrsFilter
	if filter == nil {
		filter = func(attrs ...attribute.KeyValue) []attribute.KeyValue { return attrs }
	}
	// These are just references (no increments), as just referencing will create the labels but not set values.
	_, attrs := spanInfo(mInfo.Name, ":0", metrics.ServerHostport, typeFromMethodInfo(mInfo), false)
	metrics.serverStartedCounter.Add(ctx, 0, filter(attrs...)...)
	metrics.serverStreamMsgReceived.Add(ctx, 0, filter(append(attrs, otelgrpc_.RPCMessageTypeReceived)...)...)
	metrics.serverStreamMsgSent.Add(ctx, 0, filter(append(attrs, otelgrpc_.RPCMessageTypeSent)...)...)

	for _, code := range allCodes {
		metrics.serverHandledCounter.Add(ctx, 0, filter(append(attrs, statusCodeAttr(code))...)...)
		if metrics.serverHandledTimeHistogramEnabled {
			metrics.serverHandledTimeHistogram.Record(ctx, -1, filter(append(attrs, statusCodeAttr(code))...)...)
		}
	}
	if metrics.serverStreamReceiveTimeHistogramEnabled {
		metrics.serverStreamReceiveTimeHistogram.Record(ctx, -1, filter(append(attrs, otelgrpc_.RPCMessageTypeReceived)...)...)
	}
	if metrics.serverStreamReceiveSizeHistogramEnabled {
		metrics.serverStreamReceiveSizeHistogram.Record(ctx, -1, filter(append(attrs, otelgrpc_.RPCMessageTypeReceived)...)...)
	}

	if metrics.serverStreamSendTimeHistogramEnabled {
		metrics.serverStreamSendTimeHistogram.Record(ctx, -1, filter(append(attrs, otelgrpc_.RPCMessageTypeSent)...)...)
	}
	if metrics.serverStreamSendSizeHistogramEnabled {
		metrics.serverStreamSendSizeHistogram.Record(ctx, -1, filter(append(attrs, otelgrpc_.RPCMessageTypeSent)...)...)
	}
}
