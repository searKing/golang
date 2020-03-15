// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package spliterator

import (
	"context"

	"github.com/searKing/golang/go/error/exception"
	"github.com/searKing/golang/go/util/class"
	"github.com/searKing/golang/go/util/function/consumer"
	"github.com/searKing/golang/go/util/object"
)

type TODO struct {
	class.Class
}

func (split *TODO) TryAdvance(action consumer.Consumer) bool {
	panic(exception.NewIllegalStateException1("called unimplemented TryAdvance method"))
}

func (split *TODO) ForEachRemaining(ctx context.Context, action consumer.Consumer) {
	class := split.GetDerivedElse(split).(Spliterator)
	for class.TryAdvance(ctx, action) {
	}
	return
}

func (split *TODO) TrySplit() Spliterator {
	panic(exception.NewIllegalStateException1("called unimplemented TrySplit method"))
}

func (split *TODO) EstimateSize() int {
	panic(exception.NewIllegalStateException1("called unimplemented EstimateSize method"))
}

func (split *TODO) GetExactSizeIfKnown() int {
	class := split.GetDerivedElse(split).(Spliterator)
	if class.Characteristics()&CharacteristicSized == 0 {
		return -1
	}
	return class.EstimateSize()
}

func (split *TODO) Characteristics() Characteristic {
	panic(exception.NewIllegalStateException1("called unimplemented Characteristics method"))
}

func (split *TODO) HasCharacteristics(characteristics Characteristic) bool {
	class := split.GetDerivedElse(split).(Spliterator)
	return (class.Characteristics() & characteristics) == characteristics
}

func (split *TODO) GetComparator() object.Comparator {
	panic(exception.NewIllegalStateException())
}
