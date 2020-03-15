// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package spliterator_test

import (
	"context"
	"testing"

	"github.com/searKing/golang/go/util/function/consumer"
	"github.com/searKing/golang/go/util/spliterator"
)

type SpliteratorForEachRemainingTests struct {
	input  []interface{}
	output []string
}

var spliteratorForEachRemainingTests = []SpliteratorForEachRemainingTests{
	{
		input:  []interface{}{"1", "2", "3", "4"},
		output: []string{"1", "2", "3", "4"},
	},
}

func TestSliceSpliterator_ForEachRemaining(t *testing.T) {
	for n, test := range spliteratorForEachRemainingTests {
		split := spliterator.NewSliceSpliterator2(spliterator.CharacteristicTODO, test.input...)

		var gots []string
		split.ForEachRemaining(context.Background(), consumer.ConsumerFunc(func(t interface{}) {
			gots = append(gots, t.(string))
		}))
		for i, got := range gots {
			if got != test.output[i] {
				t.Errorf("#%d[%d]: %v: got %q; expected %q", n, i, test.input, got, test.output[i])
			}
		}
	}
}

type SpliteratorTrySplitTests struct {
	input []interface{}
	ls    []string
	rs    []string
}

var spliteratorTrySplitTests = []SpliteratorTrySplitTests{
	{
		input: []interface{}{"1", "2", "3", "4"},
		ls:    []string{"1", "2"},
		rs:    []string{"3", "4"},
	},
}

func TestSliceSpliterator_TrySplit(t *testing.T) {
	for n, test := range spliteratorTrySplitTests {
		split := spliterator.NewSliceSpliterator2(spliterator.CharacteristicTODO, test.input...)

		var lsGots []string
		if ls := split.TrySplit(); ls != nil {
			ls.ForEachRemaining(context.Background(), consumer.ConsumerFunc(func(t interface{}) {
				lsGots = append(lsGots, t.(string))
			}))
			for i, got := range lsGots {
				if got != test.ls[i] {
					t.Errorf("ls #%d[%d]: %v: got %q; expected %q", n, i, test.input, got, test.ls[i])
				}
			}
		}
		var rsGots []string
		split.ForEachRemaining(context.Background(), consumer.ConsumerFunc(func(t interface{}) {
			rsGots = append(rsGots, t.(string))
		}))
		for i, got := range rsGots {
			if got != test.rs[i] {
				t.Errorf("rs #%d[%d]: %v: got %q; expected %q", n, i, test.input, got, test.rs[i])
			}
		}
	}
}

type SpliteratorTryAdvanceTests struct {
	input  []interface{}
	output []string
}

var spliteratorTryAdvanceTests = []SpliteratorTryAdvanceTests{
	{
		input:  []interface{}{"1", "2", "3", "4"},
		output: []string{"1", "2", "3", "4"},
	},
}

func TestSliceSpliterator_TryAdvance(t *testing.T) {
	for n, test := range spliteratorForEachRemainingTests {
		split := spliterator.NewSliceSpliterator2(spliterator.CharacteristicTODO, test.input...)

		var gots []string
		for split.TryAdvance(context.Background(), consumer.ConsumerFunc(func(t interface{}) {
			gots = append(gots, t.(string))
		})) {
		}

		for i, got := range gots {
			if got != test.output[i] {
				t.Errorf("#%d[%d]: %v: got %q; expected %q", n, i, test.input, got, test.output[i])
			}
		}
	}
}

type SpliteratorEstimateSizeTests struct {
	input  []interface{}
	output []int
}

var spliteratorEstimateSizeTests = []SpliteratorEstimateSizeTests{
	{
		input:  []interface{}{"1", "2", "3", "4"},
		output: []int{4, 3, 2, 1},
	},
}

func TestSliceSpliterator_EstimateSize(t *testing.T) {
	for n, test := range spliteratorEstimateSizeTests {
		split := spliterator.NewSliceSpliterator2(spliterator.CharacteristicTODO, test.input...)

		var gots []int
		for split.TryAdvance(context.Background(), consumer.ConsumerFunc(func(t interface{}) {
			gots = append(gots, split.EstimateSize())
		})) {
		}

		for i, got := range gots {
			if got != test.output[i] {
				t.Errorf("#%d[%d]: %v: got %q; expected %q", n, i, test.input, got, test.output[i])
			}
		}
	}
}
