// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package spliterator

import (
	"context"

	"github.com/searKing/golang/go/util/function/consumer"
	"github.com/searKing/golang/go/util/object"
)

type OfPrimitive struct {
	TODO
}

func NewOfPrimitive() *OfPrimitive {
	split := &OfPrimitive{}
	split.SetDerived(split)
	return split
}

func (split *OfPrimitive) TrySplit() Spliterator {
	return nil
}

/**
 * If a remaining element exists, performs the given action on it,
 * returning {@code true}; else returns {@code false}.  If this
 * Spliterator is {@link #ORDERED} the action is performed on the
 * next element in encounter order.  Exceptions thrown by the
 * action are relayed to the caller.
 *
 * @param action The action
 * @return {@code false} if no remaining elements existed
 * upon entry to this method, else {@code true}.
 * @throws NullPointerException if the specified action is null
 */
func (split *OfPrimitive) TryAdvance(ctx context.Context, action consumer.Consumer) bool {
	object.RequireNonNil(action)
	return false
}

/**
 * Performs the given action for each remaining element, sequentially in
 * the current thread, until all elements have been processed or the
 * action throws an exception.  If this Spliterator is {@link #ORDERED},
 * actions are performed in encounter order.  Exceptions thrown by the
 * action are relayed to the caller.
 *
 * @implSpec
 * The default implementation repeatedly invokes {@link #tryAdvance}
 * until it returns {@code false}.  It should be overridden whenever
 * possible.
 *
 * @param action The action
 * @throws NullPointerException if the specified action is null
 */
func (split *OfPrimitive) ForEachRemaining(ctx context.Context, action consumer.Consumer) {
	for {
		if !split.GetDerivedElse(split).(Spliterator).TryAdvance(ctx, action) {
			break
		}
	}
}
