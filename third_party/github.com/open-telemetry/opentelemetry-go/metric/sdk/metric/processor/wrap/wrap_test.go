// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wrap_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/searKing/golang/third_party/github.com/open-telemetry/opentelemetry-go/metric/sdk/metric/processor/wrap"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/export"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/metrictest"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/processor/basic"
	processorTest "go.opentelemetry.io/otel/sdk/metric/processor/processortest"
	"go.opentelemetry.io/otel/sdk/metric/sdkapi"
)

// TestProcessor tests all the non-error paths in this package.
func TestProcessor(t *testing.T) {
	type exportCase struct {
		kind aggregation.Temporality
	}
	type instrumentCase struct {
		kind sdkapi.InstrumentKind
	}
	type numberCase struct {
		kind number.Kind
	}
	type aggregatorCase struct {
		kind aggregation.Kind
	}

	for _, tc := range []exportCase{
		{kind: aggregation.CumulativeTemporality},
		{kind: aggregation.DeltaTemporality},
	} {
		t.Run(tc.kind.String(), func(t *testing.T) {
			for _, ic := range []instrumentCase{
				{kind: sdkapi.CounterInstrumentKind},
				{kind: sdkapi.UpDownCounterInstrumentKind},
				{kind: sdkapi.HistogramInstrumentKind},
				{kind: sdkapi.CounterObserverInstrumentKind},
				{kind: sdkapi.UpDownCounterObserverInstrumentKind},
				{kind: sdkapi.GaugeObserverInstrumentKind},
			} {
				t.Run(ic.kind.String(), func(t *testing.T) {
					for _, nc := range []numberCase{
						{kind: number.Int64Kind},
						{kind: number.Float64Kind},
					} {
						t.Run(nc.kind.String(), func(t *testing.T) {
							for _, ac := range []aggregatorCase{
								{kind: aggregation.SumKind},
								{kind: aggregation.HistogramKind},
								{kind: aggregation.LastValueKind},
							} {
								t.Run(ac.kind.String(), func(t *testing.T) {
									testProcessor(
										t,
										tc.kind,
										ic.kind,
										nc.kind,
										ac.kind,
									)
								})
							}
						})
					}
				})
			}
		})
	}
}

func asNumber(nkind number.Kind, value int64) number.Number {
	if nkind == number.Int64Kind {
		return number.NewInt64Number(value)
	}
	return number.NewFloat64Number(float64(value))
}

func updateFor(t *testing.T, desc *sdkapi.Descriptor, selector export.AggregatorSelector, value int64, labs ...attribute.KeyValue) export.Accumulation {
	ls := attribute.NewSet(labs...)
	var agg aggregator.Aggregator
	selector.AggregatorFor(desc, &agg)
	require.NoError(t, agg.Update(context.Background(), asNumber(desc.NumberKind(), value), desc))

	return export.NewAccumulation(desc, &ls, agg)
}

func testProcessor(
	t *testing.T,
	aggTemp aggregation.Temporality,
	mkind sdkapi.InstrumentKind,
	nkind number.Kind,
	akind aggregation.Kind,
) {
	// This code tests for errors when the export kind is Delta
	// and the instrument kind is PrecomputedSum().
	expectConversion := !(aggTemp == aggregation.DeltaTemporality && mkind.PrecomputedSum())
	requireConversion := func(t *testing.T, err error) {
		if expectConversion {
			require.NoError(t, err)
		} else {
			require.Equal(t, aggregation.ErrNoCumulativeToDelta, err)
		}
	}

	// Note: this selector uses the instrument name to dictate
	// aggregation kind.
	selector := processorTest.AggregatorSelector()

	labs1 := []attribute.KeyValue{attribute.String("L1", "V")}
	labs2 := []attribute.KeyValue{attribute.String("L2", "V")}

	testBody := func(t *testing.T, hasMemory bool, nAccum, nCheckpoint int) {
		processor := wrap.New(
			basic.New(selector, aggregation.ConstantTemporalitySelector(aggTemp), basic.WithMemory(hasMemory)),
			wrap.WithDefaultLabels(attribute.String("L3", "V"), attribute.String("L2", "M")))

		instSuffix := fmt.Sprint(".", strings.ToLower(akind.String()))

		desc1 := metrictest.NewDescriptor(fmt.Sprint("inst1", instSuffix), mkind, nkind)
		desc2 := metrictest.NewDescriptor(fmt.Sprint("inst2", instSuffix), mkind, nkind)

		for nc := 0; nc < nCheckpoint; nc++ {

			// The input is 10 per update, scaled by
			// the number of checkpoints for
			// cumulative instruments:
			input := int64(10)
			cumulativeMultiplier := int64(nc + 1)
			if mkind.PrecomputedSum() {
				input *= cumulativeMultiplier
			}

			processor.StartCollection()

			for na := 0; na < nAccum; na++ {
				requireConversion(t, processor.Process(updateFor(t, &desc1, selector, input, labs1...)))
				requireConversion(t, processor.Process(updateFor(t, &desc2, selector, input, labs2...)))
			}

			// Note: in case of !expectConversion, we still get no error here
			// because the Process() skipped entering state for those records.
			require.NoError(t, processor.FinishCollection())

			if nc < nCheckpoint-1 {
				continue
			}

			reader := processor.Reader()

			for _, repetitionAfterEmptyInterval := range []bool{false, true} {
				if repetitionAfterEmptyInterval {
					// We're repeating the test after another
					// interval with no updates.
					processor.StartCollection()
					require.NoError(t, processor.FinishCollection())
				}

				// Test the final checkpoint state.
				records1 := processorTest.NewOutput(attribute.DefaultEncoder())
				require.NoError(t, reader.ForEach(aggregation.ConstantTemporalitySelector(aggTemp), records1.AddRecord))

				if !expectConversion {
					require.EqualValues(t, map[string]float64{}, records1.Map())
					continue
				}

				var multiplier int64

				if mkind.Asynchronous() {
					// Asynchronous tests accumulate results multiply by the
					// number of Accumulators, unless LastValue aggregation.
					// If a precomputed sum, we expect cumulative inputs.
					if mkind.PrecomputedSum() {
						require.NotEqual(t, aggTemp, aggregation.DeltaTemporality)
						if akind == aggregation.LastValueKind {
							multiplier = cumulativeMultiplier
						} else {
							multiplier = cumulativeMultiplier * int64(nAccum)
						}
					} else {
						if aggTemp == aggregation.CumulativeTemporality && akind != aggregation.LastValueKind {
							multiplier = cumulativeMultiplier * int64(nAccum)
						} else if akind == aggregation.LastValueKind {
							multiplier = 1
						} else {
							multiplier = int64(nAccum)
						}
					}
				} else {
					// Synchronous accumulate results from multiple accumulators,
					// use that number as the baseline multiplier.
					multiplier = int64(nAccum)
					if aggTemp == aggregation.CumulativeTemporality {
						// If a cumulative exporter, include prior checkpoints.
						multiplier *= cumulativeMultiplier
					}
					if akind == aggregation.LastValueKind {
						// If a last-value aggregator, set multiplier to 1.0.
						multiplier = 1
					}
				}

				exp := map[string]float64{}
				if hasMemory || !repetitionAfterEmptyInterval {
					exp = map[string]float64{
						fmt.Sprintf("inst1%s/L1=V,L2=M,L3=V/", instSuffix): float64(multiplier * 10), // labels1
						fmt.Sprintf("inst2%s/L2=V,L3=V/", instSuffix):      float64(multiplier * 10), // labels2
					}
				}

				require.EqualValues(t, exp, records1.Map(), "with repetition=%v", repetitionAfterEmptyInterval)
			}
		}
	}

	for _, hasMem := range []bool{false, true} {
		t.Run(fmt.Sprintf("HasMemory=%v", hasMem), func(t *testing.T) {
			// For 1 to 3 checkpoints:
			for nAccum := 1; nAccum <= 3; nAccum++ {
				t.Run(fmt.Sprintf("NumAccum=%d", nAccum), func(t *testing.T) {
					// For 1 to 3 accumulators:
					for nCheckpoint := 1; nCheckpoint <= 3; nCheckpoint++ {
						t.Run(fmt.Sprintf("NumCkpt=%d", nCheckpoint), func(t *testing.T) {
							testBody(t, hasMem, nAccum, nCheckpoint)
						})
					}
				})
			}
		})
	}
}
