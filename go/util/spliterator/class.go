package spliterator

import (
	"context"
	"github.com/searKing/golang/go/error/exception"
	"github.com/searKing/golang/go/util/class"
	"github.com/searKing/golang/go/util/function/consumer"
	"github.com/searKing/golang/go/util/object"
)

type Class struct {
	class.Class
}

func (split *Class) TryAdvance(action consumer.Consumer) bool {
	panic(exception.NewIllegalStateException1("called unimplemented TryAdvance method"))
}

func (split *Class) ForEachRemaining(ctx context.Context, action consumer.Consumer) {
	class := split.GetDerivedElse(split).(Spliterator)
	for class.TryAdvance(ctx, action) {
	}
	return
}

func (split *Class) TrySplit() Spliterator {
	panic(exception.NewIllegalStateException1("called unimplemented TrySplit method"))
}

func (split *Class) EstimateSize() int {
	panic(exception.NewIllegalStateException1("called unimplemented EstimateSize method"))
}

func (split *Class) GetExactSizeIfKnown() int {
	class := split.GetDerivedElse(split).(Spliterator)
	if class.Characteristics()&CharacteristicSized == 0 {
		return -1
	}
	return class.EstimateSize()
}

func (split *Class) Characteristics() Characteristic {
	panic(exception.NewIllegalStateException1("called unimplemented Characteristics method"))
}

func (split *Class) HasCharacteristics(characteristics Characteristic) bool {
	class := split.GetDerivedElse(split).(Spliterator)
	return (class.Characteristics() & characteristics) == characteristics
}

func (split *Class) GetComparator() object.Comparator {
	panic(exception.NewIllegalStateException())
}
