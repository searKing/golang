// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wrap

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/metric/export"
)

var _ export.CheckpointerFactory = (*Wrap)(nil)
var _ export.Checkpointer = (*Wrap)(nil)

type AccumulationHandlerFunc func(accum export.Accumulation) export.Accumulation

var AccumulationHandler []AccumulationHandlerFunc

func RegisterAccumulationHandler(handlers ...AccumulationHandlerFunc) {
	AccumulationHandler = append(AccumulationHandler, handlers...)
}

// Wrap is a wrapped SpanProcessor.
type Wrap struct {
	Handlers []AccumulationHandlerFunc
	export.Checkpointer
}

func (p *Wrap) NewCheckpointer() export.Checkpointer {
	return p
}

var _ export.Processor = &Wrap{}
var _ export.Checkpointer = &Wrap{}

// New returns a dimensionality-reset Processor that passes data to
// the next stage in an export pipeline.
func New(ckpter export.Checkpointer, opts ...FactoryConfigFunc) *Wrap {
	var config FactoryConfig
	err := config.ApplyOptions(opts...)
	if err != nil {
		otel.Handle(err)
	}
	return &Wrap{
		Checkpointer: ckpter,
		Handlers:     config.Handlers,
	}
}

// Process implements export.Processor.
func (p *Wrap) Process(accum export.Accumulation) error {
	for _, h := range AccumulationHandler {
		accum = h(accum)
	}
	for _, h := range p.Handlers {
		accum = h(accum)
	}
	return p.Checkpointer.Process(accum)
}
