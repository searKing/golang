// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reduce

import (
	"context"

	"github.com/searKing/golang/go/container/stream/op/terminal"
	"github.com/searKing/golang/go/util/function/binary"
	"github.com/searKing/golang/go/util/object"
	"github.com/searKing/golang/go/util/optional"
	"github.com/searKing/golang/go/util/spliterator"
)

type ReduceOp struct {
	terminal.TODOOperation

	sinkNewer func() *ReducingSink
}

func NewReduceOp(sinkNewer func() *ReducingSink) *ReduceOp {
	reduceOp := &ReduceOp{
		sinkNewer: sinkNewer,
	}
	reduceOp.SetDerived(reduceOp)
	return reduceOp
}

func NewReduceOp3(seed optional.Optional, reducer binary.BiFunction, combiner binary.BiFunction) terminal.Operation {
	object.RequireNonNil(reducer)
	return NewReduceOp(func() *ReducingSink {
		return NewReducingSink(seed, reducer, combiner)
	})
}

func (op ReduceOp) MakeSink() *ReducingSink {
	object.RequireNonNull(op.sinkNewer)
	return op.sinkNewer()
}

func (op ReduceOp) EvaluateParallel(ctx context.Context, spliterator spliterator.Spliterator) optional.Optional {
	return terminal.WrapAndCopyInto(ctx, op.MakeSink(), spliterator).Get()
}

func (op ReduceOp) EvaluateSequential(ctx context.Context, spliterator spliterator.Spliterator) optional.Optional {
	return NewReduceTask(op, spliterator).Invoke(ctx).Get()
}
