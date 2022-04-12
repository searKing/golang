// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package simple

import (
	"github.com/searKing/golang/third_party/github.com/open-telemetry/opentelemetry-go/metric/sdk/metric/aggregator/exact"
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/lastvalue"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
	"go.opentelemetry.io/otel/sdk/metric/export"
	"go.opentelemetry.io/otel/sdk/metric/sdkapi"
)

type (
	selectorExact struct{}
)

var (
	_ export.AggregatorSelector = selectorExact{}
)

// NewWithExactDistribution returns a simple aggregator selector that
// uses exact aggregators for `Histogram` instruments.  This
// selector uses more memory than the others in this package because
// exact aggregators maintain the most information about the
// distribution among these choices.
func NewWithExactDistribution() export.AggregatorSelector {
	return selectorExact{}
}

func sumAggs(aggPtrs []*aggregator.Aggregator) {
	aggs := sum.New(len(aggPtrs))
	for i := range aggPtrs {
		*aggPtrs[i] = &aggs[i]
	}
}

func lastValueAggs(aggPtrs []*aggregator.Aggregator) {
	aggs := lastvalue.New(len(aggPtrs))
	for i := range aggPtrs {
		*aggPtrs[i] = &aggs[i]
	}
}

func (selectorExact) AggregatorFor(descriptor *sdkapi.Descriptor, aggPtrs ...*aggregator.Aggregator) {
	switch descriptor.InstrumentKind() {
	case sdkapi.GaugeObserverInstrumentKind:
		lastValueAggs(aggPtrs)
	case sdkapi.HistogramInstrumentKind:
		aggs := exact.New(len(aggPtrs))
		for i := range aggPtrs {
			*aggPtrs[i] = &aggs[i]
		}
	default:
		sumAggs(aggPtrs)
	}
}
