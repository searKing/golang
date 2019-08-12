package strings_test

import (
	strings_ "github.com/searKing/golang/go/strings"
	"strings"
	"testing"
)

type CamelCaseTest struct {
	input  string
	output string
}

var (
	camelCaseTests = []CamelCaseTest{
		{
			"name_2",
			"Name_2",
		},
		{
			"_my_field_name_2",
			"XMyFieldName_2",
		},
	}
)

func TestCamelCases(t *testing.T) {
	for n, test := range camelCaseTests {
		out := strings_.CamelCase(test.input)
		if !strings.EqualFold(out, test.output) {
			t.Errorf("#%d: got %v; expected %v", n, out, test.output)
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
			"Name_2",
		},
		{
			[]string{"", "my", "field", "name", "2"},
			"XMyFieldName_2",
		},
	}
)

func TestCamelCaseSlices(t *testing.T) {
	for n, test := range camelCaseSliceTests {
		out := strings_.CamelCaseSlice(test.input...)
		if !strings.EqualFold(out, test.output) {
			t.Errorf("#%d: got %v; expected %v", n, out, test.output)
		}
	}
}

var (
	smallCamelCaseTests = []CamelCaseTest{
		{
			"name_2",
			"name_2",
		},
		{
			"_my_field_name_2",
			"xMyFieldName_2",
		},
	}
)

func TestSmallCamelCases(t *testing.T) {
	for n, test := range smallCamelCaseTests {
		out := strings_.SmallCamelCase(test.input)
		if !strings.EqualFold(out, test.output) {
			t.Errorf("#%d: got %v; expected %v", n, out, test.output)
		}
	}
}

var (
	smallCamelCaseSliceTests = []CamelCaseSliceTest{
		{
			[]string{"name", "2"},
			"name_2",
		},
		{
			[]string{"", "my", "field", "name", "2"},
			"xMyFieldName_2",
		},
	}
)

func TestSmallCamelCaseSlices(t *testing.T) {
	for n, test := range smallCamelCaseSliceTests {
		out := strings_.SmallCamelCaseSlice(test.input...)
		if !strings.EqualFold(out, test.output) {
			t.Errorf("#%d: got %v; expected %v", n, out, test.output)
		}
	}
}
