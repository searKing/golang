// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strings_test

import (
	strings_ "strings"
	"testing"

	"github.com/searKing/golang/go/strings"
)

type TrimNumberTest struct {
	input  string
	output []string
}

var (
	trimNumberTests = []TrimNumberTest{
		//{"21WhoAmI", []string{"21","WhoAmI"}},
		//{"2_1WhoAmI", []string{"2_1","WhoAmI"}},
		//{"0WhoAmI", []string{"0","WhoAmI"}},
		//{"000WhoAmI", []string{"000","WhoAmI"}},
		{"0x10WhoAmI", []string{"0x10", "WhoAmI"}},
		{"0x_1_0WhoAmI", []string{"0x_1_0", "WhoAmI"}},
		{"-0x10WhoAmI", []string{"-0x10", "WhoAmI"}},
		{"0377WhoAmI", []string{"0377", "WhoAmI"}},
		{"0_3_7_7WhoAmI", []string{"0_3_7_7", "WhoAmI"}},
		{"0o377WhoAmI", []string{"0o377", "WhoAmI"}},
		{"0o_3_7_7WhoAmI", []string{"0o_3_7_7", "WhoAmI"}},
		{"-0377WhoAmI", []string{"-0377", "WhoAmI"}},
		{"-0o377WhoAmI", []string{"-0o377", "WhoAmI"}},
		{"0WhoAmI", []string{"0", "WhoAmI"}},
		{"000WhoAmI", []string{"000", "WhoAmI"}},
		{"0x10WhoAmI", []string{"0x10", "WhoAmI"}},
		{"0377WhoAmI", []string{"0377", "WhoAmI"}},
		{"22WhoAmI", []string{"22", "WhoAmI"}},
		{"23WhoAmI", []string{"23", "WhoAmI"}},
		{"24WhoAmI", []string{"24", "WhoAmI"}},
		{"25WhoAmI", []string{"25", "WhoAmI"}},
		{"127WhoAmI", []string{"127", "WhoAmI"}},
		{"-21WhoAmI", []string{"-21", "WhoAmI"}},
		{"-22WhoAmI", []string{"-22", "WhoAmI"}},
		{"-23WhoAmI", []string{"-23", "WhoAmI"}},
		{"-24WhoAmI", []string{"-24", "WhoAmI"}},
		{"-25WhoAmI", []string{"-25", "WhoAmI"}},
		{"-128WhoAmI", []string{"-128", "WhoAmI"}},
		{"+21WhoAmI", []string{"+21", "WhoAmI"}},
		{"+22WhoAmI", []string{"+22", "WhoAmI"}},
		{"+23WhoAmI", []string{"+23", "WhoAmI"}},
		{"+24WhoAmI", []string{"+24", "WhoAmI"}},
		{"+25WhoAmI", []string{"+25", "WhoAmI"}},
		{"+127WhoAmI", []string{"+127", "WhoAmI"}},
		{"2.3WhoAmI", []string{"2.3", "WhoAmI"}},
		{"2.3e1WhoAmI", []string{"2.3e1", "WhoAmI"}},
		{"2.3e2WhoAmI", []string{"2.3e2", "WhoAmI"}},
		{"2.3p2WhoAmI", []string{"2.3p2", "WhoAmI"}},
		{"2.3p+2WhoAmI", []string{"2.3p+2", "WhoAmI"}},
		{"2.3p+66WhoAmI", []string{"2.3p+66", "WhoAmI"}},
		{"2.3p-66WhoAmI", []string{"2.3p-66", "WhoAmI"}},
		{"0x2.3p-66WhoAmI", []string{"0x2.3p-66", "WhoAmI"}},
		{"2_3.4_5WhoAmI", []string{"2_3.4_5", "WhoAmI"}},
	}
)

func TestTrimNumber(t *testing.T) {
	for n, test := range trimNumberTests {
		out := strings.SplitPrefixNumber(test.input)
		if strings_.Join(out, ",") != strings_.Join(test.output, ",") {
			t.Errorf("#%d: got %v; want %v", n, out, test.output)
		}
	}
}
