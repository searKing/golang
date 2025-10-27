// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stdoutmetric

import (
	"context"
	"io"
	"net/url"
	"os"

	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"

	url_ "github.com/searKing/golang/pkg/instrumentation/otel/url"
)

// URLOpener opens stdout Metric URLs like "http://localhost?allow_stdout&pretty_print&no_timestamps&periodic_reader_interval=60s".
type URLOpener struct {
	// Options specifies the options to pass to OpenReaderURL.
	Options []Option
}

// Scheme returns the scheme supported by this metric.
func (o *URLOpener) Scheme() string { return "stdout" }

// OpenReaderURL opens a metric.Reader based on u.
func (o *URLOpener) OpenReaderURL(ctx context.Context, u *url.URL) (sdkmetric.Reader, error) {
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
	{
		readerOpts, err := parseReaderOpts(q)
		if err != nil {
			return nil, err
		}
		opts = append(opts, WithOptionPeriodicReaderOptions(readerOpts...))
	}
	return OpenReader(ctx, opts...)
}

func parseStdoutOpts(q url.Values, opts ...stdoutmetric.Option) ([]stdoutmetric.Option, error) {
	{
		b, err := url_.ParseBoolFromValues(q, "allow_stdout")
		if err != nil {
			return nil, err
		}
		w := io.Discard
		if b {
			w = os.Stdout
		}
		opts = append(opts, stdoutmetric.WithWriter(w))
		q.Del("allow_stdout")
	}
	{
		b, err := url_.ParseBoolFromValues(q, "pretty_print")
		if err != nil {
			return nil, err
		}
		if b {
			opts = append(opts, stdoutmetric.WithPrettyPrint())
		}
		q.Del("pretty_print")
	}
	{
		b, err := url_.ParseBoolFromValues(q, "no_timestamps")
		if err != nil {
			return nil, err
		}
		if b {
			opts = append(opts, stdoutmetric.WithoutTimestamps())
		}
		q.Del("no_timestamps")
	}
	return opts, nil
}

func parseReaderOpts(q url.Values, opts ...sdkmetric.PeriodicReaderOption) ([]sdkmetric.PeriodicReaderOption, error) {
	{
		d, err := url_.ParseTimeDurationFromValues(q, "periodic_reader_interval")
		if err != nil {
			return nil, err
		}
		if d != 0 {
			opts = append(opts, sdkmetric.WithInterval(d))
		}
		q.Del("periodic_reader_interval")
	}
	return opts, nil
}
