// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package otelgrpc

import (
	"context"
	"time"

	otelgrpc_ "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"
)

type serverReporter struct {
	metrics   *ServerMetrics
	attrs     []attribute.KeyValue
	startTime time.Time
}

func newServerReporter(ctx context.Context, m *ServerMetrics, rpcType grpcType, fullMethod string, peerAddress string, localAddress string) *serverReporter {
	r := &serverReporter{
		metrics: m,
	}
	if r.metrics.serverHandledTimeHistogramEnabled {
		r.startTime = time.Now()
	}
	_, attrs := spanInfo(fullMethod, peerAddress, localAddress, rpcType, false)
	r.attrs = attrs
	r.metrics.serverStartedCounter.Add(ctx, 1, r.Attrs()...)
	return r
}

func (r *serverReporter) Attrs(attrs ...attribute.KeyValue) []attribute.KeyValue {
	attrs = append(r.attrs, attrs...)
	filter := AttrsFilter
	if filter != nil {
		return filter(attrs...)
	}
	return attrs
}

func (r *serverReporter) ReceiveMessageTimer(ctx context.Context, startTime time.Time) {
	if r.metrics.serverStreamReceiveTimeHistogramEnabled {
		attrs := r.Attrs(otelgrpc_.RPCMessageTypeReceived)
		r.metrics.serverStreamReceiveTimeHistogram.Record(ctx, time.Since(startTime).Seconds(), attrs...)
		return
	}
	return
}

func (r *serverReporter) ReceivedMessage(ctx context.Context, message interface{}) {
	attrs := r.Attrs(otelgrpc_.RPCMessageTypeReceived)
	r.metrics.serverStreamMsgReceived.Add(ctx, 1, attrs...)
	if r.metrics.serverStreamReceiveSizeHistogramEnabled {
		if p, ok := message.(proto.Message); ok {
			r.metrics.serverStreamReceiveSizeHistogram.Record(ctx, int64(proto.Size(p)), attrs...)
		} else {
			r.metrics.serverStreamReceiveSizeHistogram.Record(ctx, -1, attrs...)
		}
	}
}

func (r *serverReporter) SendMessageTimer(ctx context.Context, startTime time.Time) {
	if r.metrics.serverStreamSendTimeHistogramEnabled {
		attrs := r.Attrs(otelgrpc_.RPCMessageTypeSent)
		r.metrics.serverStreamSendTimeHistogram.Record(ctx, time.Since(startTime).Seconds(), attrs...)
		return
	}
	return
}

func (r *serverReporter) SentMessage(ctx context.Context, message interface{}) {
	attrs := r.Attrs(otelgrpc_.RPCMessageTypeSent)
	r.metrics.serverStreamMsgSent.Add(ctx, 1, attrs...)
	if r.metrics.serverStreamSendSizeHistogramEnabled {
		if p, ok := message.(proto.Message); ok {
			r.metrics.serverStreamSendSizeHistogram.Record(ctx, int64(proto.Size(p)), attrs...)
		} else {
			r.metrics.serverStreamSendSizeHistogram.Record(ctx, int64(proto.Size(nil)), attrs...)
		}
	}
}

func (r *serverReporter) Handled(ctx context.Context, code codes.Code) {
	attrs := r.Attrs(statusCodeAttr(code))
	r.metrics.serverHandledCounter.Add(ctx, 1, attrs...)

	if r.metrics.serverHandledTimeHistogramEnabled {
		r.metrics.serverHandledTimeHistogram.Record(ctx, time.Since(r.startTime).Seconds(), attrs...)
	}
}
