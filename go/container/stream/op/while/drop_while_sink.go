// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package while

import (
	"github.com/searKing/golang/go/container/stream/op/terminal"
	"github.com/searKing/golang/go/util/function/predicate"
	"github.com/searKing/golang/go/util/optional"
)

/**
 * A specialization for a dropWhile sink.
 *
 * @param <T> the type of both input and output elements
 */
type DropWhileSink interface {
	terminal.Sink
	/**
	 * @return the could of elements that would have been dropped and
	 * instead were passed downstream.
	 */
	GetDropCount() int
}

type dropWhileSink struct {
	terminal.TODOSink
	dropCount int
	take      bool

	retainAndCountDroppedElements bool
	predicate                     predicate.Predicater
}

func NewDropWhileSink(retainAndCountDroppedElements bool, predicate predicate.Predicater) *dropWhileSink {
	c := &dropWhileSink{
		retainAndCountDroppedElements: retainAndCountDroppedElements,
		predicate:                     predicate,
	}
	c.SetDerived(c)
	return c
}

func (sink *dropWhileSink) Accept(t interface{}) {
	if !sink.take {
		sink.take = !sink.predicate.Test(t)
	}
	takeElement := sink.take

	// If ordered and element is dropped increment index
	// for possible future truncation
	if sink.retainAndCountDroppedElements && !takeElement {
		sink.dropCount++
	}

	// If ordered need to process element, otherwise
	// skip if element is dropped
	if sink.retainAndCountDroppedElements && takeElement {
		sink.TODOSink.Accept(t)
	}
}

func (sink *dropWhileSink) GetDropCount() int {
	return sink.dropCount
}

func (sink *dropWhileSink) Get() optional.Optional {
	return optional.Empty()
}
