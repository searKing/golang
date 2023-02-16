// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/searKing/golang/go/exp/types"
)

var testCasesPointerSlice = [][]string{
	{"a", "b", "c", "d", "e"},
	{"a", "b", "", "", "e"},
}

func TestPointerSlice(t *testing.T) {
	for idx, in := range testCasesPointerSlice {
		if in == nil {
			continue
		}
		out := types.PointerSlice(in)
		if e, a := len(out), len(in); e != a {
			t.Errorf("Unexpected len at idx %d", idx)
		}
		for i := range out {
			if e, a := in[i], *(out[i]); e != a {
				t.Errorf("Unexpected value at idx %d", idx)
			}
		}

		out2 := types.ValueSlice(out)
		if e, a := len(out2), len(in); e != a {
			t.Errorf("Unexpected len at idx %d", idx)
		}
		if e, a := in, out2; !reflect.DeepEqual(e, a) {
			t.Errorf("Unexpected value at idx %d", idx)
		}
	}
}

var testCasesValueSlice = [][]*string{
	{types.Pointer("a"), types.Pointer("b"), nil, types.Pointer("c")},
}

func TestValueSlice(t *testing.T) {
	for idx, in := range testCasesValueSlice {
		if in == nil {
			continue
		}
		out := types.ValueSlice(in)
		if e, a := len(out), len(in); e != a {
			t.Errorf("Unexpected len at idx %d", idx)
		}
		for i := range out {
			if in[i] == nil {
				if out[i] != "" {
					t.Errorf("Unexpected value at idx %d", idx)
				}
			} else {
				if e, a := *(in[i]), out[i]; e != a {
					t.Errorf("Unexpected value at idx %d", idx)
				}
			}
		}

		out2 := types.PointerSlice(out)
		if e, a := len(out2), len(in); e != a {
			t.Errorf("Unexpected len at idx %d", idx)
		}
		for i := range out2 {
			if in[i] == nil {
				if *(out2[i]) != "" {
					t.Errorf("Unexpected value at idx %d", idx)
				}
			} else {
				if e, a := *in[i], *out2[i]; e != a {
					t.Errorf("Unexpected value at idx %d", idx)
				}
			}
		}
	}
}

var testCasesPointerMap = []map[string]string{
	{"a": "1", "b": "2", "c": "3"},
}

func TestPointerMap(t *testing.T) {
	for idx, in := range testCasesPointerMap {
		if in == nil {
			continue
		}
		out := types.PointerMap(in)
		if e, a := len(out), len(in); e != a {
			t.Errorf("Unexpected len at idx %d", idx)
		}
		for i := range out {
			if e, a := in[i], *(out[i]); e != a {
				t.Errorf("Unexpected value at idx %d", idx)
			}
		}

		out2 := types.ValueMap(out)
		if e, a := len(out2), len(in); e != a {
			t.Errorf("Unexpected len at idx %d", idx)
		}
		if e, a := in, out2; !reflect.DeepEqual(e, a) {
			t.Errorf("Unexpected value at idx %d", idx)
		}
	}
}

var testCasesTimeSlice = [][]time.Time{
	{time.Now(), time.Now().AddDate(100, 0, 0)},
}

func TestTimeSlice(t *testing.T) {
	for idx, in := range testCasesTimeSlice {
		if in == nil {
			continue
		}
		out := types.PointerSlice(in)
		if e, a := len(out), len(in); e != a {
			t.Errorf("Unexpected len at idx %d", idx)
		}
		for i := range out {
			if e, a := in[i], *(out[i]); e != a {
				t.Errorf("Unexpected value at idx %d", idx)
			}
		}

		out2 := types.ValueSlice(out)
		if e, a := len(out2), len(in); e != a {
			t.Errorf("Unexpected len at idx %d", idx)
		}
		if e, a := in, out2; !reflect.DeepEqual(e, a) {
			t.Errorf("Unexpected value at idx %d", idx)
		}
	}
}

var testCasesTimeValueSlice = [][]*time.Time{}

func TestTimeValueSlice(t *testing.T) {
	for idx, in := range testCasesTimeValueSlice {
		if in == nil {
			continue
		}
		out := types.ValueSlice(in)
		if e, a := len(out), len(in); e != a {
			t.Errorf("Unexpected len at idx %d", idx)
		}
		for i := range out {
			if in[i] == nil {
				if !out[i].IsZero() {
					t.Errorf("Unexpected value at idx %d", idx)
				}
			} else {
				if e, a := *(in[i]), out[i]; e != a {
					t.Errorf("Unexpected value at idx %d", idx)
				}
			}
		}

		out2 := types.PointerSlice(out)
		if e, a := len(out2), len(in); e != a {
			t.Errorf("Unexpected len at idx %d", idx)
		}
		for i := range out2 {
			if in[i] == nil {
				if !(out2[i]).IsZero() {
					t.Errorf("Unexpected value at idx %d", idx)
				}
			} else {
				if e, a := in[i], out2[i]; e != a {
					t.Errorf("Unexpected value at idx %d", idx)
				}
			}
		}
	}
}

var testCasesTimeMap = []map[string]time.Time{
	{"a": time.Now().AddDate(-100, 0, 0), "b": time.Now()},
}

func TestTimeMap(t *testing.T) {
	for idx, in := range testCasesTimeMap {
		if in == nil {
			continue
		}
		out := types.PointerMap(in)
		if e, a := len(out), len(in); e != a {
			t.Errorf("Unexpected len at idx %d", idx)
		}
		for i := range out {
			if e, a := in[i], *(out[i]); e != a {
				t.Errorf("Unexpected value at idx %d", idx)
			}
		}

		out2 := types.ValueMap(out)
		if e, a := len(out2), len(in); e != a {
			t.Errorf("Unexpected len at idx %d", idx)
		}
		if e, a := in, out2; !reflect.DeepEqual(e, a) {
			t.Errorf("Unexpected value at idx %d", idx)
		}
	}
}

type TimeValueTestCase struct {
	in        int64
	outSecs   time.Time
	outMillis time.Time
}

var testCasesTimeValue = []TimeValueTestCase{
	{
		in:        int64(1501558289),
		outSecs:   time.Unix(1501558289, 0),
		outMillis: time.Unix(1501558, 289*1e6),
	},
	{
		in:        int64(1501558289001),
		outSecs:   time.Unix(1501558289001, 0),
		outMillis: time.Unix(1501558289, 1*1e6),
	},
}

func TestSecondsTimeValue(t *testing.T) {
	for idx, testCase := range testCasesTimeValue {
		//out := time_.UnixWithUnit(types.Value(&testCase.in), time.Second)
		out := time.Unix(types.Value(&testCase.in), 0)
		if e, a := testCase.outSecs, out; e != a {
			t.Errorf("Unexpected value for time value at %d", idx)
		}
	}
}

func TestMillisecondsTimeValue(t *testing.T) {
	for idx, testCase := range testCasesTimeValue {
		//out := time_.UnixWithUnit(types.Value(&testCase.in), time.Millisecond)
		out := time.UnixMilli(types.Value(&testCase.in))
		if e, a := testCase.outMillis, out; e != a {
			t.Errorf("Unexpected value for time value at %d", idx)
		}
	}
}
