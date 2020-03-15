// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package match

import (
	"github.com/searKing/golang/go/container/stream/op/terminal"
	"github.com/searKing/golang/go/util/function/predicate"
	"github.com/searKing/golang/go/util/object"
)

/**
 * Boolean specific terminal sink to avoid the boxing costs when returning
 * results.  Subclasses implement the shape-specific functionality.
 *
 * @param <T> The output type of the stream pipeline
 */
type booleanTerminalSink struct {
	terminal.TODOSink
	stop  bool
	value bool
}

func (sink *booleanTerminalSink) AcceptMatchKind(matchKind Kind) {
	sink.value = !matchKind.ShortCircuitResult
}

func (sink *booleanTerminalSink) GetAndClearState() bool {
	return sink.value
}

func (sink *booleanTerminalSink) CancellationRequested() bool {
	return sink.stop
}

type MatchSink struct {
	booleanTerminalSink

	matchKind Kind
	predicate predicate.Predicater
}

func NewMatchSink(matchKind Kind, predicate predicate.Predicater) *MatchSink {
	object.RequireNonNil(predicate)
	sink := &MatchSink{
		matchKind: matchKind,
		predicate: predicate,
	}
	sink.AcceptMatchKind(matchKind)
	sink.SetDerived(sink)
	return sink
}

func (sink *MatchSink) Accept(value interface{}) {
	if !sink.stop && sink.predicate.Test(value) == sink.matchKind.StopOnPredicateMatches {
		sink.stop = true
		sink.value = sink.matchKind.ShortCircuitResult
	}
}
