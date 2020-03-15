// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package find

import (
	"context"

	"github.com/searKing/golang/go/container/stream/op/terminal"
	"github.com/searKing/golang/go/util/spliterator"
)

type FindTask struct {
	terminal.TODOShortCircuitTask

	op FindOp
	// true if find first
	// false if find any
	mustFindFirst bool
}

/**
 * Constructor for root tasks.
 *
 * @param helper the {@code PipelineHelper} describing the stream pipeline
 *               up to this operation
 * @param spliterator the {@code Spliterator} describing the source for this
 *                    pipeline
 */
func (task *FindTask) WithSpliterator(op FindOp, spliterator spliterator.Spliterator) *FindTask {
	task.TODOShortCircuitTask.WithSpliterator(spliterator)
	task.op = op
	task.mustFindFirst = op.mustFindFirst
	task.SetDerived(task)
	return task
}

/**
 * Constructor for non-root nodes.
 *
 * @param parent parent task in the computation tree
 * @param spliterator the {@code Spliterator} for the portion of the
 *                    computation tree described by this task
 */
func (task *FindTask) WithParent(parent *FindTask, spliterator spliterator.Spliterator) *FindTask {
	task.TODOShortCircuitTask.WithParent(parent, spliterator)
	task.op = parent.op
	task.mustFindFirst = parent.mustFindFirst
	task.SetDerived(task)
	return task
}

func (task *FindTask) MakeChild(spliterator spliterator.Spliterator) terminal.Task {
	child := &FindTask{}
	return child.WithParent(task, spliterator)
}

func (task *FindTask) foundResult(answer terminal.Sink) {
	if task.IsLeftmostNode() {
		task.ShortCircuit(answer)
		return
	}
	task.CancelLaterNodes()
}

func (task *FindTask) DoLeaf(ctx context.Context) terminal.Sink {
	result := terminal.WrapAndCopyInto(ctx, task.op.MakeSink(), task.GetSpliterator())
	if !task.mustFindFirst {
		if result != nil && result.Get().IsPresent() {
			task.ShortCircuit(result)
		}
		return nil
	}
	if result != nil && result.Get().IsPresent() {
		task.foundResult(result)
		return result
	}
	return nil
}

func (task *FindTask) OnCompletion(caller terminal.Task) {
	if task.mustFindFirst {
		child := task.LeftChild()
		var p terminal.Task
		for child != p {
			result := child.GetLocalResult()
			if result != nil && task.op.presentPredicate.Test(result.Get().Get()) {
				task.SetLocalResult(result)
				task.foundResult(result)
				break
			}

			p = child
			child = task.RightChild()
		}
	}
	// GC spliterator, left and right child
	task.TODOShortCircuitTask.OnCompletion(caller)
}
