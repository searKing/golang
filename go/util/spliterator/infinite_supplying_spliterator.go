// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package spliterator

import (
	"context"

	"github.com/searKing/golang/go/util/function/consumer"
	"github.com/searKing/golang/go/util/function/supplier"
	"github.com/searKing/golang/go/util/object"
)

/**
 * A Spliterator that infinitely supplies elements in no particular order.
 *
 * <p>Splitting divides the estimated size in two and stops when the
 * estimate size is 0.
 *
 * <p>The {@code forEachRemaining} method if invoked will never terminate.
 * The {@code tryAdvance} method always returns true.
 *
 */
type infiniteSupplyingSpliterator struct {
	TODO
	estimate int
	// The underlying spliterator
	s supplier.Supplier
}

func NewInfiniteSupplyingSpliterator(size int, s supplier.Supplier) Spliterator {
	split := &infiniteSupplyingSpliterator{
		estimate: size,
		s:        s,
	}
	split.SetDerived(split)
	return split
}

func (split *infiniteSupplyingSpliterator) TrySplit() Spliterator {
	if split.estimate == 0 {
		return nil
	}
	return NewInfiniteSupplyingSpliterator(split.estimate>>1, split.s)
}

func (split *infiniteSupplyingSpliterator) TryAdvance(ctx context.Context, action consumer.Consumer) bool {
	object.RequireNonNil(action)
	action.Accept(split.s.Get())
	return true
}

func (split *infiniteSupplyingSpliterator) EstimateSize() int {
	return split.estimate
}

func (split *infiniteSupplyingSpliterator) Characteristics() Characteristic {
	return CharacteristicImmutable
}
