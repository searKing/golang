// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wrap

import (
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/export"
)

type FactoryConfig struct {
	// Handlers acts like a middleware of Processor.
	Handlers []AccumulationHandlerFunc
}

// FactoryConfigFunc is an alias for a function that will take in a pointer to an FactoryConfig and modify it
type FactoryConfigFunc func(os *FactoryConfig) error

func (fc *FactoryConfig) ApplyOptions(configs ...FactoryConfigFunc) error {
	// Apply all Configurations passed in
	for _, config := range configs {
		// Pass the FactoryConfig into the configuration function
		err := config(fc)
		if err != nil {
			return fmt.Errorf("failed to apply configuration function: %w", err)
		}
	}
	return nil
}

// WithDefaultLabels add default label set to Accumulation, as prepared by an Accumulator for the Processor.
// WithDefaultLabels set default attributes to all started spans.
// Duplicate keys are eliminated by taking the last value.
func WithDefaultLabels(kvs ...attribute.KeyValue) func(os *FactoryConfig) error {
	return func(os *FactoryConfig) error {
		os.Handlers = append(os.Handlers, func(accum export.Accumulation) export.Accumulation {
			if len(kvs) == 0 {
				return accum
			}
			reduced := attribute.NewSet(append(kvs, accum.Attributes().ToSlice()...)...)
			return export.NewAccumulation(
				accum.Descriptor(),
				&reduced,
				accum.Aggregator(),
			)
		})
		return nil
	}
}

// WithAccumulationHandler sets the middleware of Processor.
func WithAccumulationHandler(handlers ...AccumulationHandlerFunc) func(os *FactoryConfig) error {
	return func(os *FactoryConfig) error {
		os.Handlers = append(os.Handlers, handlers...)
		return nil
	}
}
