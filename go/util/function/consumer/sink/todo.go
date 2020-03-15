// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sink

import (
	"github.com/searKing/golang/go/util/function/consumer"
)

type TODO struct {
	consumer.TODO
}

func (_ *TODO) Begin(size int) {
	return
}

func (_ *TODO) End() {
	return
}

func (_ *TODO) CancellationRequested() bool {
	return false
}
