// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package prometheusmetric

import (
	"context"
	"net/url"

	slices_ "github.com/searKing/golang/go/exp/slices"
	"go.opentelemetry.io/otel/attribute"
	prometheusmetric "go.opentelemetry.io/otel/exporters/prometheus"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"

	url_ "github.com/searKing/golang/pkg/instrumentation/otel/url"
)

// URLOpener opens Prometheus Metric URLs like "http://endpoint?deny_keys=[]".
type URLOpener struct {
	// Options specifies the options to pass to OpenReader.
	Option option
}

// Scheme returns the scheme supported by this metric.
func (o *URLOpener) Scheme() string { return "prometheus" }

// OpenReaderURL opens a metric.Exporter based on u.
func (o *URLOpener) OpenReaderURL(ctx context.Context, u *url.URL) (sdkmetric.Reader, error) {
	var opts []Option
	q := u.Query()
	u.RawQuery = ""
	u.RawFragment = ""
	{
		prometheusOpts, err := parsePrometheusOpts(q, o.Option.PrometheusOptions...)
		if err != nil {
			return nil, err
		}
		opts = append(opts, WithOptionPrometheusOptions(prometheusOpts...))
	}
	return OpenReader(ctx, opts...)
}

func parsePrometheusOpts(q url.Values, opts ...prometheusmetric.Option) ([]prometheusmetric.Option, error) {
	opts = append(opts, prometheusmetric.WithResourceAsConstantLabels(attribute.NewDenyKeysFilter()))

	{
		ns := q.Get("namespace")
		if ns != "" {
			opts = append(opts, prometheusmetric.WithNamespace(ns))
		}
		q.Del("namespace")
	}
	if q.Has("deny_keys") {
		opts = append(opts, prometheusmetric.WithResourceAsConstantLabels(
			attribute.NewDenyKeysFilter(slices_.MapFunc(q["deny_keys"], func(e string) attribute.Key {
				return attribute.Key(e)
			})...)))
		q.Del("deny_keys")
	}
	if q.Has("allow_keys") {
		opts = append(opts, prometheusmetric.WithResourceAsConstantLabels(
			attribute.NewAllowKeysFilter(slices_.MapFunc(q["allow_keys"], func(e string) attribute.Key {
				return attribute.Key(e)
			})...)))
		q.Del("allow_keys")
	}
	{
		b, err := url_.ParseBoolFromValues(q, "no_counter_suffixes")
		if err != nil {
			return nil, err
		}
		if b {
			opts = append(opts, prometheusmetric.WithoutCounterSuffixes())
		}
		q.Del("no_counter_suffixes")
	}
	{
		b, err := url_.ParseBoolFromValues(q, "no_scope_info")
		if err != nil {
			return nil, err
		}
		if b {
			opts = append(opts, prometheusmetric.WithoutScopeInfo())
		}
		q.Del("no_scope_info")
	}
	{
		b, err := url_.ParseBoolFromValues(q, "no_target_info")
		if err != nil {
			return nil, err
		}
		if b {
			opts = append(opts, prometheusmetric.WithoutTargetInfo())
		}
		q.Del("no_target_info")
	}
	{
		b, err := url_.ParseBoolFromValues(q, "no_units")
		if err != nil {
			return nil, err
		}
		if b {
			opts = append(opts, prometheusmetric.WithoutUnits())
		}
		q.Del("no_units")
	}
	return opts, nil
}
