// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build cgo

package cgo_test

import (
	"testing"

	signal_cgo "github.com/searKing/golang/go/os/signal/cgo"
)

func TestDeprecatedSigInvokeChainCompatibility(t *testing.T) {
	signal_cgo.SetSigInvokeChain(1, 2, 3, 4)

	cancel := signal_cgo.SetSigInvokeChainCancelable(1, 2, 3, 4)
	if cancel == nil {
		t.Fatal("SetSigInvokeChainCancelable returned a nil cancel function")
	}
	cancel()
}
