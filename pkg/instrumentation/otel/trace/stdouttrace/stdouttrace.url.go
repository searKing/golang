// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stdouttrace

import (
	"context"
	"io"
	"net/url"
	"os"

	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	url_ "github.com/searKing/golang/pkg/instrumentation/otel/url"
)

// URLOpener opens stdout Metric URLs like "http://localhost?allow_stdout&pretty_print&no_timestamps".
type URLOpener struct {
	// Options specifies the options to pass to OpenExporter.
	Options []Option
}

// Scheme returns the scheme supported by this trace exporter.
func (o *URLOpener) Scheme() string { return "stdout" }

// OpenExporterURL opens a trace.Exporter based on u.
func (o *URLOpener) OpenExporterURL(ctx context.Context, u *url.URL) (sdktrace.SpanExporter, error) {
	q := u.Query()
	u.RawQuery = ""
	u.RawFragment = ""

	opts := o.Options
	{
		stdoutOpts, err := parseStdoutOpts(q)
		if err != nil {
			return nil, err
		}
		opts = append(opts, WithOptionStdoutOptions(stdoutOpts...))
	}
	return OpenExporter(ctx, opts...)
}

func parseStdoutOpts(q url.Values, opts ...stdouttrace.Option) ([]stdouttrace.Option, error) {
	{
		b, err := url_.ParseBoolFromValues(q, "allow_stdout")
		if err != nil {
			return nil, err
		}
		w := io.Discard
		if b {
			w = os.Stdout
		}
		opts = append(opts, stdouttrace.WithWriter(w))
		q.Del("allow_stdout")
	}
	{
		b, err := url_.ParseBoolFromValues(q, "pretty_print")
		if err != nil {
			return nil, err
		}
		if b {
			opts = append(opts, stdouttrace.WithPrettyPrint())
		}
		q.Del("pretty_print")
	}
	{
		b, err := url_.ParseBoolFromValues(q, "no_timestamps")
		if err != nil {
			return nil, err
		}
		if b {
			opts = append(opts, stdouttrace.WithoutTimestamps())
		}
		q.Del("no_timestamps")
	}
	return opts, nil
}
