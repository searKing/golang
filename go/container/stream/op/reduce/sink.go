// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reduce

import (
	"github.com/searKing/golang/go/container/stream/op/terminal"
	"github.com/searKing/golang/go/util/function/binary"
	"github.com/searKing/golang/go/util/object"
	"github.com/searKing/golang/go/util/optional"
)

/**
 * A type of {@code TerminalSink} that implements an associative reducing
 * operation on elements of type {@code T} and producing a result of type
 * {@code R}.
 *
 * @param <T> the type of input element to the combining operation
 * @param <R> the result type
 * @param <K> the type of the {@code AccumulatingSink}.
 */
type AccumulatingSink interface {
	Combine(other AccumulatingSink)
}

type ReducingSink struct {
	terminal.TODOSink
	seed  optional.Optional
	state optional.Optional

	reducer  binary.BiFunction
	combiner binary.BiFunction
}

func NewReducingSink(seed optional.Optional, reducer binary.BiFunction, combiner binary.BiFunction) *ReducingSink {
	object.RequireNonNull(reducer)
	sink := &ReducingSink{
		seed:     seed,
		reducer:  reducer,
		combiner: combiner,
	}
	sink.SetDerived(sink)
	return sink
}

func (sink *ReducingSink) Begin(size int) {
	sink.state = sink.seed
}

func (sink *ReducingSink) Accept(t interface{}) {
	if !sink.state.IsPresent() {
		sink.state = optional.Of(t)
		return
	}
	sink.state = optional.Of(sink.reducer.Apply(sink.state.Get(), t))
	return
}

func (sink *ReducingSink) Combine(other terminal.Sink) {
	otherOptional := other.Get()
	if !otherOptional.IsPresent() {
		return
	}

	if !sink.state.IsPresent() {
		sink.Accept(otherOptional.Get())
		return
	}
	if sink.combiner != nil {
		sink.state = optional.Of(sink.combiner.Apply(sink.state.Get(), otherOptional.Get()))
		return
	}
	sink.Accept(otherOptional.Get())
	return
}

func (sink *ReducingSink) Get() optional.Optional {
	return sink.state
}
