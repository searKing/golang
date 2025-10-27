// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package metric_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-logr/logr/funcr"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"

	metric_ "github.com/searKing/golang/pkg/instrumentation/otel/metric"

	_ "github.com/searKing/golang/pkg/instrumentation/otel/metric/otlpmetric/otlpmetricgrpc" // for otlp-grpc
	_ "github.com/searKing/golang/pkg/instrumentation/otel/metric/otlpmetric/otlpmetrichttp" // for otlp-http
	_ "github.com/searKing/golang/pkg/instrumentation/otel/metric/prometheusmetric"          // for prometheus
	_ "github.com/searKing/golang/pkg/instrumentation/otel/metric/stdoutmetric"              // for stdout
)

const testLoop = 1
const testInterval = 0 * time.Minute
const testForcePush = false
const instrumentation = "github.com/searKing/golang/pkg/instrumentation/otel/metric"

func TestNewMeterProvider(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	otel.SetLogger(funcr.New(func(prefix, args string) { t.Logf("otel: %s", fmt.Sprint(prefix, args)) }, funcr.Options{Verbosity: 1}))
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) { t.Errorf("otel: handler returned an error: %s", err.Error()) }))
	//otel.SetLogger(stdr.New(log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)))
	//otel.SetLogger(stdr.New(slog.NewLogLogger(slog.Default().Handler(), slog.LevelWarn)))
	mp, err := metric_.NewMeterProvider(ctx, metric_.WithOptionExporterEndpoints(
		"stdout://localhost?allow_stdout&pretty_print&no_timestamps",
		"prometheus://localhost",
		//`otlp-http://some_endpoint/some_path?compression=gzip&insecure`,
		//`otlp-grpc://some_endpoint/some_path?compression=gzip&insecure`,
	), metric_.WithOptionResourceAttrs())
	if err != nil {
		t.Fatalf("create meter provider failed: %s", err.Error())
		return
	}
	otel.SetMeterProvider(mp)
	defer func() {
		err := mp.Shutdown(context.Background())
		if err != nil {
			t.Fatalf("shutdown meter provider failed: %s", err.Error())
			return
		}
	}()
	defer func() {
		err := mp.ForceFlush(context.Background())
		if err != nil {
			t.Fatalf("force plush meter provider failed: %s", err.Error())
			return
		}
	}()
	for range testLoop {
		for range testLoop {
			testAllMetrics(ctx, t)
		}
		if testForcePush {
			err := mp.ForceFlush(context.Background())
			if err != nil {
				t.Fatalf("force plush meter provider failed: %s", err.Error())
				return
			}
		}
		t.Logf("sleep %s", testInterval)
		time.Sleep(testInterval)
	}
}

func testAllMetrics(ctx context.Context, t *testing.T) {
	meter := otel.Meter(instrumentation, metric.WithSchemaURL(semconv.SchemaURL))

	// Create All instruments, they should not error
	aiCounter, err := meter.Int64ObservableCounter("observable.int64.counter")
	if err != nil {
		t.Fatalf("create meter Int64ObservableCounter failed: %s", err.Error())
	}
	aiUpDownCounter, err := meter.Int64ObservableUpDownCounter("observable.int64.up.down.counter")
	if err != nil {
		t.Fatalf("create meter Int64ObservableUpDownCounter failed: %s", err.Error())
	}
	aiGauge, err := meter.Int64ObservableGauge("observable.int64.gauge")
	if err != nil {
		t.Fatalf("create meter Int64ObservableGauge failed: %s", err.Error())
	}
	afCounter, err := meter.Float64ObservableCounter("observable.float64.counter")
	if err != nil {
		t.Fatalf("create meter Float64ObservableCounter failed: %s", err.Error())
	}
	afUpDownCounter, err := meter.Float64ObservableUpDownCounter("observable.float64.up.down.counter")
	if err != nil {
		t.Fatalf("create meter Float64ObservableUpDownCounter failed: %s", err.Error())
	}
	afGauge, err := meter.Float64ObservableGauge("observable.float64.gauge")
	if err != nil {
		t.Fatalf("create meter Float64ObservableGauge failed: %s", err.Error())
	}
	siCounter, err := meter.Int64Counter("sync.int64.counter")
	if err != nil {
		t.Fatalf("create meter Int64Counter failed: %s", err.Error())
	}
	siUpDownCounter, err := meter.Int64UpDownCounter("sync.int64.up.down.counter")
	if err != nil {
		t.Fatalf("create meter Int64UpDownCounter failed: %s", err.Error())
	}
	siHistogram, err := meter.Int64Histogram("sync.int64.histogram", metric.WithExplicitBucketBoundaries(0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10))
	if err != nil {
		t.Fatalf("create meter Int64Histogram failed: %s", err.Error())
	}
	sfCounter, err := meter.Float64Counter("sync.float64.counter")
	if err != nil {
		t.Fatalf("create meter Float64Counter failed: %s", err.Error())
	}
	sfUpDownCounter, err := meter.Float64UpDownCounter("sync.float64.up.down.counter")
	if err != nil {
		t.Fatalf("create meter Float64UpDownCounter failed: %s", err.Error())
	}
	sfHistogram, err := meter.Float64Histogram("sync.float64.histogram", metric.WithExplicitBucketBoundaries(0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10))
	if err != nil {
		t.Fatalf("create meter Float64Histogram failed: %s", err.Error())
	}

	callback := func(ctx context.Context, obs metric.Observer) error {
		obs.ObserveInt64(aiCounter, 1)
		obs.ObserveInt64(aiUpDownCounter, 1)
		obs.ObserveInt64(aiGauge, 1)
		obs.ObserveFloat64(afCounter, 1)
		obs.ObserveFloat64(afUpDownCounter, 1)
		obs.ObserveFloat64(afGauge, 1)
		return nil
	}
	_, err = meter.RegisterCallback(callback, aiCounter, aiUpDownCounter, aiGauge, afCounter, afUpDownCounter, afGauge)
	if err != nil {
		t.Fatalf("register callback failed: %s", err.Error())
	}

	siCounter.Add(context.Background(), 1)
	siUpDownCounter.Add(context.Background(), 1)
	siHistogram.Record(context.Background(), 1)
	sfCounter.Add(context.Background(), 1)
	sfUpDownCounter.Add(context.Background(), 1)
	sfHistogram.Record(context.Background(), 1)
}
