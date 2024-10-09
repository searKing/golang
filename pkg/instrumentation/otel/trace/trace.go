// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package trace

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/searKing/golang/pkg/instrumentation/otel/trace/processor"
)

func NewTracerProvider(ctx context.Context, options ...Option) (*sdktrace.TracerProvider, error) {
	var o option
	o.SetDefaults()
	o.ApplyOptions(options...)

	var tracerProviderOptions []sdktrace.TracerProviderOption
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
		tracerProviderOptions = append(tracerProviderOptions, sdktrace.WithResource(res))
	}
	// add span attributes to each span, so that we can see global span attributes in all Spans.
	{
		tracerProviderOptions = append(tracerProviderOptions, sdktrace.WithSpanProcessor(processor.Annotator{
			AttrsFunc: func() []attribute.KeyValue {
				attrsSet, _ := attribute.NewSetWithFiltered(o.SpanAttrs, func(kv attribute.KeyValue) bool {
					if !kv.Valid() {
						return false
					}
					if kv.Value.Emit() == "" {
						return false
					}
					return true
				})
				return attrsSet.ToSlice()
			},
		}))
	}

	{
		exporters, err := createExporters(ctx, options...)
		if err != nil {
			return nil, err
		}
		o.Exporters = append(o.Exporters, exporters...)
		for _, exporter := range o.Exporters {
			tracerProviderOptions = append(tracerProviderOptions, sdktrace.WithBatcher(exporter))
		}
	}

	traceProvider := sdktrace.NewTracerProvider(tracerProviderOptions...)
	initPassthroughGlobals()
	return traceProvider, nil
}

func createExporters(ctx context.Context, opts ...Option) ([]sdktrace.SpanExporter, error) {
	var o option
	o.SetDefaults()
	o.ApplyOptions(opts...)

	var exporters []sdktrace.SpanExporter
	for _, v := range o.ExporterEndpoints {
		u, err := url.Parse(v)
		if err != nil {
			return nil, fmt.Errorf("malformed trace exporter endpoint %s: %w", v, err)
		}

		{
			opener := Get(u.Scheme)
			if opener == nil {
				return nil, fmt.Errorf("unknown trace exporter scheme: %s", u.Scheme)
			}
			exporter, err := opener.OpenExporterURL(ctx, u)
			if err != nil {
				return nil, err
			}
			exporters = append(exporters, exporter)
		}
	}

	return exporters, nil
}

func initPassthroughGlobals() {
	// We explicitly DO NOT set the global TracerProvider using otel.SetTracerProvider().
	// The unset TracerProvider returns a "non-recording" span, but still passes through context.
	// See also: https://github.com/open-telemetry/opentelemetry-go-contrib/blob/main/examples/passthrough
	slog.Info(`Register a global TextMapPropagator, but do not register a global TracerProvider to be in "passthrough" mode.`)
	slog.Info(`The "passthrough" mode propagates the TraceContext and Baggage, but does not record spans.`)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
}
