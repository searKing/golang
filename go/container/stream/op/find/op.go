// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package find

import (
	"context"

	"github.com/searKing/golang/go/container/stream/op/terminal"
	"github.com/searKing/golang/go/util/function/predicate"
	"github.com/searKing/golang/go/util/object"
	"github.com/searKing/golang/go/util/optional"
	"github.com/searKing/golang/go/util/spliterator"
)

/**
 * A short-circuiting {@code TerminalOp} that searches for an element in a
 * stream pipeline, and terminates when it finds one.  Implements both
 * find-first (find the first element in the encounter order) and find-any
 * (find any element, may not be the first in encounter order.)
 *
 * @param <T> the output type of the stream pipeline
 * @param <O> the result type of the find operation, typically an optional
 *        type
 */
type FindOp struct {
	terminal.TODOOperation
	mustFindFirst bool
	sinkSupplier  func() *FindSink

	presentPredicate predicate.Predicater
}

/**
 * Constructs a {@code Operation}.
 *
 * @param mustFindFirst if true, must find the first element in
 *        encounter order, otherwise can find any element
 * @param shape stream shape of elements to search
 * @param emptyValue result value corresponding to "found nothing"
 * @param presentPredicate {@code Predicate} on result value
 *        corresponding to "found something"
 * @param sinkSupplier supplier for a {@code TerminalSink} implementing
 *        the matching functionality
 */
func NewFindOp(mustFindFirst bool, presentPredicate predicate.Predicater, sinkSupplier func() *FindSink) *FindOp {
	findOp := &FindOp{
		mustFindFirst:    mustFindFirst,
		presentPredicate: presentPredicate,
		sinkSupplier:     sinkSupplier,
	}
	findOp.SetDerived(findOp)
	return findOp
}

func NewFindOp2(mustFindFirst bool, presentPredicate predicate.Predicater) *FindOp {
	object.RequireNonNil(presentPredicate)
	return NewFindOp(mustFindFirst, presentPredicate, func() *FindSink {
		return NewFindSink()
	})
}

func (op FindOp) MakeSink() *FindSink {
	object.RequireNonNull(op.sinkSupplier)
	return op.sinkSupplier()
}

func (op FindOp) EvaluateParallel(ctx context.Context, spliterator spliterator.Spliterator) optional.Optional {
	return terminal.WrapAndCopyInto(ctx, op.MakeSink(), spliterator).Get()
}

func (op FindOp) EvaluateSequential(ctx context.Context, spliterator spliterator.Spliterator) optional.Optional {
	return (&FindTask{}).WithSpliterator(op, spliterator).Invoke(ctx).Get()
}
