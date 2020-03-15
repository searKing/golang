// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flag_test

import (
	"flag"
	"fmt"
	"sort"
	"testing"

	flag_ "github.com/searKing/golang/go/flag"
	"github.com/searKing/golang/go/util/object"
)

func TestEverything(t *testing.T) {
	ResetForTesting(nil)
	flag_.StringSlice("test_[]string", nil, "[]string value")

	m := make(map[string]*flag.Flag)
	desired := "0"
	desiredStringSlice := ""
	visitor := func(f *flag.Flag) {
		if len(f.Name) > 5 && f.Name[0:5] == "test_" {
			m[f.Name] = f
			ok := false
			switch {
			case f.Value.String() == desired:
				ok = true
			case f.Name == "test_[]string" && f.Value.String() == fmt.Sprintf("[%s]", desiredStringSlice):
				ok = true
			}
			if !ok {
				t.Error("Visit: bad value", f.Value.String(), "for", f.Name)
			}
		}
	}
	flag.VisitAll(visitor)
	if len(m) != 1 {
		t.Error("VisitAll misses some flags")
		for k, v := range m {
			t.Log(k, *v)
		}
	}
	m = make(map[string]*flag.Flag)
	flag.Visit(visitor)
	if len(m) != 0 {
		t.Errorf("Visit sees unset flags")
		for k, v := range m {
			t.Log(k, *v)
		}
	}
	// Now set all flags
	flag.Set("test_[]string", "one")
	flag.Set("test_[]string", "two")
	desiredStringSlice = `"one" "two"`
	flag.Visit(visitor)
	if len(m) != 1 {
		t.Error("Visit fails after set")
		for k, v := range m {
			t.Log(k, *v)
		}
	}
	// Now test they're visited in sort order.
	var flagNames []string
	flag.Visit(func(f *flag.Flag) { flagNames = append(flagNames, f.Name) })
	if !sort.StringsAreSorted(flagNames) {
		t.Errorf("flag names not sorted: %v", flagNames)
	}
}

func TestGet(t *testing.T) {
	ResetForTesting(nil)
	flag_.StringSlice("test_[]string", []string{"one", "two"}, "[]string value")

	visitor := func(f *flag.Flag) {
		if len(f.Name) > 5 && f.Name[0:5] == "test_" {
			g, ok := f.Value.(flag.Getter)
			if !ok {
				t.Errorf("Visit: value does not satisfy Getter: %T", f.Value)
				return
			}
			switch f.Name {
			case "test_[]string":
				ok = object.DeepEquals(g.Get(), []string{"one", "two"})
			}
			if !ok {
				t.Errorf("Visit: bad value %T(%v) for %s", g.Get(), g.Get(), f.Name)
			}
		}
	}
	flag.VisitAll(visitor)
}

func testParse(f *flag.FlagSet, t *testing.T) {
	if f.Parsed() {
		t.Error("f.Parse() = true before Parse")
	}
	stringSliceFlag := flag_.StringSliceWithFlagSet(f, "test_[]string", []string{"one", "two"}, "[]string value")
	stringSlice2Flag := flag_.StringSliceWithFlagSet(f, "test_[]string2", nil, "[]string value")
	stringSlice3Flag := flag_.StringSliceWithFlagSet(f, "test_[]string3", nil, "[]string value")
	extra := "one-extra-argument"
	args := []string{
		"-test_[]string", "1",
		"-test_[]string", "2",
		"-test_[]string", "3",
		"--test_[]string2", "one",
		extra,
	}
	if err := f.Parse(args); err != nil {
		t.Fatal(err)
	}
	if !f.Parsed() {
		t.Error("f.Parse() = false after Parse")
	}
	if !object.DeepEquals(*stringSliceFlag, []string{"1", "2", "3"}) {
		t.Error("[]string flag should be [1 2 3], is ", *stringSliceFlag)
	}
	if !object.DeepEquals(*stringSlice2Flag, []string{"one"}) {
		t.Error("[]string flag should be [one], is ", *stringSlice2Flag)
	}
	if object.IsNil(stringSlice3Flag) {
		t.Error("[]string flag should be [], is ", *stringSlice3Flag)
	}
	if len(f.Args()) != 1 {
		t.Error("expected one argument, got", len(f.Args()))
	} else if f.Args()[0] != extra {
		t.Errorf("expected argument %q got %q", extra, f.Args()[0])
	}
}

func TestParse(t *testing.T) {
	ResetForTesting(func() { t.Error("bad parse") })
	testParse(flag.CommandLine, t)
}

func TestFlagSetParse(t *testing.T) {
	testParse(flag.NewFlagSet("test", flag.ContinueOnError), t)
}
