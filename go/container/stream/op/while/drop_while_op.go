// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package while

import (
	"context"

	"github.com/searKing/golang/go/container/stream/op/terminal"
	"github.com/searKing/golang/go/util/function/predicate"
	"github.com/searKing/golang/go/util/object"
	"github.com/searKing/golang/go/util/optional"
	"github.com/searKing/golang/go/util/spliterator"
)

/**
* A specialization for the dropWhile operation that controls if
* elements to be dropped are counted and passed downstream.
* <p>
* This specialization is utilized by the {@link TakeWhileTask} for
* pipelines that are ordered.  In such cases elements cannot be dropped
* until all elements have been collected.
*
* @param <T> the type of both input and output elements
 */
type DropWhileOp struct {
	terminal.TODOOperation
	sinkNewer func() DropWhileSink
}

func NewDropWhileOp(sinkNewer func() DropWhileSink) *DropWhileOp {
	reduceOp := &DropWhileOp{
		sinkNewer: sinkNewer,
	}
	return reduceOp
}

func NewDropWhileOp2(retainAndCountDroppedElements bool, predicate predicate.Predicater) *DropWhileOp {
	object.RequireNonNil(predicate)
	return NewDropWhileOp(func() DropWhileSink {
		return NewDropWhileSink(retainAndCountDroppedElements, predicate)
	})
}

func (op DropWhileOp) MakeSink() DropWhileSink {
	object.RequireNonNull(op.sinkNewer)
	return op.sinkNewer()
}

func (op DropWhileOp) EvaluateParallel(ctx context.Context, spliterator spliterator.Spliterator) optional.Optional {
	return terminal.WrapAndCopyInto(ctx, op.MakeSink(), spliterator).Get()
}

func (op DropWhileOp) EvaluateSequential(ctx context.Context, spliterator spliterator.Spliterator) optional.Optional {
	// Approach for parallel implementation:
	// - Decompose as per usual
	// - run match on leaf chunks, call result "b"
	// - if b == matchKind.shortCircuitOn, complete early and return b
	// - else if we complete normally, return !shortCircuitOn
	return (&DropWhileTask{}).WithSpliterator(op, spliterator).Invoke(ctx).Get()
}
