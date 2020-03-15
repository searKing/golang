// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package match

import (
	"context"

	"github.com/searKing/golang/go/container/stream/op/terminal"
	"github.com/searKing/golang/go/util/spliterator"
)

type MatchTask struct {
	terminal.TODOShortCircuitTask

	op MatchOp
}

func NewMatchTask(op MatchOp, spliterator spliterator.Spliterator) *MatchTask {
	matcher := &MatchTask{
		op: op,
	}
	matcher.WithSpliterator(spliterator)
	matcher.SetDerived(matcher)
	return matcher
}

func NewTaskFromParent(parent MatchTask, spliterator spliterator.Spliterator) *MatchTask {
	matcher := &MatchTask{
		op: parent.op,
	}
	matcher.WithSpliterator(spliterator)
	matcher.SetDerived(matcher)
	return matcher
}

func (task MatchTask) MakeChild(spliterator spliterator.Spliterator) terminal.Task {
	return NewTaskFromParent(task, spliterator)
}

func (task MatchTask) DoLeaf(ctx context.Context) terminal.Sink {
	result := terminal.WrapAndCopyInto(ctx, task.op.MakeSink(), task.GetSpliterator())
	b := result.(*booleanTerminalSink).GetAndClearState()
	if b == task.op.matchKind.ShortCircuitResult {
		task.ShortCircuit(result)
	}
	return nil
}
