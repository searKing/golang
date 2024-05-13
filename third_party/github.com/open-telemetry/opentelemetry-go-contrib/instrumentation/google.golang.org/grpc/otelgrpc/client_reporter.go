// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package otelgrpc

import (
	"context"
	"time"

	otelgrpc_ "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"
)

type clientReporter struct {
	metrics   *ClientMetrics
	attrs     []attribute.KeyValue
	startTime time.Time
}

func newClientReporter(ctx context.Context, m *ClientMetrics, rpcType grpcType, fullMethod string, peerAddress string, localAddress string) *clientReporter {
	r := &clientReporter{
		metrics: m,
	}
	if r.metrics.clientHandledTimeHistogramEnabled {
		r.startTime = time.Now()
	}

	_, attrs := spanInfo(fullMethod, peerAddress, localAddress, rpcType, true)
	r.attrs = attrs
	r.metrics.clientStartedCounter.Add(ctx, 1, metric.WithAttributes(r.Attrs()...))
	return r
}

func (r *clientReporter) Attrs(attrs ...attribute.KeyValue) []attribute.KeyValue {
	attrs = append(r.attrs, attrs...)
	filter := AttrsFilter
	if filter != nil {
		return filter(attrs...)
	}
	return attrs
}

func (r *clientReporter) ReceiveMessageTimer(ctx context.Context, startTime time.Time) {
	if r.metrics.clientStreamReceiveTimeHistogramEnabled {
		attrs := r.Attrs(otelgrpc_.RPCMessageTypeReceived)
		r.metrics.clientStreamReceiveTimeHistogram.Record(ctx, time.Since(startTime).Seconds(), metric.WithAttributes(attrs...))
		return
	}
	return
}

func (r *clientReporter) ReceivedMessage(ctx context.Context, message any) {
	attrs := r.Attrs(otelgrpc_.RPCMessageTypeReceived)
	r.metrics.clientStreamMsgReceived.Add(ctx, 1, metric.WithAttributes(attrs...))
	if r.metrics.clientStreamReceiveSizeHistogramEnabled {
		if p, ok := message.(proto.Message); ok {
			r.metrics.clientStreamReceiveSizeHistogram.Record(ctx, int64(proto.Size(p)), metric.WithAttributes(attrs...))
		} else {
			r.metrics.clientStreamReceiveSizeHistogram.Record(ctx, -1, metric.WithAttributes(attrs...))
		}
	}
}

func (r *clientReporter) SendMessageTimer(ctx context.Context, startTime time.Time) {
	if r.metrics.clientStreamSendTimeHistogramEnabled {
		attrs := r.Attrs(otelgrpc_.RPCMessageTypeSent)
		r.metrics.clientStreamSendTimeHistogram.Record(ctx, time.Since(startTime).Seconds(), metric.WithAttributes(attrs...))
		return
	}
	return
}

func (r *clientReporter) SentMessage(ctx context.Context, message any) {
	attrs := r.Attrs(otelgrpc_.RPCMessageTypeSent)
	r.metrics.clientStreamMsgSent.Add(ctx, 1, metric.WithAttributes(attrs...))
	if r.metrics.clientStreamSendSizeHistogramEnabled {
		if p, ok := message.(proto.Message); ok {
			r.metrics.clientStreamSendSizeHistogram.Record(ctx, int64(proto.Size(p)), metric.WithAttributes(attrs...))
		} else {
			r.metrics.clientStreamSendSizeHistogram.Record(ctx, int64(proto.Size(nil)), metric.WithAttributes(attrs...))
		}
	}
}

func (r *clientReporter) Handled(ctx context.Context, code codes.Code) {
	attrs := r.Attrs(statusCodeAttr(code))
	r.metrics.clientHandledCounter.Add(ctx, 1, metric.WithAttributes(attrs...))

	if r.metrics.clientHandledTimeHistogramEnabled {
		r.metrics.clientHandledTimeHistogram.Record(ctx, time.Since(r.startTime).Seconds(), metric.WithAttributes(attrs...))
	}
}
