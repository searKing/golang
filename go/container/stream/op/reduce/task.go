// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reduce

import (
	"context"

	"github.com/searKing/golang/go/container/stream/op/terminal"
	"github.com/searKing/golang/go/util/function/binary"
	"github.com/searKing/golang/go/util/spliterator"
)

type ReduceTask struct {
	terminal.TODOTask

	op      ReduceOp
	combine binary.BiFunction
}

func NewReduceTask(op ReduceOp, spliterator spliterator.Spliterator) *ReduceTask {
	reducer := &ReduceTask{
		op: op,
	}
	reducer.WithSpliterator(spliterator)
	reducer.SetDerived(reducer)
	return reducer
}

func NewReduceTaskFromParent(parent *ReduceTask, spliterator spliterator.Spliterator) *ReduceTask {
	reducer := &ReduceTask{
		op: parent.op,
	}
	reducer.WithParent(parent, spliterator)
	reducer.SetDerived(reducer)
	return reducer
}

func (t *ReduceTask) MakeChild(spliterator spliterator.Spliterator) terminal.Task {
	return NewReduceTaskFromParent(t, spliterator)
}

func (t *ReduceTask) DoLeaf(ctx context.Context) terminal.Sink {
	return terminal.WrapAndCopyInto(ctx, t.op.MakeSink(), t.GetSpliterator())
}

func (t *ReduceTask) OnCompletion(caller terminal.Task) {
	task := t.GetDerivedElse(t).(terminal.Task)
	if !task.IsLeaf() {
		leftResult := task.LeftChild().GetLocalResult()
		rightResult := task.RightChild().GetLocalResult()
		sink := t.op.MakeSink()
		sink.Begin(-1)
		sink.Combine(leftResult)
		sink.Combine(rightResult)
		sink.End()
		task.SetLocalResult(sink)
	}
	// GC spliterator, left and right child
	t.TODOTask.OnCompletion(caller)
}
