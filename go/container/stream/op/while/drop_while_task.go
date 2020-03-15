// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package while

import (
	"context"

	"github.com/searKing/golang/go/container/stream/op/terminal"
	"github.com/searKing/golang/go/util/spliterator"
)

type DropWhileTask struct {
	terminal.TODOShortCircuitTask

	op DropWhileOp

	isOrdered    bool
	thisNodeSize int
	// The index from which elements of the node should be taken
	// i.e. the node should be truncated from [takeIndex, thisNodeSize)
	// Equivalent to the count of dropped elements
	index int
}

/**
 * Constructor for root tasks.
 *
 * @param helper the {@code PipelineHelper} describing the stream pipeline
 *               up to this operation
 * @param spliterator the {@code Spliterator} describing the source for this
 *                    pipeline
 */
func (task *DropWhileTask) WithSpliterator(op DropWhileOp, spliterator spliterator.Spliterator) *DropWhileTask {
	task.TODOTask.WithSpliterator(spliterator)
	task.op = op
	//task.isOrdered = op
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
func (task *DropWhileTask) WithParent(parent *DropWhileTask, spliterator spliterator.Spliterator) *DropWhileTask {
	task.TODOTask.WithParent(parent, spliterator)
	task.op = parent.op
	task.isOrdered = parent.isOrdered
	task.SetDerived(task)
	return task
}

func (task *DropWhileTask) MakeChild(spliterator spliterator.Spliterator) terminal.Task {
	child := &DropWhileTask{}
	return child.WithParent(task, spliterator)
}

func (task *DropWhileTask) foundResult(answer terminal.Sink) {
	if task.IsLeftmostNode() {
		task.ShortCircuit(answer)
		return
	}
	task.CancelLaterNodes()
}

func (task *DropWhileTask) DoLeaf(ctx context.Context) terminal.Sink {
	result := terminal.WrapAndCopyInto(ctx, task.op.MakeSink(), task.GetSpliterator())
	if !task.isOrdered {
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

func (task *DropWhileTask) OnCompletion(caller terminal.Task) {
	if task.isOrdered {
		child := task.LeftChild()
		var p terminal.Task
		for child != p {
			result := child.GetLocalResult()
			if result != nil && result.Get().IsPresent() {
				task.SetLocalResult(result)
				task.foundResult(result)
				break
			}

			p = child
			child = task.RightChild()
		}
	}
	// GC spliterator, left and right child
	task.TODOTask.OnCompletion(caller)
}
