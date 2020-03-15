// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package terminal

import (
	"context"

	"github.com/searKing/golang/go/util/object"
	"github.com/searKing/golang/go/util/spliterator"
)

/**
 * Applies the pipeline stages described by this {@code PipelineHelper} to
 * the provided {@code Spliterator} and send the results to the provided
 * {@code Sink}.
 *
 * @implSpec
 * The implementation behaves as if:
 * <pre>{@code
 *     copyInto(wrapSink(sink), spliterator);
 * }</pre>
 *
 * @param sink the {@code Sink} to receive the results
 * @param spliterator the spliterator describing the source input to process
 */
func WrapAndCopyInto(ctx context.Context, sink Sink, spliterator spliterator.Spliterator) Sink {
	CopyInto(ctx, object.RequireNonNull(sink).(Sink), spliterator)
	return sink
}

func CopyInto(ctx context.Context, wrappedSink Sink, spliterator spliterator.Spliterator) {
	object.RequireNonNil(wrappedSink)
	wrappedSink.Begin(spliterator.GetExactSizeIfKnown())
	spliterator.ForEachRemaining(ctx, wrappedSink)
	wrappedSink.End()
}
