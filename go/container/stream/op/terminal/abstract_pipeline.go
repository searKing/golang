// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package terminal

import (
	"context"
	"runtime"

	"github.com/searKing/golang/go/error/exception"
	"github.com/searKing/golang/go/util/class"
	"github.com/searKing/golang/go/util/optional"
	"github.com/searKing/golang/go/util/spliterator"
)

const (
	MsgStreamLinked = "stream has already been operated upon or closed"
	MsgConsumed     = "source already consumed or closed"
)

var (
	ParallelTargetSize = runtime.GOMAXPROCS(-1)
)

/**
 * Abstract base class for "pipeline" classes, which are the core
 * implementations of the Stream interface and its primitive specializations.
 * Manages construction and evaluation of stream pipelines.
 *
 * <p>An {@code AbstractPipeline} represents an initial portion of a stream
 * pipeline, encapsulating a stream source and zero or more intermediate
 * operations.  The individual {@code AbstractPipeline} objects are often
 * referred to as <em>stages</em>, where each stage describes either the stream
 * source or an intermediate operation.
 *
 * <p>A concrete intermediate stage is generally built from an
 * {@code AbstractPipeline}, a shape-specific pipeline class which extends it
 * (e.g., {@code IntPipeline}) which is also abstract, and an operation-specific
 * concrete class which extends that.  {@code AbstractPipeline} contains most of
 * the mechanics of evaluating the pipeline, and implements methods that will be
 * used by the operation; the shape-specific classes add helper methods for
 * dealing with collection of results into the appropriate shape-specific
 * containers.
 *
 * <p>After chaining a new intermediate operation, or executing a terminal
 * operation, the stream is considered to be consumed, and no more intermediate
 * or terminal operations are permitted on this stream instance.
 *
 * @implNote
 * <p>For sequential streams, and parallel streams without
 * <a href="package-summary.html#StreamOps">stateful intermediate
 * operations</a>, parallel streams, pipeline evaluation is done in a single
 * pass that "jams" all the operations together.  For parallel streams with
 * stateful operations, execution is divided into segments, where each
 * stateful operations marks the end of a segment, and each segment is
 * evaluated separately and the result used as the input to the next
 * segment.  In all cases, the source data is not consumed until a terminal
 * operation begins.
 *
 * @param <E_IN>  type of input elements
 * @param <E_OUT> type of output elements
 * @param <S> type of the subclass implementing {@code BaseStream}
 * @since 1.8
 */
type AbstractPipeline struct {
	class.Class
	/**
	 * True if pipeline is parallel, otherwise the pipeline is sequential; only
	 * valid for the source stage.
	 */
	parallel bool

	/**
	 * True if this pipeline has been linked or consumed
	 */
	linkedOrConsumed bool

	/** Target split size, common to all tasks in a computation */
	targetSize int // may be lazily initialized
}

func (p *AbstractPipeline) IsParallel() bool {
	return p.parallel
}

func (p *AbstractPipeline) SetParallel(parallel bool) {
	p.parallel = parallel
}

func (p *AbstractPipeline) IsLinkedOrConsumed() bool {
	return p.linkedOrConsumed
}

func (p *AbstractPipeline) SetLinkedOrConsumed(linkedOrConsumed bool) {
	p.linkedOrConsumed = linkedOrConsumed
}

/**
 * Returns the targetSize, initializing it via the supplied
 * size estimate if not already initialized.
 */
func (p *AbstractPipeline) GetTargetSize(sizeEstimate int) int {
	if p.targetSize == 0 {
		p.targetSize = p.SuggestTargetSize(sizeEstimate)
	}
	return p.targetSize
}

func (p *AbstractPipeline) SetTargetSize(targetSize int) {
	p.targetSize = targetSize
}

/**
 * Returns a suggested target leaf size based on the initial size estimate.
 *
 * @return suggested target leaf size
 */
func (p *AbstractPipeline) SuggestTargetSize(sizeEstimate int) int {
	est := sizeEstimate / ParallelTargetSize
	if est > 0 {
		return est
	}
	return 1
}

// Terminal evaluation methods

/**
 * Evaluate the pipeline with a terminal operation to produce a result.
 *
 * @param <R> the type of result
 * @param terminalOp the terminal operation to be applied to the pipeline.
 * @return the result
 */
func (p *AbstractPipeline) Evaluate(ctx context.Context, terminalOp Operation, split spliterator.Spliterator) optional.Optional {
	if p.IsLinkedOrConsumed() {
		panic(exception.NewIllegalStateException1(MsgStreamLinked))
	}
	p.SetLinkedOrConsumed(true)
	defer p.SetLinkedOrConsumed(false)

	if p.IsParallel() {
		return terminalOp.EvaluateParallel(ctx, split)
	}
	return terminalOp.EvaluateSequential(ctx, split)
}
