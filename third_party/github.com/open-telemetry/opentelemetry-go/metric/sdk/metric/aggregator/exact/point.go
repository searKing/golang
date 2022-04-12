// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exact

import (
	"time"

	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/number"
)

type (
	// Points returns the raw values that were aggregated.
	Points interface {
		aggregation.Aggregation

		// Points returns points in the order they were
		// recorded.  Points are approximately ordered by
		// timestamp, but this is not guaranteed.
		Points() ([]Point, error)
	}

	// Point is a raw data point, consisting of a number and value.
	Point struct {
		number.Number
		time.Time
	}
)
