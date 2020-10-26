// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reflect_test

import (
	"reflect"
	"strings"
	"testing"

	reflect_ "github.com/searKing/golang/go/reflect"
)

func TestParseStructTag(t *testing.T) {
	test := []struct {
		name    string
		tag     string
		exp     []reflect_.SubStructTag
		invalid bool
	}{
		{
			name: "empty tag",
			tag:  "",
		},
		{
			name:    "tag with one key (invalid)",
			tag:     "json",
			invalid: true,
		},
		{
			name: "tag with one key (valid)",
			tag:  `json:""`,
			exp: []reflect_.SubStructTag{
				{
					Key: "json",
				},
			},
		},
		{
			name: "tag with one key and dash name",
			tag:  `json:"-"`,
			exp: []reflect_.SubStructTag{
				{
					Key: "json", Name: "-",
				},
			},
		},
		{
			name: "tag with key and name",
			tag:  `json:"foo"`,
			exp: []reflect_.SubStructTag{
				{
					Key: "json", Name: "foo",
				},
			},
		},
		{
			name: "tag with key, name and option",
			tag:  `json:"foo,omitempty"`,
			exp: []reflect_.SubStructTag{
				{
					Key: "json", Name: "foo", Options: []string{"omitempty"},
				},
			},
		},
		{
			name: "tag with multiple keys",
			tag:  `json:"" yaml:""`,
			exp: []reflect_.SubStructTag{
				{Key: "json"},
				{Key: "yaml"},
			},
		},
		{
			name: "tag with multiple keys and names",
			tag:  `json:"foo" yaml:"bar"`,
			exp: []reflect_.SubStructTag{
				{Key: "json", Name: "foo"},
				{Key: "yaml", Name: "bar"},
			},
		},
		{
			name: "tag with multiple keys, different names and options",
			tag:  `json:"foo,omitempty" yaml:"bar,omitempty"`,
			exp: []reflect_.SubStructTag{
				{Key: "json", Name: "foo", Options: []string{"omitempty"}},
				{Key: "yaml", Name: "bar", Options: []string{"omitempty"}},
			},
		},
		{
			name: "tag with multiple keys, different names and options",
			tag:  `json:"foo" yaml:"bar,omitempty" xml:"-"`,
			exp: []reflect_.SubStructTag{
				{Key: "json", Name: "foo"},
				{Key: "yaml", Name: "bar", Options: []string{"omitempty"}},
				{Key: "xml", Name: "-"},
			},
		},
		{
			name: "tag with quoted name",
			tag:  `json:"foo,bar:\"baz\""`,
			exp: []reflect_.SubStructTag{
				{
					Key: "json", Name: "foo", Options: []string{`bar:"baz"`},
				},
			},
		},
		{
			name:    "tag with trailing space",
			tag:     `json:"foo   " `,
			invalid: true,
		},
	}

	for i, ts := range test {
		t.Run(ts.name, func(t *testing.T) {
			tags, err := reflect_.ParseStructTag(ts.tag)
			invalid := err != nil

			if invalid != ts.invalid {
				t.Errorf("#%d, invalid case\n\twant: %+v\n\tgot : %+v\n\terr : %s", i, ts.invalid, invalid, err)
			}

			if invalid {
				return
			}

			for _, tag := range ts.exp {
				got, _ := tags.Get(tag.Key)
				if !reflect.DeepEqual(tag, got) {
					t.Errorf("#%d, parse\n\twant: %#v\n\tgot : %#v", i, tag, got)
				}
			}

			trimmedInput := strings.TrimSpace(ts.tag)
			got := tags.String()

			if len(trimmedInput) != len(got) {
				t.Errorf("#%d, parse string\n\twant: %#v\n\tgot : %#v", i, trimmedInput, got)
			}
		})
	}
}

func TestTags_Get(t *testing.T) {
	tag := `json:"foo,omitempty" yaml:"bar,omitempty"`

	tags, err := reflect_.ParseStructTag(tag)
	if err != nil {
		t.Fatal(err)
	}

	found, ok := tags.Get("json")
	if !ok {
		t.Fatalf("expect %q, go %q", "json", "")
	}

	t.Run("String", func(t *testing.T) {
		want := `json:"foo,omitempty"`
		if found.String() != want {
			t.Errorf("get\n\twant: %#v\n\tgot : %#v", want, found.String())
		}
	})
	t.Run("Value", func(t *testing.T) {
		want := `foo,omitempty`
		if found.Value() != want {
			t.Errorf("get\n\twant: %#v\n\tgot : %#v", want, found.Value())
		}
	})
}

func TestTags_Set(t *testing.T) {
	tag := `json:"foo,omitempty" yaml:"bar,omitempty"`

	tags, err := reflect_.ParseStructTag(tag)
	if err != nil {
		t.Fatal(err)
	}
	err = tags.Set(reflect_.SubStructTag{
		Key:     "json",
		Name:    "bar",
		Options: []string{},
	})
	if err != nil {
		t.Fatal(err)
	}

	found, ok := tags.Get("json")
	if !ok {
		t.Fatalf("expect %q, go %q", "json", "")
	}

	want := `json:"bar"`
	if found.String() != want {
		t.Errorf("set\n\twant: %#v\n\tgot : %#v", want, found.String())
	}
}

func TestTags_Set_Append(t *testing.T) {
	tag := `json:"foo,omitempty"`

	tags, err := reflect_.ParseStructTag(tag)
	if err != nil {
		t.Fatal(err)
	}

	err = tags.Set(reflect_.SubStructTag{
		Key:     "yaml",
		Name:    "bar",
		Options: []string{"omitempty"},
	})
	if err != nil {
		t.Fatal(err)
	}

	found, ok := tags.Get("yaml")
	if !ok {
		t.Fatalf("expect %q, go %q", "json", "")
	}

	want := `yaml:"bar,omitempty"`
	if found.String() != want {
		t.Errorf("set append\n\twant: %#v\n\tgot : %#v", want, found.String())
	}

	wantFull := `json:"foo,omitempty" yaml:"bar,omitempty"`
	got := tags.String()
	if len(got) != len(wantFull) {
		t.Errorf("set append\n\twant: %#v\n\tgot : %#v", wantFull, got)
	}
}

func TestTags_Set_KeyDoesNotExist(t *testing.T) {
	tag := `json:"foo,omitempty" yaml:"bar,omitempty"`

	tags, err := reflect_.ParseStructTag(tag)
	if err != nil {
		t.Fatal(err)
	}

	err = tags.Set(reflect_.SubStructTag{
		Key:     "",
		Name:    "bar",
		Options: []string{},
	})
	if err == nil {
		t.Fatal("setting tag with a nonexisting key should error")
	}
}

func TestTags_Delete(t *testing.T) {
	tag := `json:"foo,omitempty" yaml:"bar,omitempty" xml:"-"`

	tags, err := reflect_.ParseStructTag(tag)
	if err != nil {
		t.Fatal(err)
	}

	tags.Delete("xml")
	if len(tags.Keys()) != 2 {
		t.Fatalf("tag length should be 2, have %d", len(tags.Keys()))
	}

	found, ok := tags.Get("json")
	if !ok {
		t.Fatalf("expect %q, go %q", "json", "")
	}

	want := `json:"foo,omitempty"`
	if found.String() != want {
		t.Errorf("delete\n\twant: %#v\n\tgot : %#v", want, found.String())
	}

	wantFull := `json:"foo,omitempty" yaml:"bar,omitempty"`
	if len(tags.String()) != len(wantFull) {
		t.Errorf("delete\n\twant: %#v\n\tgot : %#v", wantFull, tags.String())
	}
}

func TestTags_DeleteOptions(t *testing.T) {
	tag := `json:"foo,omitempty" yaml:"bar,omitempty,omitempty" xml:"-"`

	tags, err := reflect_.ParseStructTag(tag)
	if err != nil {
		t.Fatal(err)
	}

	tags.DeleteOptions("json", "omitempty")

	want := `json:"foo" yaml:"bar,omitempty,omitempty" xml:"-"`
	if len(tags.String()) != len(want) {
		t.Errorf("delete option\n\twant: %#v\n\tgot : %#v", want, tags.String())
	}

	tags.DeleteOptions("yaml", "omitempty")
	want = `json:"foo" yaml:"bar" xml:"-"`
	if len(tags.String()) != len(want) {
		t.Errorf("delete option\n\twant: %#v\n\tgot : %#v", want, tags.String())
	}
}

func TestTags_AddOption(t *testing.T) {
	tag := `json:"foo" yaml:"bar,omitempty" xml:"-"`

	tags, err := reflect_.ParseStructTag(tag)
	if err != nil {
		t.Fatal(err)
	}

	tags.AddOptions("json", "omitempty")

	want := `json:"foo,omitempty" yaml:"bar,omitempty" xml:"-"`
	if len(tags.String()) != len(want) {
		t.Errorf("add options\n\twant: %#v\n\tgot : %#v", want, tags.String())
	}

	// this shouldn't change anything
	tags.AddOptions("yaml", "omitempty")

	want = `json:"foo,omitempty" yaml:"bar,omitempty" xml:"-"`
	if len(tags.String()) != len(want) {
		t.Errorf("add options\n\twant: %#v\n\tgot : %#v", want, tags.String())
	}

	// this should append to the existing
	tags.AddOptions("yaml", "omitempty", "flatten")
	want = `json:"foo,omitempty" yaml:"bar,omitempty,flatten" xml:"-"`
	if len(tags.String()) != len(want) {
		t.Errorf("add options\n\twant: %#v\n\tgot : %#v", want, tags.String())
	}
}

func TestTags_String(t *testing.T) {
	tag := `json:"foo" yaml:"bar,omitempty" xml:"-"`

	tags, err := reflect_.ParseStructTag(tag)
	if err != nil {
		t.Fatal(err)
	}

	got := tags.String()
	if len(got) != len(tag) {
		t.Errorf("string\n\twant: %#v\n\tgot : %#v", tag, got)
	}
}

func TestTags_OrderedString(t *testing.T) {
	tag := `json:"foo" yaml:"bar,omitempty" xml:"-"`

	tags, err := reflect_.ParseStructTag(tag)
	if err != nil {
		t.Fatal(err)
	}

	got := tags.OrderedString()
	want := `json:"foo" yaml:"bar,omitempty" xml:"-"`

	if got != want {
		t.Errorf("string\n\twant: %#v\n\tgot : %#v", want, got)
	}
}

func TestTags_SortedString(t *testing.T) {
	tag := `json:"foo" yaml:"bar,omitempty" xml:"-"`

	tags, err := reflect_.ParseStructTag(tag)
	if err != nil {
		t.Fatal(err)
	}

	got := tags.SortedString()
	want := `json:"foo" xml:"-" yaml:"bar,omitempty"`

	if got != want {
		t.Errorf("string\n\twant: %#v\n\tgot : %#v", want, got)
	}
}

func TestTags_AstString(t *testing.T) {
	tag := `json:"foo" yaml:"bar,omitempty" xml:"-"`

	tags, err := reflect_.ParseStructTag(tag)
	if err != nil {
		t.Fatal(err)
	}

	got := tags.AstString()
	want := "`json:\"foo\" yaml:\"bar,omitempty\" xml:\"-\"`"
	if len(got) != len(want) {
		t.Errorf("string\n\twant: %#v\n\tgot : %#v", want, got)
	}
}
