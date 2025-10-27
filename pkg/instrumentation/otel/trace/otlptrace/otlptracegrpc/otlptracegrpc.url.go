// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package otlptracegrpc

import (
	"context"
	"fmt"
	"net/url"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	url_ "github.com/searKing/golang/pkg/instrumentation/otel/url"
)

// URLOpener opens OTLP Trace gRPC URLs like "http://endpoint:4317?compression=gzip".
type URLOpener struct {
	// Options specifies the options to pass to OpenExporter.
	Options []Option
}

// Scheme returns the scheme supported by this trace exporter.
func (o *URLOpener) Scheme() string { return "otlp-grpc" }

// OpenExporterURL opens a trace.SpanExporter based on u.
func (o *URLOpener) OpenExporterURL(ctx context.Context, u *url.URL) (sdktrace.SpanExporter, error) {
	q := u.Query()
	u.RawQuery = ""
	u.RawFragment = ""

	{
		scheme, err := parseScheme(q)
		if err != nil {
			return nil, err
		}
		u.Scheme = scheme
	}
	opts := o.Options
	{
		var otlpOpts []otlptracegrpc.Option
		otlpOpts = append(otlpOpts, otlptracegrpc.WithEndpointURL(u.String()))
		{
			opts, err := parseOtlpOpts(q)
			if err != nil {
				return nil, err
			}
			otlpOpts = append(otlpOpts, opts...)
		}
		opts = append(opts, WithOptionOtlpOptions(otlpOpts...))
	}
	return OpenExporter(ctx, opts...)
}

func parseScheme(q url.Values) (scheme string, err error) {
	b, err := url_.ParseBoolFromValues(q, "insecure")
	if err != nil {
		return "", err
	}
	if b {
		scheme = "http"
	} else {
		scheme = "https"
	}
	q.Del("insecure")
	return
}

func parseOtlpOpts(q url.Values, opts ...otlptracegrpc.Option) ([]otlptracegrpc.Option, error) {
	{
		v := q.Get("compression")
		switch v {
		case "gzip":
			opts = append(opts, otlptracegrpc.WithCompressor("gzip"))
		case "none":
			opts = append(opts, otlptracegrpc.WithCompressor("none"))
		case "":
		default:
			return nil, fmt.Errorf("unknown quary parameter compression: %s", v)
		}
		q.Del("compression")
	}
	return opts, nil
}
