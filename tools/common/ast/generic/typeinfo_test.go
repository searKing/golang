// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains tests for some of the internal functions.

package generic_test

import (
	"testing"

	"github.com/searKing/golang/tools/common/ast"
	"github.com/searKing/golang/tools/common/ast/generic"
)

type NewTest struct {
	input  string
	output []generic.TypeInfo
}

var (
	newTests = []NewTest{
		// No need for a test for the empty case; that's picked off before splitIntoRuns.
		// Single value.
		{"NumMap<int, *[][]*[]time.Time, *interface{}, map[[]map[string]int8][]map[int][]string>", []generic.TypeInfo{{
			Name:   "NumMap",
			Import: "",
			TemplateTypes: []generic.TemplateType{{
				Type:   "int",
				Import: "",
			}, {
				Type:       "time.Time",
				Import:     "time",
				IsPointer:  true,
				TypePrefix: "[][]*[]",
			}, {
				Type:      "interface{}",
				IsPointer: true,
			}, {
				Type:      "map[[]map[string]int8][]map[int][]string",
				IsPointer: false,
			}}}}},
		{"NumMap<a.b, *a.b/c.d>", []generic.TypeInfo{{
			Name:   "NumMap",
			Import: "",
			TemplateTypes: []generic.TemplateType{{
				Import: "a",
				Type:   "a.b",
			}, {
				Import:    "a.b/c",
				Type:      "c.d",
				IsPointer: true,
			}}}}}}
)

func TestNew(t *testing.T) {
Outer:
	for n, test := range newTests {
		runs := generic.New(test.input)
		if len(runs) != len(test.output) {
			t.Errorf("#%d: %v: got %d runs; expected %d", n, test.input, len(runs), len(test.output))
			continue
		}
		for i, run := range runs {
			if run.Name != test.output[i].Name {
				t.Errorf("#%d: .Name got %v; expected %v", n, run.Name, test.output[i].Name)
				continue Outer
			}
			if run.Import != test.output[i].Import {
				t.Errorf("#%d: .Import got %v; expected %v", n, run.Import, test.output[i].Import)
				continue Outer
			}

			if len(run.TemplateTypes) != len(test.output[i].TemplateTypes) {
				t.Errorf("#%d: len(.TemplateTypes) got %v; expected %v", n, len(run.TemplateTypes), len(test.output[i].TemplateTypes))
				continue Outer
			}

			for j, outTmpl := range test.output[i].TemplateTypes {
				runTmpl := run.TemplateTypes[j]
				if runTmpl.Type != outTmpl.Type {
					t.Errorf("#%d: .TemplateTypes[%d].Type got %v; expected %v", n, j, runTmpl.Type, outTmpl.Type)
					continue Outer
				}
				if runTmpl.Import != outTmpl.Import {
					t.Errorf("#%d: .TemplateTypes[%d].Import got %v; expected %v", n, j, runTmpl.Import, outTmpl.Import)
					continue Outer
				}
				if runTmpl.IsPointer != outTmpl.IsPointer {
					t.Errorf("#%d: .TemplateTypes[%d].IsPointer got %v; expected %v", n, j, runTmpl.IsPointer, outTmpl.IsPointer)
					continue Outer
				}
				if runTmpl.TypePrefix != outTmpl.TypePrefix {
					t.Errorf("#%d: .TemplateTypes[%d].Prefix got %v; expected %v", n, j, runTmpl.TypePrefix, outTmpl.TypePrefix)
					continue Outer
				}
			}
		}
	}
}

type ParserTest struct {
	input  []ast.Token
	output []generic.TypeInfo
}

var (
	parserTests = []ParserTest{
		// No need for a test for the empty case; that's picked off before splitIntoRuns.
		// Single value.
		{[]ast.Token{{
			Type:  ast.TokenTypeName,
			Value: "NumMap",
		}, {
			Type:  ast.TokenTypeParen,
			Value: "<",
		}, {
			Type:  ast.TokenTypeName,
			Value: "int",
		}, {
			Type:  ast.TokenTypeParen,
			Value: ",",
		}, {
			Type:  ast.TokenTypeName,
			Value: "string",
		}, {
			Type:  ast.TokenTypeParen,
			Value: ">",
		}}, []generic.TypeInfo{{
			Name:   "NumMap",
			Import: "",
			TemplateTypes: []generic.TemplateType{{
				Type:   "int",
				Import: "",
			}, {
				Type:   "string",
				Import: "",
			}}}}},
		{[]ast.Token{{
			Type:  ast.TokenTypeName,
			Value: "NumMap",
		}, {
			Type:  ast.TokenTypeParen,
			Value: "<",
		}, {
			Type:  ast.TokenTypeName,
			Value: "a.b",
		}, {
			Type:  ast.TokenTypeParen,
			Value: ",",
		}, {
			Type:  ast.TokenTypeName,
			Value: "a.b/c.d",
		}, {
			Type:  ast.TokenTypeParen,
			Value: ">",
		}}, []generic.TypeInfo{{
			Name:   "NumMap",
			Import: "",
			TemplateTypes: []generic.TemplateType{{
				Import: "a",
				Type:   "a.b",
			}, {
				Import: "a.b/c",
				Type:   "c.d",
			}}}},
		}}
)

func TestParserTests(t *testing.T) {
Outer:
	for n, test := range parserTests {
		runs := generic.Parser(test.input)
		if len(runs) != len(test.output) {
			t.Errorf("#%d: %v: got %d runs; expected %d", n, test.input, len(runs), len(test.output))
			continue
		}
		for i, run := range runs {
			if run.Name != test.output[i].Name {
				t.Errorf("#%d: .Name got %v; expected %v", n, run.Name, test.output[i].Name)
				continue Outer
			}
			if run.Import != test.output[i].Import {
				t.Errorf("#%d: .Import got %v; expected %v", n, run.Import, test.output[i].Import)
				continue Outer
			}

			if len(run.TemplateTypes) != len(test.output[i].TemplateTypes) {
				t.Errorf("#%d: len(.TemplateTypes) got %v; expected %v", n, len(run.TemplateTypes), len(test.output[i].TemplateTypes))
				continue Outer
			}

			for j, outTmpl := range test.output[i].TemplateTypes {
				runTmpl := run.TemplateTypes[j]
				if runTmpl.Type != outTmpl.Type {
					t.Errorf("#%d: .TemplateTypes[%d].Type got %v; expected %v", n, j, runTmpl.Type, outTmpl.Type)
					continue Outer
				}
				if runTmpl.Import != outTmpl.Import {
					t.Errorf("#%d: .TemplateTypes[%d].Import got %v; expected %v", n, j, runTmpl.Import, outTmpl.Import)
					continue Outer
				}
				if runTmpl.IsPointer != outTmpl.IsPointer {
					t.Errorf("#%d: .TemplateTypes[%d].IsPointer got %v; expected %v", n, j, runTmpl.IsPointer, outTmpl.IsPointer)
					continue Outer
				}
			}
		}
	}
}
