// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package spliterator

import (
	"context"

	"github.com/searKing/golang/go/error/exception"
	"github.com/searKing/golang/go/util/function/consumer"
	"github.com/searKing/golang/go/util/object"
)

/**
 * A slice Spliterator from a source Spliterator that reports
 * {@code SUBSIZED}.
 *
 */
type sliceSpliterator struct {
	OfPrimitive

	array []interface{}

	index int // current index, modified on advance/split
	fence int // one past last index
	cs    Characteristic
}

func NewSliceSpliterator2(acs Characteristic, arrays ...interface{}) Spliterator {
	return NewSliceSpliterator4(0, len(arrays), acs, arrays...)
}

func NewSliceSpliterator4(origin int, fence int, acs Characteristic, arrays ...interface{}) Spliterator {
	split := &sliceSpliterator{
		array: arrays,
		index: origin,
		fence: fence,
		cs:    acs | CharacteristicOrdered | CharacteristicSized | CharacteristicSubsized,
	}
	split.SetDerived(split)
	return split
}

func (split *sliceSpliterator) TrySplit() Spliterator {
	lo := split.index
	mid := (lo + split.fence) >> 1
	if lo >= mid {
		return nil
	}
	split.index = mid
	return NewSliceSpliterator4(lo, mid, split.cs, split.array...)
}

func (split *sliceSpliterator) ForEachRemaining(ctx context.Context, action consumer.Consumer) {
	object.RequireNonNil(action)

	var a []interface{}
	var i, hi int // hoist accesses and checks from loop
	a = split.array

	hi = split.fence
	i = split.index
	split.index = hi
	if len(a) >= hi && i >= 0 && i < hi {
		for ; i < hi; i++ {
			select {
			case <-ctx.Done():
				return
			default:
			}
			action.Accept(a[i])
		}
	}
	return
}

func (split *sliceSpliterator) TryAdvance(ctx context.Context, action consumer.Consumer) bool {
	if action == nil {
		panic(exception.NewNullPointerException())
	}
	if split.index >= 0 && split.index < split.fence {
		action.Accept(split.array[split.index])
		split.index++
		return true
	}
	return false
}

func (split *sliceSpliterator) EstimateSize() int {
	return split.fence - split.index
}

func (split *sliceSpliterator) Characteristics() Characteristic {
	return split.cs
}

func (split *sliceSpliterator) GetComparator() object.Comparator {
	if split.GetDerivedElse(split).(Spliterator).HasCharacteristics(CharacteristicSorted) {
		return nil
	}
	panic(exception.NewIllegalStateException())
}
