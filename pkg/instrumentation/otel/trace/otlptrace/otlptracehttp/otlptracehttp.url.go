// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package otlptracehttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	url_ "github.com/searKing/golang/pkg/instrumentation/otel/url"
)

// URLOpener opens OTLP Trace HTTP URLs like "otlp-http://endpoint".
type URLOpener struct {
	// Options specifies the options to pass to OpenExporter.
	Option option
}

// Scheme returns the scheme supported by this trace exporter.
func (o *URLOpener) Scheme() string { return "otlp-http" }

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

	var otlpOpts []otlptracehttp.Option
	otlpOpts = append(otlpOpts, otlptracehttp.WithEndpointURL(u.String()))
	{
		opts, err := parseOtlpOpts(q, o.Option.OtlpOptions...)
		if err != nil {
			return nil, err
		}
		otlpOpts = append(otlpOpts, opts...)
	}
	return OpenExporter(ctx, WithOptionOtlpOptions(otlpOpts...))
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

func parseOtlpOpts(q url.Values, opts ...otlptracehttp.Option) ([]otlptracehttp.Option, error) {
	{
		v := q.Get("compression")
		switch v {
		case "gzip":
			opts = append(opts, otlptracehttp.WithCompression(otlptracehttp.GzipCompression))
		case "none":
			opts = append(opts, otlptracehttp.WithCompression(otlptracehttp.NoCompression))
		case "":
		default:
			return nil, fmt.Errorf("unknown quary parameter compression: %s", v)
		}
		q.Del("compression")
	}
	{
		v := q["headers"]
		headers := make(map[string]string)
		for _, data := range v {
			err := json.Unmarshal([]byte(data), &headers)
			if err != nil {
				return nil, fmt.Errorf("unknown quary parameter headers: %w", err)
			}
		}
		if len(headers) > 0 {
			opts = append(opts, otlptracehttp.WithHeaders(headers))
		}
		q.Del("headers")
	}
	return opts, nil
}
