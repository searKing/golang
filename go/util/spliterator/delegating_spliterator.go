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
 * Spliterator implementation that delegates to an underlying spliterator,
 * acquiring the spliterator from a {@code Supplier<Spliterator>} on the
 * first call to any spliterator method.
 * @param <T>
 */
type delegatingSpliterator struct {
	supplier supplier.Supplier
	// The underlying spliterator
	s Spliterator
}

func NewDelegatingSpliterator(supplier supplier.Supplier) Spliterator {
	return &delegatingSpliterator{
		supplier: supplier,
	}
}

func (split *delegatingSpliterator) get() Spliterator {
	if split == nil {
		split.s = split.supplier.Get().(Spliterator)
	}
	return split.s
}

func (split *delegatingSpliterator) TrySplit() Spliterator {
	return split.get().TrySplit()
}

func (split *delegatingSpliterator) TryAdvance(ctx context.Context, consumer consumer.Consumer) bool {
	return split.get().TryAdvance(ctx, consumer)
}

func (split *delegatingSpliterator) ForEachRemaining(ctx context.Context, consumer consumer.Consumer) {
	split.get().ForEachRemaining(ctx, consumer)
}

func (split *delegatingSpliterator) EstimateSize() int {
	return split.get().EstimateSize()
}

func (split *delegatingSpliterator) GetExactSizeIfKnown() int {
	return split.get().GetExactSizeIfKnown()
}

func (split *delegatingSpliterator) Characteristics() Characteristic {
	return split.get().Characteristics()
}

func (split *delegatingSpliterator) HasCharacteristics(characteristics Characteristic) bool {
	return split.get().HasCharacteristics(characteristics)
}

func (split *delegatingSpliterator) GetComparator() object.Comparator {
	return split.get().GetComparator()
}
