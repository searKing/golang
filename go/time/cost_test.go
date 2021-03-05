// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time_test

import (
	"sort"
	"testing"
	"time"

	time_ "github.com/searKing/golang/go/time"
)

func TestCostTick_Costs(t *testing.T) {
	tests := []struct {
		msg       string
		repeat    int
		wantCosts int
	}{
		{
			msg:       "",
			repeat:    10,
			wantCosts: 10,
		},
	}

	for _, tt := range tests {
		var cost time_.CostTick
		for i := 0; i < tt.repeat; i++ {
			cost.Tick("")
		}
		if len(cost.Costs()) != tt.wantCosts {
			t.Errorf("%s: expected %q got %q", tt.msg, tt.wantCosts, len(cost.Costs()))
		}
	}
}

func TestCostTick_Sort(t *testing.T) {
	tests := []struct {
		msg    string
		repeat int
		Lesser func(i time.Duration, j time.Duration) bool
	}{
		{
			msg:    "",
			repeat: 10,
			Lesser: nil,
		},
		{
			msg:    "",
			repeat: 10,
			Lesser: func(i time.Duration, j time.Duration) bool {
				return i < j
			},
		},
		{
			msg:    "",
			repeat: 10,
			Lesser: func(i time.Duration, j time.Duration) bool {
				return i > j
			},
		},
	}

	for _, tt := range tests {
		var cost time_.CostTick
		for i := 0; i < tt.repeat; i++ {
			cost.Tick("")
			time.Sleep(time.Duration(i) * time.Millisecond)
		}
		cost.Sort()

		if !sort.IsSorted(&cost) {
			t.Errorf("%s: expected sorted got %q", tt.msg, cost)
		}
	}
}
