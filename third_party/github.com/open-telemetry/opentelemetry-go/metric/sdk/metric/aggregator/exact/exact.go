// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exact

import (
	"context"
	"sync"
	"time"

	"go.opentelemetry.io/otel/sdk/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/sdkapi"
)

type (
	// Aggregator aggregates events that form a distribution, keeping
	// an array with the exact set of values.
	Aggregator struct {
		lock    sync.Mutex
		samples []Point
	}
)

var _ aggregator.Aggregator = &Aggregator{}
var _ Points = &Aggregator{}
var _ aggregation.Count = &Aggregator{}

// New returns cnt many new exact aggregators, which aggregate recorded
// measurements by storing them in an array.  This type uses a mutex
// for Update() and SynchronizedMove() concurrency.
func New(cnt int) []Aggregator {
	return make([]Aggregator, cnt)
}

// Aggregation returns an interface for reading the state of this aggregator.
func (c *Aggregator) Aggregation() aggregation.Aggregation {
	return c
}

// Kind returns aggregation.ExactKind.
func (c *Aggregator) Kind() aggregation.Kind {
	return aggregation.HistogramKind
}

// Count returns the number of values in the checkpoint.
func (c *Aggregator) Count() (uint64, error) {
	return uint64(len(c.samples)), nil
}

// Points returns access to the raw data set.
func (c *Aggregator) Points() ([]Point, error) {
	return c.samples, nil
}

// SynchronizedMove saves the current state to oa and resets the current state to
// the empty set, taking a lock to prevent concurrent Update() calls.
func (c *Aggregator) SynchronizedMove(oa aggregator.Aggregator, desc *sdkapi.Descriptor) error {
	o, _ := oa.(*Aggregator)

	if oa != nil && o == nil {
		return aggregator.NewInconsistentAggregatorError(c, oa)
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	if o != nil {
		o.samples = c.samples
	}
	c.samples = nil

	return nil
}

// Update adds the recorded measurement to the current data set.
// Update takes a lock to prevent concurrent Update() and SynchronizedMove()
// calls.
func (c *Aggregator) Update(_ context.Context, number number.Number, desc *sdkapi.Descriptor) error {
	now := time.Now()
	c.lock.Lock()
	defer c.lock.Unlock()
	c.samples = append(c.samples, Point{
		Number: number,
		Time:   now,
	})

	return nil
}

// Merge combines two data sets into one.
func (c *Aggregator) Merge(oa aggregator.Aggregator, desc *sdkapi.Descriptor) error {
	o, _ := oa.(*Aggregator)
	if o == nil {
		return aggregator.NewInconsistentAggregatorError(c, oa)
	}

	c.samples = combine(c.samples, o.samples)
	return nil
}

func combine(a, b []Point) []Point {
	result := make([]Point, 0, len(a)+len(b))

	for len(a) != 0 && len(b) != 0 {
		if a[0].Time.Before(b[0].Time) {
			result = append(result, a[0])
			a = a[1:]
		} else {
			result = append(result, b[0])
			b = b[1:]
		}
	}
	result = append(result, a...)
	result = append(result, b...)
	return result
}
