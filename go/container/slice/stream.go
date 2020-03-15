// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slice

import (
	"github.com/searKing/golang/go/util/optional"
)

type Stream struct {
	s interface{}
}

func NewStream() *Stream {
	return &Stream{}
}
func (stream *Stream) WithSlice(s interface{}) *Stream {
	stream.s = s
	return stream
}

func (stream *Stream) Value() interface{} {
	return stream.s
}

func (stream *Stream) Filter(f func(interface{}) bool) *Stream {
	return stream.WithSlice(FilterFunc(stream.s, f))
}

func (stream *Stream) Map(f func(interface{}) interface{}) *Stream {
	return stream.WithSlice(MapFunc(stream.s, f))
}

func (stream *Stream) Distinct(f func(interface{}, interface{}) int) *Stream {
	return stream.WithSlice(DistinctFunc(stream.s, f))
}

func (stream *Stream) Sorted(f func(interface{}, interface{}) int) *Stream {
	return stream.WithSlice(SortedFunc(stream.s, f))
}

func (stream *Stream) Peek(f func(interface{})) *Stream {
	return stream.WithSlice(PeekFunc(stream.s, f))
}

func (stream *Stream) Limit(maxSize int) *Stream {
	return stream.WithSlice(LimitFunc(stream.s, maxSize))
}

func (stream *Stream) Skip(n int) *Stream {
	return stream.WithSlice(SkipFunc(stream.s, n))
}

func (stream *Stream) TakeWhile(f func(interface{}) bool) *Stream {
	return stream.WithSlice(TakeWhileFunc(stream.s, f))
}

func (stream *Stream) TakeUntil(f func(interface{}) bool) *Stream {
	return stream.WithSlice(TakeUntilFunc(stream.s, f))
}

func (stream *Stream) DropWhile(f func(interface{}) bool) *Stream {
	return stream.WithSlice(DropWhileFunc(stream.s, f))
}

func (stream *Stream) DropUntil(f func(interface{}) bool) *Stream {
	return stream.WithSlice(DropUntilFunc(stream.s, f))
}

func (stream *Stream) ForEach(f func(interface{})) {
	ForEachFunc(stream.s, f)
}

func (stream *Stream) ForEachOrdered(f func(interface{})) {
	ForEachOrderedFunc(stream.s, f)
}

func (stream *Stream) ToSlice(ifStringAsRune ...bool) interface{} {
	return ToSliceFunc(stream.s)
}

func (stream *Stream) Reduce(f func(left, right interface{}) interface{}) *optional.Optional {
	return optional.OfNillable(ReduceFunc(stream.s, f))
}

func (stream *Stream) Min(f func(interface{}, interface{}) int) *optional.Optional {
	return optional.OfNillable(MinFunc(stream.s, f))
}

func (stream *Stream) Max(f func(interface{}, interface{}) int) *optional.Optional {
	return optional.OfNillable(MaxFunc(stream.s, f))
}

func (stream *Stream) Count(ifStringAsRune ...bool) int {
	return CountFunc(stream.s)
}

func (stream *Stream) AnyMatch(f func(interface{}) bool) bool {
	return AnyMatchFunc(stream.s, f)
}

func (stream *Stream) AllMatch(f func(interface{}) bool) bool {
	return AllMatchFunc(stream.s, f)
}

func (stream *Stream) NoneMatch(f func(interface{}) bool) bool {
	return NoneMatchFunc(stream.s, f)
}

func (stream *Stream) FindFirst(f func(interface{}) bool) *optional.Optional {
	return optional.OfNillable(FindFirstFunc(stream.s, f))
}

func (stream *Stream) FindFirstIndex(f func(interface{}) bool) int {
	return FindFirstIndexFunc(stream.s, f)
}

func (stream *Stream) FindAny(f func(interface{}) bool) *optional.Optional {
	return optional.OfNillable(FindAnyFunc(stream.s, f))
}

func (stream *Stream) FindAnyIndex(f func(interface{}) bool) int {
	return FindAnyIndexFunc(stream.s, f)
}

func (stream *Stream) Empty(ifStringAsRune ...bool) interface{} {
	return EmptyFunc(stream.s)
}

func (stream *Stream) Of(ifStringAsRune ...bool) *Stream {
	return stream.WithSlice(Of(stream.s))
}

func (stream *Stream) Concat(s2 *Stream) *Stream {
	return stream.WithSlice(ConcatFunc(stream.s, s2.s))
}
func (stream *Stream) ConcatWithValue(v interface{}) *Stream {
	return stream.Concat(NewStream().WithSlice(v))
}

//grammar surger for count
func (stream *Stream) Size() int {
	return stream.Count()
}

func (stream *Stream) Length() int {
	return stream.Count()
}
