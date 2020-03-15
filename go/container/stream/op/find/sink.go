// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package find

import (
	"github.com/searKing/golang/go/container/stream/op/terminal"
	"github.com/searKing/golang/go/util/optional"
)

type FindSink struct {
	terminal.TODOSink
	value optional.Optional
}

func NewFindSink() *FindSink {
	sink := &FindSink{}
	sink.SetDerived(sink)
	return sink
}

// stores only first value
func (sink *FindSink) Accept(value interface{}) {
	if sink.value.IsPresent() {
		return
	}
	sink.value = optional.Of(value)
}

// return true if found
func (sink *FindSink) CancellationRequested() bool {
	return sink.value.IsPresent()
}

func (sink *FindSink) Get() optional.Optional {
	return sink.value
}
