// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package metric

import (
	"context"
	"fmt"
	"net/url"

	"go.opentelemetry.io/otel/sdk/metric"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
)

func NewMeterProvider(ctx context.Context, options ...Option) (*metric.MeterProvider, error) {
	var o option
	o.SetDefaults()
	o.ApplyOptions(options...)

	var metricProviderOptions []metric.Option
	{
		readers, err := createReaders(ctx, options...)
		if err != nil {
			return nil, err
		}
		o.Readers = append(o.Readers, readers...)
	}

	for _, reader := range o.Readers {
		metricProviderOptions = append(metricProviderOptions, metric.WithReader(reader)) // 默认cumulative
	}
	//{
	//	res := sdkresource.NewSchemaless(o.ResourceAttrs...)
	//	metricProviderOptions = append(metricProviderOptions, metric.WithResource(res))
	//}
	{
		res, err := sdkresource.New(ctx,
			sdkresource.WithFromEnv(),      // Discover and provide attributes from OTEL_RESOURCE_ATTRIBUTES and OTEL_SERVICE_NAME environment variables.
			sdkresource.WithTelemetrySDK(), // Discover and provide information about the OpenTelemetry SDK used.
			sdkresource.WithProcess(),      // Discover and provide process information.
			sdkresource.WithOS(),           // Discover and provide OS information.
			sdkresource.WithContainer(),    // Discover and provide container information.
			sdkresource.WithHost(),         // Discover and provide host information.
			sdkresource.WithProcess(),
			// sdkresource.WithDetectors(thirdparty.Detector{}), // Bring your own external Detector implementation.
			sdkresource.WithAttributes(o.ResourceAttrs...)) // Add custom resource attributes.
		if err != nil {
			return nil, err
		}
		metricProviderOptions = append(metricProviderOptions, metric.WithResource(res))
	}
	return metric.NewMeterProvider(metricProviderOptions...), nil
}

func createReaders(ctx context.Context, opts ...Option) ([]metric.Reader, error) {
	var o option
	o.SetDefaults()
	o.ApplyOptions(opts...)

	var readers []metric.Reader
	for _, v := range o.ExporterEndpoints {
		u, err := url.Parse(v)
		if err != nil {
			return nil, fmt.Errorf("malformed metric exporter endpoint %s: %w", v, err)
		}
		{
			// handle reader, as manual puller
			opener := Get(u.Scheme)
			if opener == nil {
				return nil, fmt.Errorf("unknown metric exporter scheme: %s", u.Scheme)
			}
			reader, err := opener.OpenReaderURL(ctx, u)
			if err != nil {
				return nil, err
			}
			readers = append(readers, reader)
			continue
		}
	}

	return readers, nil
}
