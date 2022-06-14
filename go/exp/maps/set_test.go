// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package maps_test

import (
	"testing"

	maps_ "github.com/searKing/golang/go/exp/maps"
	"golang.org/x/exp/maps"
)

var m1 = map[int]struct{}{1: {}, 2: {}, 4: {}, 8: {}}

func TestSet(t *testing.T) {
	mc := maps_.Set(maps.Keys(m1)...)
	if !maps.Equal(mc, m1) {
		t.Errorf("Set(%v) = %v, want %v", m1, mc, m1)
	}
	mc[16] = struct{}{}
	if maps.Equal(mc, m1) {
		t.Errorf("Equal(%v, %v) = true, want false", mc, m1)
	}
}
