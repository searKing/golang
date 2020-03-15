// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// A sequence of elements supporting sequential and parallel aggregate
// operations.  The following example illustrates an aggregate operation using
// SEE java/util/function/Consumer.java
package consumer

import (
	"github.com/searKing/golang/go/error/exception"
	"github.com/searKing/golang/go/util/class"
	"github.com/searKing/golang/go/util/object"
)

/**
 * Represents an operation that accepts a single input argument and returns no
 * result. Unlike most other functional interfaces, {@code Consumer} is expected
 * to operate via side-effects.
 *
 * <p>This is a <a href="package-summary.html">functional interface</a>
 * whose functional method is {@link #accept(Object)}.
 *
 * @param <T> the type of the input to the operation
 *
 * @since 1.8
 */
type Consumer interface {
	/**
	 * Performs this operation on the given argument.
	 *
	 * @param t the input argument
	 */
	Accept(t interface{})

	/**
	 * Returns a composed {@code Consumer} that performs, in sequence, this
	 * operation followed by the {@code after} operation. If performing either
	 * operation throws an exception, it is relayed to the caller of the
	 * composed operation.  If performing this operation throws an exception,
	 * the {@code after} operation will not be performed.
	 *
	 * @param after the operation to perform after this operation
	 * @return a composed {@code Consumer} that performs in sequence this
	 * operation followed by the {@code after} operation
	 * @throws NullPointerException if {@code after} is null
	 */
	AndThen(after Consumer) Consumer
}

type ConsumerFunc func(t interface{})

// Accept calls f(t).
func (f ConsumerFunc) Accept(t interface{}) {
	f(t)
}

func (f ConsumerFunc) AndThen(after Consumer) Consumer {
	object.RequireNonNil(after)
	return ConsumerFunc(func(t interface{}) {
		f.Accept(t)
		after.Accept(t)
	})
}

type TODO struct {
	class.Class
}

func (consumer *TODO) Accept(t interface{}) {
	panic(exception.NewIllegalStateException1("called wrong Accept method"))
}

func (consumer *TODO) AndThen(after Consumer) Consumer {
	object.RequireNonNil(after)
	return ConsumerFunc(func(t interface{}) {
		consumer.GetDerived().(Consumer).Accept(t)
		after.Accept(t)
	})
}
