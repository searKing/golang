// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package spliterator

import (
	"context"

	"github.com/searKing/golang/go/util/function/consumer"
	"github.com/searKing/golang/go/util/object"
)

type EmptySpliterator struct {
	TODO
}

func NewEmptySpliterator() *EmptySpliterator {
	split := &EmptySpliterator{}
	split.SetDerived(split)
	return split
}

// Helper
func (split *EmptySpliterator) follow() Spliterator {
	derived := split.GetDerived()
	if derived == nil {
		return split
	}
	return derived.(Spliterator)
}

func (split *EmptySpliterator) TrySplit() Spliterator {
	return nil
}

func (split *EmptySpliterator) TryAdvance(ctx context.Context, consumer consumer.Consumer) bool {
	object.RequireNonNil(consumer)
	return false
}

func (split *EmptySpliterator) ForEachRemaining(ctx context.Context, consumer consumer.Consumer) {
	object.RequireNonNil(consumer)
}

func (split *EmptySpliterator) EstimateSize() int {
	return 0
}

func (split *EmptySpliterator) Characteristics() Characteristic {
	return CharacteristicSized | CharacteristicSubsized
}
