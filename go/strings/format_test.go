// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strings_test

import (
	"strings"
	"testing"
	"unicode"

	strings_ "github.com/searKing/golang/go/strings"
	unicode_ "github.com/searKing/golang/go/unicode"
)

type TransformCaseTest struct {
	input  string
	seps   []rune
	f      func(r string) string
	output string
}

var (
	transformCaseTests = []TransformCaseTest{
		{
			"name____+++2",
			[]rune{'_', '+'},
			strings.ToUpper,
			"NAME2",
		},
		{
			"_my__field__Name2y_2age.gender",
			[]rune{'_', '.'},
			strings.ToUpper,
			"MYFIELDNAME2Y2AGEGENDER",
		},
		{
			"one__two_+_+three.four__",
			[]rune{'_', '.', '+'},
			strings.ToUpper,
			"ONETWOTHREEFOUR",
		},
		{
			"ONE__two_+_+three.four__",
			[]rune{'_', '.', '+'},
			strings.ToLower,
			"onetwothreefour",
		},
	}
)

func TestTransformCase(t *testing.T) {
	for n, test := range transformCaseTests {
		out := strings_.TransformCase(test.input, test.f, test.seps...)
		if strings.Compare(out, test.output) != 0 {
			t.Errorf("#%d: src %v; sep %s; got %v; expected %v", n, test.input, string(test.seps), out, test.output)
		}
	}
}

type CamelCaseTest struct {
	input  string
	seps   []rune
	output string
}

var (
	upperCamelCaseTests = []CamelCaseTest{
		{
			"name____+++2",
			[]rune{'_', '+'},
			"Name2",
		},
		{
			"_my__field__Name2y_2age.gender",
			[]rune{'_', '.'},
			"XMyFieldName2y2ageGender",
		},
		{
			"one__two_+_+three.four__",
			[]rune{'_', '.', '+'},
			"OneTwoThreeFour",
		},
	}
)

func TestUpperCamelCases(t *testing.T) {
	for n, test := range upperCamelCaseTests {
		out := strings_.UpperCamelCase(test.input, test.seps...)
		if strings.Compare(out, test.output) != 0 {
			t.Errorf("#%d: src %v; sep %s; got %v; expected %v", n, test.input, string(test.seps), out, test.output)
		}
	}
}

type CamelCaseSliceTest struct {
	input  []string
	output string
}

var (
	upperCamelCaseSliceTests = []CamelCaseSliceTest{
		{
			[]string{"name", "2"},
			"Name2",
		},
		{
			[]string{"", "my", "field", "name", "2"},
			"XMyFieldName2",
		},
	}
)

func TestUpperCamelCaseSlices(t *testing.T) {
	for n, test := range upperCamelCaseSliceTests {
		out := strings_.UpperCamelCaseSlice(test.input...)
		if strings.Compare(out, test.output) != 0 {
			t.Errorf("#%d: got %v; expected %v", n, out, test.output)
		}
	}
}

var (
	lowerCamelCaseTests = []CamelCaseTest{
		{
			"name_2",
			[]rune{'_'},
			"name2",
		},
		{
			"_my_field_name_2",
			[]rune{'_'},
			"xMyFieldName2",
		},
	}
)

func TestLowerCamelCases(t *testing.T) {
	for n, test := range lowerCamelCaseTests {
		out := strings_.LowerCamelCase(test.input, test.seps...)
		if strings.Compare(out, test.output) != 0 {
			t.Errorf("#%d: got %v; expected %v", n, out, test.output)
		}
	}
}

var (
	snakeCamelCaseTests = []CamelCaseTest{
		{
			"name_2",
			nil,
			"name_2",
		},
		{
			"_my_field_name_2",
			[]rune{'_'},
			"x_my_field_name_2",
		},
	}
)

func TestSnakeCamelCases(t *testing.T) {
	for n, test := range snakeCamelCaseTests {
		out := strings_.SnakeCase(test.input, test.seps...)
		if strings.Compare(out, test.output) != 0 {
			t.Errorf("#%d: got %v; expected %v", n, out, test.output)
		}
	}
}

var (
	darwinCamelCaseTests = []CamelCaseTest{
		{
			"name_2",
			[]rune{'_'},
			"Name_2",
		},
		{
			"_my_field_name_2",
			[]rune{'_'},
			"X_My_Field_Name_2",
		},
	}
)

func TestDarwinCamelCases(t *testing.T) {
	for n, test := range darwinCamelCaseTests {
		out := strings_.DarwinCase(test.input, test.seps...)
		if strings.Compare(out, test.output) != 0 {
			t.Errorf("#%d: got %v; expected %v", n, out, test.output)
		}
	}
}

var (
	kebabCamelCaseTests = []CamelCaseTest{
		{
			"name_2",
			[]rune{'_'},
			"name-2",
		},
		{
			"_my_field_name_2",
			[]rune{'_'},
			"x-my-field-name-2",
		},
	}
)

func TestKebabCamelCases(t *testing.T) {
	for n, test := range kebabCamelCaseTests {
		out := strings_.KebabCase(test.input, test.seps...)
		if strings.Compare(out, test.output) != 0 {
			t.Errorf("#%d: got %v; expected %v", n, out, test.output)
		}
	}
}

var (
	dotCamelCaseTests = []CamelCaseTest{
		{
			"name_2",
			[]rune{'_'},
			"name.2",
		},
		{
			"_my_field_name_2",
			[]rune{'_'},
			"x.my.field.name.2",
		},
	}
)

type studlyCapsCaseTest struct {
	input     string
	upperCase unicode.SpecialCase
	output    string
}

var (
	studlyCapsCaseTests = []studlyCapsCaseTest{
		{
			"abcdefghijklmnopqrstuvwxyz",
			unicode_.VowelCase(nil, func(r rune) rune {
				return unicode.ToUpper(r)
			}, nil),
			"AbcdEfghIjklmnOpqrstUvwxyz",
		},
		{
			"abcdefghijklmnopqrstuvwxyz",
			unicode_.ConsonantCase(nil, func(r rune) rune {
				return unicode.ToUpper(r)
			}, nil),
			"aBCDeFGHiJKLMNoPQRSTuVWXYZ",
		},
		{
			"the quick brown fox jumps over the lazy dog",
			unicode_.VowelCase(nil, func(r rune) rune {
				return unicode.ToUpper(r)
			}, nil),
			"thE qUIck brOwn fOx jUmps OvEr thE lAzy dOg",
		},
		{
			"the quick brown fox jumps over the lazy dog",
			unicode_.ConsonantCase(nil, func(r rune) rune {
				return unicode.ToUpper(r)
			}, nil),
			"THe QuiCK BRoWN FoX JuMPS oVeR THe LaZY DoG",
		},
		{
			"name_2",
			unicode_.VowelCase(nil, func(r rune) rune {
				return unicode.ToUpper(r)
			}, nil),
			"nAmE_2",
		},
		{
			"_i_love_you_2",
			unicode_.VowelCase(nil, func(r rune) rune {
				return unicode.ToUpper(r)
			}, nil),
			"_I_lOvE_yOU_2",
		},
	}
)

func TestStudlyCapsCases(t *testing.T) {
	for n, test := range studlyCapsCaseTests {
		out := strings_.StudlyCapsCase(test.upperCase, test.input)
		if strings.Compare(out, test.output) != 0 {
			t.Errorf("#%d: got %v; expected %v", n, out, test.output)
		}
	}
}

type studlyCapsVowelUpperCaseTest struct {
	input  string
	output string
}

var (
	studlyCapsVowelUpperCaseTests = []studlyCapsVowelUpperCaseTest{
		{
			"abcdefghijklmnopqrstuvwxyz",
			"AbcdEfghIjklmnOpqrstUvwxyz",
		},
		{
			"the quick brown fox jumps over the lazy dog",
			"thE qUIck brOwn fOx jUmps OvEr thE lAzy dOg",
		},
		{
			"name_2",
			"nAmE_2",
		},
		{
			"_i_love_you_2",
			"_I_lOvE_yOU_2",
		},
	}
)

func TestStudlyCapsVowelUpperCase(t *testing.T) {
	for n, test := range studlyCapsVowelUpperCaseTests {
		out := strings_.StudlyCapsVowelUpperCase(test.input)
		if strings.Compare(out, test.output) != 0 {
			t.Errorf("#%d: got %v; expected %v", n, out, test.output)
		}
	}
}

var (
	studlyCapsConsonantUpperCaseTests = []studlyCapsVowelUpperCaseTest{
		{
			"abcdefghijklmnopqrstuvwxyz",
			"aBCDeFGHiJKLMNoPQRSTuVWXYZ",
		},
		{
			"the quick brown fox jumps over the lazy dog",
			"THe QuiCK BRoWN FoX JuMPS oVeR THe LaZY DoG",
		},
		{
			"name_2",
			"NaMe_2",
		},
		{
			"_i_love_you_2",
			"_i_LoVe_You_2",
		},
	}
)

func TestStudlyCapsConsonantUpperCase(t *testing.T) {
	for n, test := range studlyCapsConsonantUpperCaseTests {
		out := strings_.StudlyCapsConsonantUpperCase(test.input)
		if strings.Compare(out, test.output) != 0 {
			t.Errorf("#%d: got %v; expected %v", n, out, test.output)
		}
	}
}

func TestDotCamelCases(t *testing.T) {
	for n, test := range dotCamelCaseTests {
		out := strings_.DotCase(test.input, test.seps...)
		if strings.Compare(out, test.output) != 0 {
			t.Errorf("#%d: got %v; expected %v", n, out, test.output)
		}
	}
}

var (
	smallCamelCaseSliceTests = []CamelCaseSliceTest{
		{
			[]string{"name", "2"},
			"name2",
		},
		{
			[]string{"", "my", "field", "name", "2"},
			"xMyFieldName2",
		},
	}
)

func TestSmallCamelCaseSlices(t *testing.T) {
	for n, test := range smallCamelCaseSliceTests {
		out := strings_.LowerCamelCaseSlice(test.input...)
		if strings.Compare(out, test.output) != 0 {
			t.Errorf("#%d: got %v; expected %v", n, out, test.output)
		}
	}
}
