// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package otelgrpc

import (
	"context"
	"net"
	"net/http"
	"time"

	net_ "github.com/searKing/golang/go/net"
	otelgrpc_ "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type clientReporter struct {
	metrics   *ClientMetrics
	attrs     []attribute.KeyValue
	startTime time.Time
}

func newClientReporter(ctx context.Context, m *ClientMetrics, req *http.Request) *clientReporter {
	if m.ClientHostport == "" {
		var addr string
		if ip, err := net_.ListenIP(); err == nil {
			addr = ip.String()
		}
		m.ClientHostport = net.JoinHostPort(addr, "0")
	}

	r := &clientReporter{
		metrics: m,
	}
	if r.metrics.clientHandledTimeHistogramEnabled {
		r.startTime = time.Now()
	}

	attrs := AttrsFromRequest(req, m.ClientHostport)
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

func (r *clientReporter) ReceivedResponse(ctx context.Context, resp *http.Response) {
	attrs := r.Attrs(otelgrpc_.RPCMessageTypeReceived)
	r.metrics.clientStreamRequestReceived.Add(ctx, 1, metric.WithAttributes(attrs...))
	if r.metrics.clientStreamReceiveSizeHistogramEnabled {
		if resp != nil {
			r.metrics.clientStreamReceiveSizeHistogram.Record(ctx, resp.ContentLength, metric.WithAttributes(attrs...))
		} else {
			r.metrics.clientStreamReceiveSizeHistogram.Record(ctx, -1, metric.WithAttributes(attrs...))
		}
	}
}

func (r *clientReporter) SentRequest(ctx context.Context, req *http.Request) {
	attrs := r.Attrs(otelgrpc_.RPCMessageTypeSent)

	r.metrics.clientStreamRequestSent.Add(ctx, 1, metric.WithAttributes(attrs...))
	if r.metrics.clientStreamSendSizeHistogramEnabled {
		if req != nil {
			r.metrics.clientStreamSendSizeHistogram.Record(ctx, req.ContentLength, metric.WithAttributes(attrs...))
		} else {
			r.metrics.clientStreamSendSizeHistogram.Record(ctx, -1, metric.WithAttributes(attrs...))
		}
	}
}

func (r *clientReporter) Handled(ctx context.Context, resp *http.Response) {
	attrs := r.Attrs(AttrsFromResponse(resp)...)
	r.metrics.clientHandledCounter.Add(ctx, 1, metric.WithAttributes(attrs...))

	if r.metrics.clientHandledTimeHistogramEnabled {
		r.metrics.clientHandledTimeHistogram.Record(ctx, time.Since(r.startTime).Seconds(), metric.WithAttributes(attrs...))
	}
}
