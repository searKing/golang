// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package otlpmetricgrpc

import (
	"context"
	"fmt"
	"net/url"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"

	url_ "github.com/searKing/golang/pkg/instrumentation/otel/url"
)

// URLOpener opens OTLP Metric gRPC URLs like "http://endpoint:4317?compression=gzip&temporality_selector=delta".
type URLOpener struct {
	// Options specifies the options to pass to OpenExporter.
	Option option
}

// Scheme returns the scheme supported by this metric.
func (o *URLOpener) Scheme() string { return "otlp-grpc" }

// OpenReaderURL opens a metric.Reader based on u.
func (o *URLOpener) OpenReaderURL(ctx context.Context, u *url.URL) (sdkmetric.Reader, error) {
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

	var opts []Option
	{
		var otlpOpts []otlpmetricgrpc.Option
		otlpOpts = append(otlpOpts, otlpmetricgrpc.WithEndpointURL(u.String()))
		{
			var err error
			otlpOpts, err = parseOtlpOpts(q, otlpOpts...)
			if err != nil {
				return nil, err
			}
		}
		opts = append(opts, WithOptionOtlpOptions(otlpOpts...))
	}

	{
		var readerOpts []sdkmetric.PeriodicReaderOption
		readerOpts, err := parseReaderOpts(q, readerOpts...)
		if err != nil {
			return nil, err
		}
		opts = append(opts, WithOptionPeriodicReaderOptions(readerOpts...))
	}
	return OpenReader(ctx, opts...)
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

func parseOtlpOpts(q url.Values, opts ...otlpmetricgrpc.Option) ([]otlpmetricgrpc.Option, error) {
	{
		v := q.Get("compression")
		switch v {
		case "gzip":
			opts = append(opts, otlpmetricgrpc.WithCompressor("gzip"))
		case "none":
			opts = append(opts, otlpmetricgrpc.WithCompressor("none"))
		case "":
		default:
			return nil, fmt.Errorf("unknown quary parameter compression: %s", v)
		}
		q.Del("compression")
	}

	{
		v := q.Get("temporality_selector")
		switch v {
		case "delta":
			opts = append(opts, otlpmetricgrpc.WithTemporalitySelector(func(kind sdkmetric.InstrumentKind) metricdata.Temporality { return metricdata.DeltaTemporality }))
		case "cumulative":
			opts = append(opts, otlpmetricgrpc.WithTemporalitySelector(func(kind sdkmetric.InstrumentKind) metricdata.Temporality { return metricdata.CumulativeTemporality }))
		case "default":
			opts = append(opts, otlpmetricgrpc.WithTemporalitySelector(sdkmetric.DefaultTemporalitySelector))
		case "":
		default:
			return nil, fmt.Errorf("unknown quary parameter temporality_selector: %s", v)
		}
		q.Del("temporality_selector")
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
