// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package processor

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"
)

// Annotator is a SpanProcessor that adds attributes to all started spans.
type Annotator struct {
	// AttrsFunc is called when a span is started. The attributes it returns
	// are set on the Span being started.
	AttrsFunc func() []attribute.KeyValue
}

// OnStart ...
func (a Annotator) OnStart(_ context.Context, s trace.ReadWriteSpan) {
	if a.AttrsFunc == nil {
		return
	}
	s.SetAttributes(a.AttrsFunc()...)
}

// Shutdown ...
func (a Annotator) Shutdown(context.Context) error { return nil }

// ForceFlush ...
func (a Annotator) ForceFlush(context.Context) error { return nil }

// OnEnd ...
func (a Annotator) OnEnd(s trace.ReadOnlySpan) {}
