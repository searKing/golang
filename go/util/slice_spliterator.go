package util

import (
	"context"
	"github.com/searKing/golang/go/error/exception"
	"github.com/searKing/golang/go/util/function/consumer"
	"github.com/searKing/golang/go/util/object"
)

type AbstractSliceSpliteratorClass struct {
	AbstractSpliteratorClass

	array []interface{}

	index int // current index, modified on advance/split
	fence int // one past last index
	cs    SpliteratorCharacteristic
}

func NewSliceSpliterator2(array []interface{}, acs SpliteratorCharacteristic) *AbstractSliceSpliteratorClass {
	return NewSliceSpliterator4(array, 0, len(array), acs)
}

func NewSliceSpliterator4(array []interface{}, origin int, fence int, acs SpliteratorCharacteristic) *AbstractSliceSpliteratorClass {
	split := &AbstractSliceSpliteratorClass{
		array: array,
		index: origin,
		fence: fence,
		cs:    acs | SpliteratorOrdered | SpliteratorSized | SpliteratorSubsized,
	}
	split.SetDerived(split)
	return split
}

// Helper
func (split *AbstractSliceSpliteratorClass) follow() Spliterator {
	derived := split.GetDerived()
	if derived == nil {
		return split
	}
	return derived.(Spliterator)
}

func (split *AbstractSliceSpliteratorClass) TrySplit() Spliterator {
	lo := split.index
	mid := (lo + split.fence) >> 1
	if lo >= mid {
		return nil
	}
	split.index = mid
	return NewSliceSpliterator4(split.array, lo, mid, split.cs)
}

func (split *AbstractSliceSpliteratorClass) ForEachRemaining(ctx context.Context, action consumer.Consumer) {
	if action == nil {
		panic(exception.NewNullPointerException())
	}
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

func (split *AbstractSliceSpliteratorClass) TryAdvance(action consumer.Consumer) bool {
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

func (split *AbstractSliceSpliteratorClass) EstimateSize() int {
	return split.fence - split.index
}

func (split *AbstractSliceSpliteratorClass) Characteristics() SpliteratorCharacteristic {
	return split.cs
}

func (split *AbstractSliceSpliteratorClass) GetComparator() object.Comparator {
	class := split.follow()
	if class.HasCharacteristics(SpliteratorSorted) {
		return nil
	}
	panic(exception.NewIllegalStateException())
}
