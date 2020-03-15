// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package terminal

import (
	"context"

	"github.com/searKing/golang/go/error/exception"
	"github.com/searKing/golang/go/util/class"
	"github.com/searKing/golang/go/util/optional"
	"github.com/searKing/golang/go/util/spliterator"
)

/**
 * An operation in a stream pipeline that takes a stream as input and produces
 * a result or side-effect.  A {@code Op} has an input type and stream
 * shape, and a result type.  A {@code Op} also has a set of
 * <em>operation flags</em> that describes how the operation processes elements
 * of the stream (such as short-circuiting or respecting encounter order; see
 * {@link StreamOpFlag}).
 *
 * <p>A {@code Op} must provide a sequential and parallel implementation
 * of the operation relative to a given stream source and set of intermediate
 * operations.
 *
 * @param <E_IN> the type of input elements
 * @param <R>    the type of the result
 * @since 1.8
 */
type Operation interface {
	/**
	 * Gets the stream flags of the operation.  Terminal operations may set a
	 * limited subset of the stream flags defined in {@link StreamOpFlag}, and
	 * these flags are combined with the previously combined stream and
	 * intermediate operation flags for the pipeline.
	 *
	 * @implSpec The default implementation returns zero.
	 *
	 * @return the stream flags for this operation
	 * @see StreamOpFlag
	 */
	GetOpFlags() int

	/**
	 * Performs a parallel evaluation of the operation using the specified
	 * {@code PipelineHelper}, which describes the upstream intermediate
	 * operations.
	 *
	 * @implSpec The default performs a sequential evaluation of the operation
	 * using the specified {@code PipelineHelper}.
	 *
	 * @param helper the pipeline helper
	 * @param spliterator the source spliterator
	 * @return the result of the evaluation
	 */
	EvaluateParallel(ctx context.Context, spliterator spliterator.Spliterator) optional.Optional

	/**
	 * Performs a sequential evaluation of the operation using the specified
	 * {@code PipelineHelper}, which describes the upstream intermediate
	 * operations.
	 *
	 * @param helper the pipeline helper
	 * @param spliterator the source spliterator
	 * @return the result of the evaluation
	 */
	EvaluateSequential(ctx context.Context, spliterator spliterator.Spliterator) optional.Optional
}

type TODOOperation struct {
	class.Class
}

func (op *TODOOperation) GetOpFlags() int {
	return 0
}

func (op *TODOOperation) EvaluateParallel(ctx context.Context, spliterator spliterator.Spliterator) optional.Optional {
	c := op.GetDerivedElse(op).(Operation)
	return c.EvaluateSequential(ctx, spliterator)
}

func (op *TODOOperation) EvaluateSequential(ctx context.Context, spliterator spliterator.Spliterator) optional.Optional {
	panic(exception.NewIllegalStateException())
}
