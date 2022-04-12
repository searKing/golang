// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package simple_test

import (
	"testing"

	"github.com/searKing/golang/third_party/github.com/open-telemetry/opentelemetry-go/metric/sdk/metric/aggregator/exact"
	"github.com/searKing/golang/third_party/github.com/open-telemetry/opentelemetry-go/metric/sdk/metric/selector/simple"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/sdk/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/lastvalue"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
	"go.opentelemetry.io/otel/sdk/metric/export"
	"go.opentelemetry.io/otel/sdk/metric/metrictest"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/sdkapi"
)

var (
	testCounterDesc               = metrictest.NewDescriptor("counter", sdkapi.CounterInstrumentKind, number.Int64Kind)
	testUpDownCounterDesc         = metrictest.NewDescriptor("updowncounter", sdkapi.UpDownCounterInstrumentKind, number.Int64Kind)
	testCounterObserverDesc       = metrictest.NewDescriptor("counterobserver", sdkapi.CounterObserverInstrumentKind, number.Int64Kind)
	testUpDownCounterObserverDesc = metrictest.NewDescriptor("updowncounterobserver", sdkapi.UpDownCounterObserverInstrumentKind, number.Int64Kind)
	testHistogramDesc             = metrictest.NewDescriptor("histogram", sdkapi.HistogramInstrumentKind, number.Int64Kind)
	testGaugeObserverDesc         = metrictest.NewDescriptor("gauge", sdkapi.GaugeObserverInstrumentKind, number.Int64Kind)
)

func oneAgg(sel export.AggregatorSelector, desc *sdkapi.Descriptor) aggregator.Aggregator {
	var agg aggregator.Aggregator
	sel.AggregatorFor(desc, &agg)
	return agg
}

func testFixedSelectors(t *testing.T, sel export.AggregatorSelector) {
	require.IsType(t, (*lastvalue.Aggregator)(nil), oneAgg(sel, &testGaugeObserverDesc))
	require.IsType(t, (*sum.Aggregator)(nil), oneAgg(sel, &testCounterDesc))
	require.IsType(t, (*sum.Aggregator)(nil), oneAgg(sel, &testUpDownCounterDesc))
	require.IsType(t, (*sum.Aggregator)(nil), oneAgg(sel, &testCounterObserverDesc))
	require.IsType(t, (*sum.Aggregator)(nil), oneAgg(sel, &testUpDownCounterObserverDesc))
}

func TestExactDistribution(t *testing.T) {
	ex := simple.NewWithExactDistribution()
	require.IsType(t, (*exact.Aggregator)(nil), oneAgg(ex, &testHistogramDesc))
	testFixedSelectors(t, ex)
}
