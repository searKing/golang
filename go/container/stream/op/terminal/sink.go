// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package terminal

import (
	"github.com/searKing/golang/go/util/function/consumer/sink"
	"github.com/searKing/golang/go/util/function/supplier"
	"github.com/searKing/golang/go/util/optional"
)

/**
 * A {@link Sink} which accumulates state as elements are accepted, and allows
 * a result to be retrieved after the computation is finished.
 *
 * @param <T> the type of elements to be accepted
 * @param <R> the type of the result
 *
 * @since 1.8
 */
type Sink interface {
	sink.Sink
	supplier.OptionalSupplier
}

type TODOSink struct {
	sink.TODO
}

func (sink *TODOSink) Get() optional.Optional {
	return optional.Empty()
}
