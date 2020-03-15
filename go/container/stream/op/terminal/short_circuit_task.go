// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package terminal

import (
	"context"
	"sync/atomic"

	"github.com/searKing/golang/go/error/exception"
	"github.com/searKing/golang/go/util/spliterator"
)

//go:generate go-atomicvalue -type "sharedResult<Sink>"
type sharedResult atomic.Value

/**
 * Abstract class for fork-join tasks used to implement short-circuiting
 * stream ops, which can produce a result without processing all elements of the
 * stream.
 *
 * @param <P_IN> type of input elements to the pipeline
 * @param <P_OUT> type of output elements from the pipeline
 * @param <R> type of intermediate result, may be different from operation
 *        result type
 * @param <K> type of child and sibling tasks
 * @since 1.8
 */
type ShortCircuitTask interface {
	Task

	SharedResult() *sharedResult
	/**
	 * The result for this computation; this is shared among all tasks and set
	 * exactly once
	 */
	GetSharedResult() Sink

	/**
	 * Declares that a globally valid result has been found.  If another task has
	 * not already found the answer, the result is installed in
	 * {@code sharedResult}.  The {@code compute()} method will check
	 * {@code sharedResult} before proceeding with computation, so this causes
	 * the computation to terminate early.
	 *
	 * @param result the result found
	 */
	ShortCircuit(result Sink)

	/**
	 * Mark this task as canceled
	 */
	Cancel()

	/**
	 * Queries whether this task is canceled.  A task is considered canceled if
	 * it or any of its parents have been canceled.
	 *
	 * @return {@code true} if this task or any parent is canceled.
	 */
	TaskCanceled() bool

	/**
	 * Cancels all tasks which succeed this one in the encounter order.  This
	 * includes canceling all the current task's right sibling, as well as the
	 * later right siblings of all its parents.
	 */
	CancelLaterNodes()
}

type TODOShortCircuitTask struct {
	TODOTask

	/**
	 * The result for this computation; this is shared among all tasks and set
	 * exactly once
	 */
	sharedResult *sharedResult

	/**
	 * Indicates whether this task has been canceled.  Tasks may cancel other
	 * tasks in the computation under various conditions, such as in a
	 * find-first operation, a task that finds a value will cancel all tasks
	 * that are later in the encounter order.
	 */
	canceled bool
}

/**
 * Constructor for root tasks.
 *
 * @param helper the {@code PipelineHelper} describing the stream pipeline
 *               up to this operation
 * @param spliterator the {@code Spliterator} describing the source for this
 *                    pipeline
 */
func (task *TODOShortCircuitTask) WithSpliterator(spliterator spliterator.Spliterator) *TODOShortCircuitTask {
	task.TODOTask.WithSpliterator(spliterator)
	task.sharedResult = &sharedResult{}
	return task
}

func (task *TODOShortCircuitTask) SharedResult() *sharedResult {
	return task.sharedResult
}

/**
 * Constructor for non-root nodes.
 *
 * @param parent parent task in the computation tree
 * @param spliterator the {@code Spliterator} for the portion of the
 *                    computation tree described by this task
 */
func (task *TODOShortCircuitTask) WithParent(parent ShortCircuitTask, spliterator spliterator.Spliterator) *TODOShortCircuitTask {
	task.TODOTask.WithParent(parent, spliterator)
	task.sharedResult = parent.SharedResult()
	task.SetDerived(task)
	return task
}

// Helper

/**
 * Overrides TODOTask version to include checks for early
 * exits while splitting or computing.
 */
func (task *TODOShortCircuitTask) Compute(ctx context.Context) {
	rs := task.spliterator
	var ls spliterator.Spliterator
	sizeEstimate := rs.EstimateSize()
	sizeThreshold := task.getTargetSize(sizeEstimate)

	var this = task.GetDerivedElse(task).(Task)
	var forkRight bool
	var sr = task.sharedResult
	for sr.Load() == nil {
		select {
		case <-ctx.Done():
			return
		default:
		}

		if sizeEstimate <= sizeThreshold {
			break
		}

		ls = rs.TrySplit()
		if ls == nil {
			break
		}

		var leftChild, rightChild, taskToFork Task
		leftChild = this.MakeChild(ls)
		rightChild = this.MakeChild(rs)
		this.SetLeftChild(leftChild)
		this.SetRightChild(rightChild)

		if forkRight {
			forkRight = false
			rs = ls
			this = leftChild
			taskToFork = rightChild
		} else {
			forkRight = true
			this = rightChild
			taskToFork = leftChild
		}

		// fork
		taskToFork.Fork(ctx)

		sizeEstimate = rs.EstimateSize()
	}
	this.SetLocalResult(this.DoLeaf(ctx))
}

func (task *TODOShortCircuitTask) ShortCircuit(result Sink) {
	if result != nil {
		task.sharedResult.Store(result)
	}
}

func (task *TODOShortCircuitTask) GetSharedResult() Sink {
	return task.sharedResult.Load()
}

/**
 * Does nothing; instead, subclasses should use
 * {@link #setLocalResult(Object)}} to manage results.
 *
 * @param result must be null, or an exception is thrown (this is a safety
 *        tripwire to detect when {@code setRawResult()} is being used
 *        instead of {@code setLocalResult()}
 */
func (task TODOShortCircuitTask) SetRawResult(result Sink) {
	if result != nil {
		panic(exception.NewIllegalStateException())
	}
}

/**
 * Sets a local result for this task.  If this task is the root, set the
 * shared result instead (if not already set).
 *
 * @param localResult The result to set for this task
 */
func (task *TODOShortCircuitTask) SetLocalResult(localResult Sink) {
	var this = task.GetDerivedElse(task).(Task)
	if this.IsRoot() {
		if localResult != nil {
			task.sharedResult.Store(localResult)
		}
	} else {
		task.TODOTask.SetLocalResult(localResult)
	}
}

/**
 * Retrieves the local result for this task
 */
func (task *TODOShortCircuitTask) getRawResult() Sink {
	var this = task.GetDerivedElse(task).(Task)
	return this.GetLocalResult()
}

/**
 * Retrieves the local result for this task.  If this task is the root,
 * retrieves the shared result instead.
 */
func (task *TODOShortCircuitTask) GetLocalResult() Sink {
	var this = task.GetDerivedElse(task).(Task)
	if this.IsRoot() {
		answer := task.sharedResult.Load()
		return answer
	}
	return task.TODOTask.GetLocalResult()
}

func (task *TODOShortCircuitTask) Cancel() {
	task.canceled = true
}

func (task *TODOShortCircuitTask) TaskCanceled() bool {
	if task.canceled {
		return true
	}
	var this = task.GetDerivedElse(task).(Task)
	parent := this.GetParent()
	if parent != nil {
		return parent.(ShortCircuitTask).TaskCanceled()
	}
	return false
}

func (task *TODOShortCircuitTask) CancelLaterNodes() {
	// Go up the tree, cancel right siblings of this node and all parents
	parent := task.GetParent()
	var node = task.GetDerivedElse(task).(Task)

	for parent != nil {
		// If node is a left child of parent, then has a right sibling
		if parent.LeftChild() == node {
			rightSibling := parent.RightChild()
			if !rightSibling.(*TODOShortCircuitTask).canceled {
				rightSibling.(ShortCircuitTask).Cancel()
			}
		}

		node = parent
		parent = parent.GetParent()
	}
}
