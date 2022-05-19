// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wrap

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/metric/export"
)

var _ export.CheckpointerFactory = (*Factory)(nil)
var _ export.Checkpointer = (*Processor)(nil)
var _ export.Processor = (*Processor)(nil)

type AccumulationHandlerFunc func(accum export.Accumulation) export.Accumulation

var AccumulationHandler []AccumulationHandlerFunc

func RegisterAccumulationHandler(handlers ...AccumulationHandlerFunc) {
	AccumulationHandler = append(AccumulationHandler, handlers...)
}

// Processor is a wrapped SpanProcessor.
type Processor struct {
	Handlers []AccumulationHandlerFunc
	export.Checkpointer
}

type Factory struct {
	opts []FactoryConfigFunc
	export.CheckpointerFactory
}

func (f *Factory) NewCheckpointer() export.Checkpointer {
	return New(f.CheckpointerFactory.NewCheckpointer(), f.opts...)
}

// New returns a dimensionality-reset Processor that passes data to
// the next stage in an export pipeline.
func New(ckpter export.Checkpointer, opts ...FactoryConfigFunc) *Processor {
	var config FactoryConfig
	err := config.ApplyOptions(opts...)
	if err != nil {
		otel.Handle(err)
	}
	return &Processor{
		Checkpointer: ckpter,
		Handlers:     config.Handlers,
	}
}

// NewFactory returns a dimensionality-reset Processor Factory that passes data to
// the next stage in an export pipeline.
func NewFactory(checkpointerFactory export.CheckpointerFactory, opts ...FactoryConfigFunc) *Factory {
	return &Factory{
		opts:                opts,
		CheckpointerFactory: checkpointerFactory,
	}
}

// Process implements export.Processor.
func (p *Processor) Process(accum export.Accumulation) error {
	for _, h := range AccumulationHandler {
		accum = h(accum)
	}
	for _, h := range p.Handlers {
		accum = h(accum)
	}
	return p.Checkpointer.Process(accum)
}
