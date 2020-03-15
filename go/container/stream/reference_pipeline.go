// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stream

import (
	"context"
	"sort"
	"strconv"
	"sync"

	"github.com/searKing/golang/go/container/stream/op/find"
	"github.com/searKing/golang/go/container/stream/op/match"
	"github.com/searKing/golang/go/container/stream/op/reduce"
	"github.com/searKing/golang/go/container/stream/op/terminal"
	"github.com/searKing/golang/go/error/exception"
	"github.com/searKing/golang/go/util"
	"github.com/searKing/golang/go/util/function/binary"
	"github.com/searKing/golang/go/util/function/consumer"
	"github.com/searKing/golang/go/util/function/predicate"
	"github.com/searKing/golang/go/util/object"
	"github.com/searKing/golang/go/util/optional"
	"github.com/searKing/golang/go/util/spliterator"
)

type ReferencePipeline struct {
	upstreams []interface{}

	terminal.AbstractPipeline
}

func New(upstreams ...interface{}) *ReferencePipeline {
	pipe := &ReferencePipeline{
		upstreams: upstreams,
	}
	pipe.SetDerived(pipe)
	return pipe
}

func (r *ReferencePipeline) Filter(ctx context.Context, predicate predicate.Predicater) *ReferencePipeline {
	object.RequireNonNil(predicate)
	var sFiltered []interface{}
	for _, u := range r.upstreams {
		select {
		case <-ctx.Done():
			return &ReferencePipeline{}
		default:
		}
		if predicate.Test(u) {
			sFiltered = append(sFiltered, u)
		}
	}
	return New(sFiltered...)
}

func (r *ReferencePipeline) Map(ctx context.Context, mapper func(interface{}) interface{}) *ReferencePipeline {
	object.RequireNonNil(mapper)
	var sMapped []interface{}
	for _, u := range r.upstreams {
		select {
		case <-ctx.Done():
			return &ReferencePipeline{}
		default:
		}
		sMapped = append(sMapped, mapper(u))
	}
	return New(sMapped...)
}

func (r *ReferencePipeline) Distinct(ctx context.Context, distincter func(interface{}, interface{}) int) *ReferencePipeline {
	object.RequireNonNil(distincter)

	sDistinctMap := map[interface{}]struct{}{}
	var sDistinct []interface{}
	for _, u := range r.upstreams {
		select {
		case <-ctx.Done():
			return &ReferencePipeline{}
		default:
		}
		if uu, ok := sDistinctMap[u]; ok {
			if distincter(uu, u) == 0 {
				continue
			}
		}
		sDistinctMap[u] = struct{}{}
		sDistinct = append(sDistinct, u)
	}
	return New(sDistinct...)
}

func (r *ReferencePipeline) Sorted(lesser func(interface{}, interface{}) int) *ReferencePipeline {
	object.RequireNonNil(lesser)
	sSorted := make([]interface{}, len(r.upstreams))
	copy(sSorted, r.upstreams)

	less := func(i, j int) bool {
		if lesser(sSorted[i], sSorted[j]) < 0 {
			return true
		}
		return false
	}
	sort.Slice(sSorted, less)
	return New(sSorted...)

}

func (r *ReferencePipeline) Peek(ctx context.Context, action consumer.Consumer) *ReferencePipeline {
	object.RequireNonNil(action)

	for _, u := range r.upstreams {
		select {
		case <-ctx.Done():
			return &ReferencePipeline{}
		default:
		}
		action.Accept(u)
	}
	return r
}

func (r *ReferencePipeline) Limit(maxSize int) *ReferencePipeline {
	if maxSize < 0 {
		panic(exception.NewIllegalArgumentException1(strconv.Itoa(maxSize)))
	}
	m := len(r.upstreams)
	if m > maxSize {
		m = maxSize
	}
	return &ReferencePipeline{upstreams: r.upstreams[:m]}
}

func (r *ReferencePipeline) Skip(n int) *ReferencePipeline {
	if n < 0 {
		panic(exception.NewIllegalArgumentException1(strconv.Itoa(n)))
	}
	if n == 0 {
		return r
	}
	m := len(r.upstreams)
	if n > m {
		n = m
	}
	return &ReferencePipeline{upstreams: r.upstreams[n:]}
}

func (r *ReferencePipeline) TakeWhile(ctx context.Context, predicate predicate.Predicater) *ReferencePipeline {
	object.RequireNonNil(predicate)

	var sTaken []interface{}
	for _, u := range r.upstreams {
		select {
		case <-ctx.Done():
			return &ReferencePipeline{}
		default:
		}
		if predicate.Test(u) {
			sTaken = append(sTaken, u)
			continue
		}
		break
	}
	return &ReferencePipeline{upstreams: sTaken}
}

func (r *ReferencePipeline) TakeUntil(ctx context.Context, predicate predicate.Predicater) *ReferencePipeline {
	object.RequireNonNil(predicate)

	return r.TakeWhile(ctx, predicate.Negate())
}

func (r *ReferencePipeline) DropWhile(ctx context.Context, predicate predicate.Predicater) *ReferencePipeline {
	object.RequireNonNil(predicate)

	var sTaken []interface{}
	for _, u := range r.upstreams {
		select {
		case <-ctx.Done():
			return &ReferencePipeline{}
		default:
		}
		if predicate.Test(u) {
			continue
		}
		sTaken = append(sTaken, r)
	}
	return &ReferencePipeline{upstreams: sTaken}
}

func (r *ReferencePipeline) DropUntil(ctx context.Context, predicate predicate.Predicater) *ReferencePipeline {
	object.RequireNonNil(predicate)

	return r.DropWhile(ctx, predicate.Negate())
}

func (r *ReferencePipeline) ForEach(ctx context.Context, action consumer.Consumer) {
	object.RequireNonNil(action)

	var wg sync.WaitGroup
	for _, u := range r.upstreams {
		select {
		case <-ctx.Done():
			return
		default:
		}
		wg.Add(1)
		go func(uu interface{}) {
			defer wg.Done()
			action.Accept(uu)
		}(u)
	}
	wg.Wait()
}

func (r *ReferencePipeline) ForEachOrdered(ctx context.Context, action consumer.Consumer) {
	object.RequireNonNil(action)

	for _, u := range r.upstreams {
		select {
		case <-ctx.Done():
			return
		default:
		}
		action.Accept(u)
	}
}

func (r *ReferencePipeline) ToSlice(generator func(interface{}) interface{}) []interface{} {
	if generator == nil {
		return r.upstreams
	}

	downstream := make([]interface{}, len(r.upstreams))

	for _, u := range r.upstreams {
		downstream = append(downstream, generator(u))
	}
	return downstream
}

func (r *ReferencePipeline) Reduce(ctx context.Context,
	accumulator binary.BiFunction, combiner binary.BiFunction,
	identity ...interface{}) optional.Optional {
	split := spliterator.NewSliceSpliterator2(spliterator.CharacteristicTODO, append(identity, r.upstreams...))
	return r.Evaluate(ctx, reduce.NewReduceOp3(optional.Empty(), accumulator, combiner), split)
}

func (r *ReferencePipeline) Max(ctx context.Context, comparator util.Comparator) optional.Optional {
	return r.Reduce(ctx, binary.MaxBy(comparator), nil)
}

func (r *ReferencePipeline) Min(ctx context.Context, comparator util.Comparator) optional.Optional {
	return r.Reduce(ctx, binary.MinBy(comparator), nil)
}

func (r *ReferencePipeline) Count(ctx context.Context, comparator util.Comparator) int {
	return len(r.upstreams)
}

func (r *ReferencePipeline) AnyMatch(ctx context.Context, predicate predicate.Predicater) bool {
	split := spliterator.NewSliceSpliterator2(spliterator.CharacteristicTODO, r.upstreams...)
	return r.Evaluate(ctx, match.NewMatchOp2(match.KindAny, predicate), split).IsPresent()
}

func (r *ReferencePipeline) AllMatch(ctx context.Context, predicate predicate.Predicater) bool {
	split := spliterator.NewSliceSpliterator2(spliterator.CharacteristicTODO, r.upstreams...)
	return r.Evaluate(ctx, match.NewMatchOp2(match.KindAll, predicate), split).IsPresent()
}

func (r *ReferencePipeline) NoneMatch(ctx context.Context, predicate predicate.Predicater) bool {
	split := spliterator.NewSliceSpliterator2(spliterator.CharacteristicTODO, r.upstreams...)
	return r.Evaluate(ctx, match.NewMatchOp2(match.KindNone, predicate), split).IsPresent()
}

func (r *ReferencePipeline) FindFirst(ctx context.Context, predicate predicate.Predicater) optional.Optional {
	split := spliterator.NewSliceSpliterator2(spliterator.CharacteristicTODO, r.upstreams...)

	return r.Evaluate(ctx, find.NewFindOp2(true, predicate), split)
}

func (r *ReferencePipeline) FindAny(ctx context.Context, predicate predicate.Predicater) optional.Optional {
	split := spliterator.NewSliceSpliterator2(spliterator.CharacteristicTODO, r.upstreams...)

	return r.Evaluate(ctx, find.NewFindOp2(false, predicate), split)
}

func (r *ReferencePipeline) Close() error {
	return nil
}
