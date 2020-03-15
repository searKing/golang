// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package terminal

import (
	"context"
	"runtime"
	"sync"

	"github.com/searKing/golang/go/error/exception"
	"github.com/searKing/golang/go/util/class"
	"github.com/searKing/golang/go/util/spliterator"
)

type Task interface {
	GetSpliterator() spliterator.Spliterator
	/**
	 * Returns the parent of this task, or null if this task is the root
	 *
	 * @return the parent of this task, or null if this task is the root
	 */
	GetParent() Task

	/**
	 * The left child.
	 * null if no children
	 * if non-null rightChild is non-null
	 *
	 * @return the right child
	 */
	LeftChild() Task

	SetLeftChild(task Task)

	/**
	 * The right child.
	 * null if no children
	 * if non-null rightChild is non-null
	 *
	 * @return the right child
	 */
	RightChild() Task

	SetRightChild(task Task)

	/**
	 * Indicates whether this task is a leaf node.  (Only valid after
	 * {@link #compute} has been called on this node).  If the node is not a
	 * leaf node, then children will be non-null and numChildren will be
	 * positive.
	 *
	 * @return {@code true} if this task is a leaf node
	 */
	IsLeaf() bool

	/**
	 * Indicates whether this task is the root node
	 *
	 * @return {@code true} if this task is the root node.
	 */
	IsRoot() bool

	/**
	 * Target leaf size, common to all tasks in a computation
	 *
	 * @return target leaf size.
	 */
	TargetSize() int

	/**
	 * Default target of leaf tasks for parallel decomposition.
	 * To allow load balancing, we over-partition, currently to approximately
	 * four tasks per processor, which enables others to help out
	 * if leaf tasks are uneven or some processors are otherwise busy.
	 *
	 * @return the default target size of leaf tasks
	 */
	GetLeafTarget() int

	/**
	 * Constructs a new node of type T whose parent is the receiver; must call
	 * the TODOTask(T, Spliterator) constructor with the receiver and the
	 * provided Spliterator.
	 *
	 * @param spliterator {@code Spliterator} describing the subtree rooted at
	 *        this node, obtained by splitting the parent {@code Spliterator}
	 * @return newly constructed child node
	 */
	MakeChild(spliterator spliterator.Spliterator) Task

	/**
	 * Computes the result associated with a leaf node.  Will be called by
	 * {@code compute()} and the result passed to @{code setLocalResult()}
	 *
	 * @return the computed result of a leaf node
	 */
	DoLeaf(ctx context.Context) Sink

	/**
	 * Retrieves a result previously stored with {@link #setLocalResult}
	 *
	 * @return local result for this node previously stored with
	 * {@link #setLocalResult}
	 */
	GetLocalResult() Sink

	/**
	 * Associates the result with the task, can be retrieved with
	 * {@link #GetLocalResult}
	 *
	 * @param localResult local result for this node
	 */
	SetLocalResult(localResult Sink)

	/**
	 * Decides whether or not to split a task further or compute it
	 * directly. If computing directly, calls {@code doLeaf} and pass
	 * the result to {@code setRawResult}. Otherwise splits off
	 * subtasks, forking one and continuing as the other.
	 *
	 * <p> The method is structured to conserve resources across a
	 * range of uses.  The loop continues with one of the child tasks
	 * when split, to avoid deep recursion. To cope with spliterators
	 * that may be systematically biased toward left-heavy or
	 * right-heavy splits, we alternate which child is forked versus
	 * continued in the loop.
	 */
	Compute(ctx context.Context)

	/**
	 * {@inheritDoc}
	 *
	 * @implNote
	 * Clears spliterator and children fields.  Overriders MUST call
	 * {@code super.onCompletion} as the last thing they do if they want these
	 * cleared.
	 */
	OnCompletion(caller Task)

	/**
	 * Returns whether this node is a "leftmost" node -- whether the path from
	 * the root to this node involves only traversing leftmost child links.  For
	 * a leaf node, this means it is the first leaf node in the encounter order.
	 *
	 * @return {@code true} if this node is a "leftmost" node
	 */
	IsLeftmostNode() bool

	/**
	 * Arranges to asynchronously execute this task in the pool the
	 * current task is running in, if applicable, or using the {@link
	 * ForkJoinPool#commonPool()} if not {@link #inForkJoinPool}.  While
	 * it is not necessarily enforced, it is a usage error to fork a
	 * task more than once unless it has completed and been
	 * reinitialized.  Subsequent modifications to the state of this
	 * task or any data it operates on are not necessarily
	 * consistently observable by any thread other than the one
	 * executing it unless preceded by a call to {@link #join} or
	 * related methods, or a call to {@link #isDone} returning {@code
	 * true}.
	 *
	 * @return {@code this}, to simplify usage
	 */
	Fork(ctx context.Context)

	/**
	 * Returns the result of the computation when it
	 * {@linkplain #isDone is done}.
	 * This method differs from {@link #get()} in that abnormal
	 * completion results in {@code RuntimeException} or {@code Error},
	 * not {@code ExecutionException}, and that interrupts of the
	 * calling thread do <em>not</em> cause the method to abruptly
	 * return by throwing {@code InterruptedException}.
	 *
	 * @return the computed result
	 */
	Join() Sink

	/**
	 * Commences performing this task, awaits its completion if
	 * necessary, and returns its result, or throws an (unchecked)
	 * {@code RuntimeException} or {@code Error} if the underlying
	 * computation did so.
	 *
	 * @return the computed result
	 */
	Invoke(ctx context.Context) Sink
}

type TODOTask struct {
	class.Class

	parent Task

	/**
	 * The spliterator for the portion of the input associated with the subtree
	 * rooted at this task
	 */
	spliterator spliterator.Spliterator

	/** Target leaf size, common to all tasks in a computation */
	targetSize int // may be lazily initialized
	/**
	 * The left child.
	 * null if no children
	 * if non-null rightChild is non-null
	 */
	leftChild Task

	/**
	 * The right child.
	 * null if no children
	 * if non-null leftChild is non-null
	 */
	rightChild Task

	/** The result of this node, if completed */
	localResult Sink

	computer func(ctx context.Context)

	wg sync.WaitGroup
}

/**
 * Constructor for root nodes.
 *
 * @param helper The {@code PipelineHelper} describing the stream pipeline
 *               up to this operation
 * @param spliterator The {@code Spliterator} describing the source for this
 *                    pipeline
 */
func (task *TODOTask) WithSpliterator(spliterator spliterator.Spliterator) *TODOTask {
	task.spliterator = spliterator
	return task
}

/**
 * Constructor for non-root nodes.
 *
 * @param parent this node's parent task
 * @param spliterator {@code Spliterator} describing the subtree rooted at
 *        this node, obtained by splitting the parent {@code Spliterator}
 */
func (task *TODOTask) WithParent(parent Task, spliterator spliterator.Spliterator) *TODOTask {
	task.parent = parent
	task.spliterator = spliterator
	task.targetSize = parent.TargetSize()
	return task
}

func (task *TODOTask) GetSpliterator() spliterator.Spliterator {
	return task.spliterator
}

func (task *TODOTask) GetLeafTarget() int {
	return runtime.GOMAXPROCS(-1) << 2
}

func (task *TODOTask) LeftChild() Task {
	return task.leftChild
}

func (task *TODOTask) SetLeftChild(task_ Task) {
	task.leftChild = task_
}

func (task *TODOTask) RightChild() Task {
	return task.rightChild
}

func (task *TODOTask) SetRightChild(task_ Task) {
	task.rightChild = task_
}

func (task *TODOTask) TargetSize() int {
	return task.targetSize
}

func (task *TODOTask) MakeChild(spliterator spliterator.Spliterator) Task {
	panic(exception.NewUnsupportedOperationException())
	return nil
}

func (task *TODOTask) DoLeaf(ctx context.Context) Sink {
	panic(exception.NewUnsupportedOperationException())
	return nil
}

/**
 * Returns a suggested target leaf size based on the initial size estimate.
 *
 * @return suggested target leaf size
 */
func (task *TODOTask) suggestTargetSize(sizeEstimate int) int {
	this := task.GetDerivedElse(task).(Task)
	est := sizeEstimate / this.GetLeafTarget()
	if est > 0 {
		return est
	}
	return 1
}

/**
 * Returns the targetSize, initializing it via the supplied
 * size estimate if not already initialized.
 */
func (task *TODOTask) getTargetSize(sizeEstimate int) int {
	if task.targetSize != 0 {
		return task.targetSize
	}

	task.targetSize = task.suggestTargetSize(sizeEstimate)
	return task.targetSize
}

/**
 * Returns the local result, if any. Subclasses should use
 * {@link #setLocalResult(Object)} and {@link #GetLocalResult()} to manage
 * results.  This returns the local result so that calls from within the
 * fork-join framework will return the correct result.
 *
 * @return local result for this node previously stored with
 * {@link #setLocalResult}
 */
func (task *TODOTask) getRawResult() Sink {
	return task.localResult
}

/**
 * Does nothing; instead, subclasses should use
 * {@link #setLocalResult(Object)}} to manage results.
 *
 * @param result must be null, or an exception is thrown (this is a safety
 *        tripwire to detect when {@code setRawResult()} is being used
 *        instead of {@code setLocalResult()}
 */
func (task *TODOTask) SetRawResult(result Sink) {
	if result != nil {
		panic(exception.NewIllegalStateException())
	}
}

/**
 * Retrieves a result previously stored with {@link #setLocalResult}
 *
 * @return local result for this node previously stored with
 * {@link #setLocalResult}
 */
func (task *TODOTask) GetLocalResult() Sink {
	return task.localResult
}

/**
 * Associates the result with the task, can be retrieved with
 * {@link #GetLocalResult}
 *
 * @param localResult local result for this node
 */
func (task *TODOTask) SetLocalResult(localResult Sink) {
	task.localResult = localResult
}

func (task *TODOTask) IsLeaf() bool {
	return task.leftChild == nil
}

func (task *TODOTask) IsRoot() bool {
	this := task.GetDerivedElse(task).(Task)
	return this.GetParent() == nil
}

func (task *TODOTask) GetParent() Task {
	return task.parent
}

func (task *TODOTask) Compute(ctx context.Context) {
	rs := task.spliterator
	var ls spliterator.Spliterator
	sizeEstimate := rs.EstimateSize()
	sizeThreshold := task.getTargetSize(sizeEstimate)

	var this = task.GetDerivedElse(task).(Task)
	var forkRight bool
	for {
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

func (task *TODOTask) Fork(ctx context.Context) {
	go func() {
		task.wg.Add(1)
		defer task.wg.Done()
		this := task.GetDerivedElse(task).(Task)
		this.Compute(ctx)
		this.OnCompletion(task.GetDerivedElse(task).(Task))
	}()
}

func (task *TODOTask) Join() Sink {
	task.wg.Wait()
	return task.getRawResult()
}

func (task *TODOTask) Invoke(ctx context.Context) Sink {
	this := task.GetDerivedElse(task).(Task)

	this.Fork(ctx)
	return this.Join()
}

func (task *TODOTask) OnCompletion(caller Task) {
	task.spliterator = nil
	task.leftChild = nil
	task.rightChild = nil
}

func (task *TODOTask) IsLeftmostNode() bool {
	var node = task.GetDerivedElse(task).(Task)
	for node != nil {
		parent := node.GetParent()
		if parent != nil && parent.LeftChild() != node {
			return false
		}
		node = parent
	}
	return true
}
