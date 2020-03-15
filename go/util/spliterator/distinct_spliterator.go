// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package spliterator

import (
	"context"
	"sync"

	"github.com/searKing/golang/go/util/function/consumer"
	"github.com/searKing/golang/go/util/object"
	"github.com/searKing/golang/go/util/optional"
)

//go:generate go-syncmap -type "seenMap<interface{}, struct{}>"
type seenMap sync.Map

/**
 * A wrapping spliterator that only reports distinct elements of the
 * underlying spliterator. Does not preserve size and encounter order.
 */
type distinctSpliterator struct {
	consumer.TODO
	// The underlying spliterator
	s Spliterator
	// ConcurrentHashMap holding distinct elements as keys
	seen *seenMap
	// Temporary element, only used with tryAdvance
	tmpSlot optional.Optional
}

func NewDistinctSpliterator(s Spliterator) Spliterator {
	return NewDistinctSpliterator2(s, &seenMap{})
}

func NewDistinctSpliterator2(s Spliterator, seen *seenMap) Spliterator {
	split := &distinctSpliterator{
		s:    s,
		seen: seen,
	}
	split.SetDerived(split)
	return split
}

func (split *distinctSpliterator) Accept(t interface{}) {
	split.tmpSlot = optional.Of(t)
}

func (split *distinctSpliterator) TrySplit() Spliterator {
	if ls := split.s.TrySplit(); ls != nil {
		return NewDistinctSpliterator2(ls, split.seen)
	}
	return nil
}

func (split *distinctSpliterator) TryAdvance(ctx context.Context, action consumer.Consumer) bool {
	for {
		if !split.s.TryAdvance(ctx, split) {
			break
		}
		if !split.tmpSlot.IsPresent() {
			break
		}
		if _, loaded := split.seen.LoadOrStore(split.tmpSlot.Get(), struct{}{}); !loaded {
			action.Accept(split.tmpSlot.Get())
			split.tmpSlot = optional.Empty()
			return true
		}
	}
	return false
}

func (split *distinctSpliterator) ForEachRemaining(ctx context.Context, action consumer.Consumer) {
	split.s.ForEachRemaining(ctx, consumer.ConsumerFunc(func(t interface{}) {
		if _, loaded := split.seen.LoadOrStore(t, struct{}{}); !loaded {
			action.Accept(t)
		}
	}))
}

func (split *distinctSpliterator) EstimateSize() int {
	return split.s.EstimateSize()
}

func (split *distinctSpliterator) GetExactSizeIfKnown() int {
	return split.s.GetExactSizeIfKnown()
}

func (split *distinctSpliterator) Characteristics() Characteristic {
	return (split.s.Characteristics() & (^(CharacteristicSized |
		CharacteristicSubsized |
		CharacteristicSorted |
		CharacteristicOrdered))) |
		CharacteristicDistinct
}

func (split *distinctSpliterator) HasCharacteristics(characteristics Characteristic) bool {
	return split.s.HasCharacteristics(characteristics)
}

func (split *distinctSpliterator) GetComparator() object.Comparator {
	return split.s.GetComparator()
}
