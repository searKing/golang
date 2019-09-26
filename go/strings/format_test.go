package strings_test

import (
	strings_ "github.com/searKing/golang/go/strings"
	"strings"
	"testing"
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
	camelCaseTests = []CamelCaseTest{
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

func TestCamelCases(t *testing.T) {
	for n, test := range camelCaseTests {
		out := strings_.CamelCase(test.input, test.seps...)
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
	camelCaseSliceTests = []CamelCaseSliceTest{
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

func TestCamelCaseSlices(t *testing.T) {
	for n, test := range camelCaseSliceTests {
		out := strings_.CamelCaseSlice(test.input...)
		if strings.Compare(out, test.output) != 0 {
			t.Errorf("#%d: got %v; expected %v", n, out, test.output)
		}
	}
}

var (
	smallCamelCaseTests = []CamelCaseTest{
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

func TestSmallCamelCases(t *testing.T) {
	for n, test := range smallCamelCaseTests {
		out := strings_.SmallCamelCase(test.input, test.seps...)
		if strings.Compare(out, test.output) != 0 {
			t.Errorf("#%d: got %v; expected %v", n, out, test.output)
		}
	}
}

var (
	snakeCamelCaseTests = []CamelCaseTest{
		{
			"name_2",
			[]rune{'_'},
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
		out := strings_.SmallCamelCaseSlice(test.input...)
		if strings.Compare(out, test.output) != 0 {
			t.Errorf("#%d: got %v; expected %v", n, out, test.output)
		}
	}
}
