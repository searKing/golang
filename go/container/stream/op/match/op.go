// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package match

import (
	"context"

	"github.com/searKing/golang/go/container/stream/op/terminal"
	"github.com/searKing/golang/go/util/function/predicate"
	"github.com/searKing/golang/go/util/object"
	"github.com/searKing/golang/go/util/optional"
	"github.com/searKing/golang/go/util/spliterator"
)

type MatchOp struct {
	terminal.TODOOperation
	matchKind Kind
	sinkNewer func() *MatchSink
}

func NewMatchOp(matchKind Kind, sinkNewer func() *MatchSink) *MatchOp {
	reduceOp := &MatchOp{
		matchKind: matchKind,
		sinkNewer: sinkNewer,
	}
	reduceOp.SetDerived(reduceOp)
	return reduceOp
}

func NewMatchOp2(matchKind Kind, predicate predicate.Predicater) *MatchOp {
	object.RequireNonNil(predicate)
	return NewMatchOp(matchKind, func() *MatchSink {
		return NewMatchSink(matchKind, predicate)
	})
}

func (op MatchOp) MakeSink() *MatchSink {
	object.RequireNonNull(op.sinkNewer)
	return op.sinkNewer()
}

func (op MatchOp) EvaluateParallel(ctx context.Context, spliterator spliterator.Spliterator) optional.Optional {
	return terminal.WrapAndCopyInto(ctx, op.MakeSink(), spliterator).Get()
}

func (op MatchOp) EvaluateSequential(ctx context.Context, spliterator spliterator.Spliterator) optional.Optional {
	// Approach for parallel implementation:
	// - Decompose as per usual
	// - run match on leaf chunks, call result "b"
	// - if b == matchKind.shortCircuitOn, complete early and return b
	// - else if we complete normally, return !shortCircuitOn
	return NewMatchTask(op, spliterator).Invoke(ctx).Get()
}
