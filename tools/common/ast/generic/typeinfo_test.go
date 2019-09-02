// Copyright 2019 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains tests for some of the internal functions.

package generic_test

import (
	ast2 "github.com/searKing/golang/tools/common/ast"
	"testing"
)

type NewTest struct {
	input  string
	output []TypeInfo
}

var (
	newTests = []NewTest{
		// No need for a test for the empty case; that's picked off before splitIntoRuns.
		// Single value.
		{"NumMap<int, *time.Time, *interface{}>", []TypeInfo{{
			Name:   "NumMap",
			Import: "",
			TemplateTypes: []TemplateType{{
				Type:   "int",
				Import: "",
			}, {
				Type:      "time.Time",
				Import:    "time",
				IsPointer: true,
			}, {
				Type:      "interface{}",
				IsPointer: true,
			}}}}},
		{"NumMap<a.b, *a.b.c>", []TypeInfo{{
			Name:   "NumMap",
			Import: "",
			TemplateTypes: []TemplateType{{
				Import: "a",
				Type:   "a.b",
			}, {
				Import:    "a.b",
				Type:      "b.c",
				IsPointer: true,
			}}}},
		}}
)

func TestNew(t *testing.T) {
Outer:
	for n, test := range newTests {
		runs := New(test.input)
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

type ParserTest struct {
	input  []ast2.Token
	output []TypeInfo
}

var (
	parserTests = []ParserTest{
		// No need for a test for the empty case; that's picked off before splitIntoRuns.
		// Single value.
		{[]ast2.Token{{
			Type:  ast2.TokenTypeName,
			Value: "NumMap",
		}, {
			Type:  ast2.TokenTypeParen,
			Value: "<",
		}, {
			Type:  ast2.TokenTypeName,
			Value: "int",
		}, {
			Type:  ast2.TokenTypeParen,
			Value: ",",
		}, {
			Type:  ast2.TokenTypeName,
			Value: "string",
		}, {
			Type:  ast2.TokenTypeParen,
			Value: ">",
		}}, []TypeInfo{{
			Name:   "NumMap",
			Import: "",
			TemplateTypes: []TemplateType{{
				Type:   "int",
				Import: "",
			}, {
				Type:   "string",
				Import: "",
			}}}}},
		{[]ast2.Token{{
			Type:  ast2.TokenTypeName,
			Value: "NumMap",
		}, {
			Type:  ast2.TokenTypeParen,
			Value: "<",
		}, {
			Type:  ast2.TokenTypeName,
			Value: "a.b",
		}, {
			Type:  ast2.TokenTypeParen,
			Value: ",",
		}, {
			Type:  ast2.TokenTypeName,
			Value: "a.b.c",
		}, {
			Type:  ast2.TokenTypeParen,
			Value: ">",
		}}, []TypeInfo{{
			Name:   "NumMap",
			Import: "",
			TemplateTypes: []TemplateType{{
				Import: "a",
				Type:   "a.b",
			}, {
				Import: "a.b",
				Type:   "b.c",
			}}}},
		}}
)

func TestParserTests(t *testing.T) {
Outer:
	for n, test := range parserTests {
		runs := Parser(test.input)
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
