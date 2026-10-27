// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package prettyjson_test

import (
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"
	"maps"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/searKing/golang/go/encoding/prettyjson"
)

func TestPrettyJsonMarshal(t *testing.T) {
	const longString = "The quick brown fox jumps over the lazy dog"
	const longUrl = "https://example.com/tests/1.html?foo=1&bar=baz&a=0&b=1&c=2&d=3#paragraph"
	var longBytes = []byte(longString)
	_ = longBytes
	var longSlice = []string{longString, "Crimson", "Red", "Ruby", "Maroon"}
	_ = longSlice
	var longMap = map[string]string{"0": "red", "1": "green", "2": "blue", "3": "white", "4": "black"}
	_ = longMap
	tests := []struct {
		data any
		opts []prettyjson.EncOptsOption
		want string
	}{
		// string
		{
			data: longString,
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateString(4)},
			want: `"The ...43 chars"`},
		{
			data: longString,
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateString(4), prettyjson.WithEncOptsOmitStatistics(true)},
			want: `"The "`},
		{
			data: longString,
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateString(1), prettyjson.WithEncOptsTruncateStringIfMoreThan(4)},
			want: `"T...43 chars"`},
		{
			data: longString,
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateString(4), prettyjson.WithEncOptsTruncateStringIfMoreThan(1)},
			want: `"The ...43 chars"`},
		{
			data: longString,
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateString(len(longString)), prettyjson.WithEncOptsTruncateStringIfMoreThan(1)},
			want: `"The quick brown fox jumps over the lazy dog"`},
		{
			data: longString,
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateString(1), prettyjson.WithEncOptsTruncateStringIfMoreThan(4), prettyjson.WithEncOptsOmitStatistics(true)},
			want: `"T"`},
		// url
		{
			data: longUrl,
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateString(1), prettyjson.WithEncOptsTruncateStringIfMoreThan(4)},
			want: `"h...72 chars"`}, // take url as string
		{
			data: longUrl,
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateString(1), prettyjson.WithEncOptsTruncateStringIfMoreThan(4), prettyjson.WithEncOptsTruncateUrl(true)},
			want: `"https://example.com/tests/1.html...72 chars,6Q9F"`}, // take url as url
		{
			data: longUrl,
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateString(1), prettyjson.WithEncOptsTruncateStringIfMoreThan(4),
				prettyjson.WithEncOptsForceLongUrl(true), prettyjson.WithEncOptsEscapeHTML(true)},
			want: `"https://example.com/tests/1.html?foo=1\u0026bar=baz\u0026a=0\u0026b=1\u0026c=2\u0026d=3#paragraph"`}, // take url as url
		{
			data: longUrl,
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateString(1), prettyjson.WithEncOptsTruncateStringIfMoreThan(4),
				prettyjson.WithEncOptsForceLongUrl(true), prettyjson.WithEncOptsEscapeHTML(false)},
			want: `"https://example.com/tests/1.html?foo=1&bar=baz&a=0&b=1&c=2&d=3#paragraph"`}, // take url as url
		//[]byte
		{
			data: longBytes,
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateBytes(4)},
			want: `"VGhlIA==...43 bytes"`},
		{
			data: longBytes,
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateBytes(4), prettyjson.WithEncOptsOmitStatistics(true)},
			want: `"VGhlIA=="`},
		{
			data: longBytes,
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateBytes(1), prettyjson.WithEncOptsTruncateBytesIfMoreThan(4)},
			want: `"VA==...43 bytes"`},
		{
			data: longBytes,
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateBytes(4), prettyjson.WithEncOptsTruncateBytesIfMoreThan(1)},
			want: `"VGhlIA==...43 bytes"`},
		{
			data: longBytes,
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateBytes(len("The quick brown fox jumps over the lazy dog")), prettyjson.WithEncOptsTruncateBytesIfMoreThan(1)},
			want: `"VGhlIHF1aWNrIGJyb3duIGZveCBqdW1wcyBvdmVyIHRoZSBsYXp5IGRvZw=="`},
		{
			data: longBytes,
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateBytes(1), prettyjson.WithEncOptsTruncateBytesIfMoreThan(4), prettyjson.WithEncOptsOmitStatistics(true)},
			want: `"VA=="`},
		// slice
		{
			data: longSlice,
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateSliceOrArray(1)},
			want: `["The quick brown fox jumps over the lazy dog","...5 elems"]`},
		{
			data: longSlice,
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateSliceOrArray(1), prettyjson.WithEncOptsTruncateSliceOrArrayIfMoreThan(4)},
			want: `["The quick brown fox jumps over the lazy dog","...5 elems"]`},
		{
			data: longSlice,
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateSliceOrArray(2), prettyjson.WithEncOptsTruncateSliceOrArrayIfMoreThan(4)},
			want: `["The quick brown fox jumps over the lazy dog","Crimson","...5 elems"]`},
		{
			data: longSlice,
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateSliceOrArray(2), prettyjson.WithEncOptsTruncateSliceOrArrayIfMoreThan(4),
				prettyjson.WithEncOptsOmitStatistics(true)},
			want: `["The quick brown fox jumps over the lazy dog","Crimson"]`},
		{
			data: longSlice,
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateSliceOrArray(4), prettyjson.WithEncOptsTruncateSliceOrArrayIfMoreThan(1)},
			want: `["The quick brown fox jumps over the lazy dog","Crimson","Red","Ruby","...5 elems"]`},
		{
			data: longSlice,
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateSliceOrArray(5), prettyjson.WithEncOptsTruncateSliceOrArrayIfMoreThan(1)},
			want: `["The quick brown fox jumps over the lazy dog","Crimson","Red","Ruby","Maroon"]`},
		{
			data: []int{1, 2, 3, 4, 5},
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateSliceOrArray(1), prettyjson.WithEncOptsTruncateSliceOrArrayIfMoreThan(4)},
			want: `[1,"...5 elems"]`},
		// map
		{
			data: longMap,
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateMap(1), prettyjson.WithEncOptsOmitStatistics(true)},
			want: `{"0":"red"}`},
		{
			data: longMap,
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateMap(1)},
			want: `{"0":"red","1...5 pairs":"green"}`},
		{
			data: longMap,
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateMap(1), prettyjson.WithEncOptsTruncateMapIfMoreThan(4)},
			want: `{"0":"red","1...5 pairs":"green"}`},
		{
			data: longMap,
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateMap(1), prettyjson.WithEncOptsTruncateMapIfMoreThan(5)},
			want: `{"0":"red","1":"green","2":"blue","3":"white","4":"black"}`},
		{
			data: longMap,
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateMap(3), prettyjson.WithEncOptsTruncateMapIfMoreThan(1)},
			want: `{"0":"red","1":"green","2":"blue","3...5 pairs":"white"}`},
		{
			data: longMap,
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateMap(3), prettyjson.WithEncOptsTruncateMapIfMoreThan(1),
				prettyjson.WithEncOptsOmitStatistics(true)},
			want: `{"0":"red","1":"green","2":"blue"}`},
		{
			data: longMap,
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateMap(4), prettyjson.WithEncOptsTruncateMapIfMoreThan(1)},
			want: `{"0":"red","1":"green","2":"blue","3":"white","4":"black"}`},
		{
			data: longMap,
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateMap(4), prettyjson.WithEncOptsTruncateMapIfMoreThan(1),
				prettyjson.WithEncOptsOmitStatistics(true)},
			want: `{"0":"red","1":"green","2":"blue","3":"white"}`},
		// nested struct
		{
			data: struct {
				ID            int
				Name          string
				LongName      string
				Colors        []string
				ColorById     map[string]string
				ColorEnumById map[string]int
				Url           string
				LongUrl       string
				Empty         struct {
					ID   int
					Name string
				}
			}{
				ID:            1,
				Name:          "Reds",
				LongName:      "The quick brown fox jumps over the lazy dog",
				Colors:        []string{"The quick brown fox jumps over the lazy dog", "Crimson", "Red", "Ruby", "Maroon"},
				ColorById:     map[string]string{"0": "red", "1": "green", "2": "blue", "3": "white", "4": "black"},
				ColorEnumById: map[string]int{"0": 0, "1": 1, "2": 2, "3": 3},
				Url:           "https://example.com/tests/1?foo=1&bar=baz",
				LongUrl:       "https://example.com/tests/1.html?foo=1&bar=baz&a=0&b=1&c=2&d=3#paragraph",
			},
			opts: []prettyjson.EncOptsOption{
				prettyjson.WithEncOptsTruncateString(10),
				prettyjson.WithEncOptsTruncateStringIfMoreThan(20),
				prettyjson.WithEncOptsTruncateSliceOrArray(2),
				prettyjson.WithEncOptsTruncateSliceOrArrayIfMoreThan(4),
				prettyjson.WithEncOptsTruncateMap(2),
				prettyjson.WithEncOptsTruncateMapIfMoreThan(4),
				prettyjson.WithEncOptsTruncateUrl(true),
				prettyjson.WithEncOptsEscapeHTML(false),
				prettyjson.WithEncOptsForceLongUrl(false),
				prettyjson.WithEncOptsOmitEmpty(true),
				prettyjson.WithEncOptsOmitStatistics(true)},
			want: `{"ID":1,"Name":"Reds","LongName":"The quick ","Colors":["The quick ","Crimson"],"ColorById":{"0":"red","1":"green"},"ColorEnumById":{"0":0,"1":1,"2":2,"3":3},"Url":"https://example.com/tests/1","LongUrl":"https://example.com/tests/1.html"}`},

		{
			data: struct {
				ID            int
				Name          string
				LongName      string
				Colors        []string
				ColorById     map[string]string
				ColorEnumById map[string]int
				Url           string
				LongUrl       string
				Empty         struct {
					ID   int
					Name string
				}
			}{
				ID:            1,
				Name:          "Reds",
				LongName:      "The quick brown fox jumps over the lazy dog",
				Colors:        []string{"The quick brown fox jumps over the lazy dog", "Crimson", "Red", "Ruby", "Maroon"},
				ColorById:     map[string]string{"0": "red", "1": "green", "2": "blue", "3": "white", "4": "black"},
				ColorEnumById: map[string]int{"0": 0, "1": 1, "2": 2, "3": 3},
				Url:           "https://example.com/tests/1?foo=1&bar=baz",
				LongUrl:       "https://example.com/tests/1.html?foo=1&bar=baz&a=0&b=1&c=2&d=3#paragraph",
			},
			opts: []prettyjson.EncOptsOption{
				prettyjson.WithEncOptsTruncateString(10),
				prettyjson.WithEncOptsTruncateStringIfMoreThan(20),
				prettyjson.WithEncOptsTruncateSliceOrArray(2),
				prettyjson.WithEncOptsTruncateSliceOrArrayIfMoreThan(4),
				prettyjson.WithEncOptsTruncateMap(2),
				prettyjson.WithEncOptsTruncateMapIfMoreThan(4),
				prettyjson.WithEncOptsTruncateUrl(true),
				prettyjson.WithEncOptsEscapeHTML(false),
				prettyjson.WithEncOptsForceLongUrl(false),
				prettyjson.WithEncOptsOmitEmpty(true)},
			want: `{"ID":1,"Name":"Reds","LongName":"The quick ...43 chars","Colors":["The quick ...43 chars","Crimson","...5 elems"],"ColorById":{"0":"red","1":"green","2...5 pairs":"blue"},"ColorEnumById":{"0":0,"1":1,"2":2,"3":3},"Url":"https://example.com/tests/1?foo=1&bar=baz","LongUrl":"https://example.com/tests/1.html...72 chars,6Q9F"}`},
	}
	for i, tt := range tests {
		got, err := prettyjson.Marshal(tt.data, tt.opts...)
		if err != nil {
			t.Errorf("#%d: Marshal(%v) error: %v", i, tt.data, err)
			continue
		}
		if tt.want != string(got) {
			t.Errorf("#%d: Marshal(%v) = `%v`, want `%s`", i, tt.data, string(got), tt.want)
		}
	}
}

type textUUID [4]byte

func (textUUID) MarshalText() ([]byte, error) { return []byte("uuid"), nil }

type joinedStrings []string

func (s joinedStrings) MarshalJSON() ([]byte, error) { return jsonv2.Marshal(strings.Join(s, ",")) }

type ptrJoinedStrings []string

func (s *ptrJoinedStrings) MarshalJSON() ([]byte, error) {
	return jsonv2.Marshal(strings.Join(*s, ","))
}

type namedMap map[string]int

func (namedMap) MarshalJSON() ([]byte, error) { return []byte(`"namedMap"`), nil }

func TestMarshal_ContainerTruncation(t *testing.T) {
	strs := []string{"a", "b", "c", "d"}
	tests := []struct {
		name string
		data any
		opts []prettyjson.EncOptsOption
		want string
	}{
		{
			name: "TextMarshaler array is not truncated",
			data: textUUID{1, 2, 3, 4},
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateSliceOrArray(2)},
			want: `"uuid"`,
		},
		{
			name: "Marshaler slice is not truncated",
			data: joinedStrings(strs),
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateSliceOrArray(2)},
			want: `"a,b,c,d"`,
		},
		{
			name: "pointer receiver Marshaler slice is not truncated",
			data: struct{ V ptrJoinedStrings }{V: ptrJoinedStrings(strs)},
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateSliceOrArray(2)},
			want: `{"V":"a,b,c,d"}`,
		},
		{
			name: "Marshaler map is not truncated",
			data: namedMap{"a": 1, "b": 2, "c": 3},
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateMap(1)},
			want: `"namedMap"`,
		},
		{
			name: "byte array stays base64",
			data: [4]byte{1, 2, 3, 4},
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateSliceOrArray(2)},
			want: `"AQIDBA=="`,
		},
		{
			name: "slice behind pointer is truncated",
			data: struct{ P *[]string }{P: &strs},
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateSliceOrArray(1)},
			want: `{"P":["a","...4 elems"]}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := prettyjson.Marshal(tt.data, tt.opts...)
			if err != nil {
				t.Fatalf("Marshal(%v) error: %v", tt.data, err)
			}
			if string(got) != tt.want {
				t.Errorf("Marshal(%v) = `%s`, want `%s`", tt.data, got, tt.want)
			}
		})
	}
}

// stringSet marshals as a sorted JSON array of its keys.
type stringSet map[string]struct{}

func (s stringSet) MarshalJSON() ([]byte, error) { return jsonv2.Marshal(slices.Sorted(maps.Keys(s))) }

// maskedHeaders marshals as a JSON object with the Authorization value masked.
type maskedHeaders map[string]string

func (h maskedHeaders) MarshalJSON() ([]byte, error) {
	m := maps.Clone(map[string]string(h))
	if _, ok := m["Authorization"]; ok {
		m["Authorization"] = "***"
	}
	return jsonv2.Marshal(m, jsonv2.Deterministic(true))
}

func TestMarshal_TruncateMarshalerOutput(t *testing.T) {
	long := strings.Repeat("x", 40)
	words := joinedStrings(slices.Repeat([]string{"abc"}, 10))
	on := prettyjson.WithEncOptsTruncateMarshalerOutput(true)
	tests := []struct {
		name string
		data any
		opts []prettyjson.EncOptsOption
		want string
	}{
		{
			name: "disabled by default",
			data: words,
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateString(5)},
			want: `"abc,abc,abc,abc,abc,abc,abc,abc,abc,abc"`,
		},
		{
			name: "string output follows string rules",
			data: words,
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateString(5), on},
			want: `"abc,a...39 chars"`,
		},
		{
			name: "array output follows slice rules",
			data: stringSet{"a": {}, "b": {}, "c": {}, "d": {}},
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateSliceOrArray(2), on},
			want: `["a","b","...4 elems"]`,
		},
		{
			name: "strings nested in array output are truncated",
			data: stringSet{long: {}, "b": {}},
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateString(5), on},
			want: `["b","xxxxx...40 chars"]`,
		},
		{
			name: "object output keeps members and masked values",
			data: maskedHeaders{"Authorization": "Bearer " + long, "Long": long, "A": "1"},
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateMap(1), prettyjson.WithEncOptsTruncateString(5), on},
			want: `{"A":"1","Authorization":"***","Long":"xxxxx...40 chars"}`,
		},
		{
			name: "byte array keeps its text",
			data: textUUID{1, 2, 3, 4},
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateString(2), prettyjson.WithEncOptsTruncateSliceOrArray(2), on},
			want: `"uuid"`,
		},
		{
			name: "plain slice is truncated as before",
			data: []string{"a", "b", "c", "d"},
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateSliceOrArray(1), on},
			want: `["a","...4 elems"]`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := prettyjson.Marshal(tt.data, tt.opts...)
			if err != nil {
				t.Fatalf("Marshal(%v) error: %v", tt.data, err)
			}
			if string(got) != tt.want {
				t.Errorf("Marshal(%v) = `%s`, want `%s`", tt.data, got, tt.want)
			}
		})
	}
}

type upperKey string

func (k upperKey) MarshalText() ([]byte, error) { return []byte(strings.ToUpper(string(k))), nil }

type ptrUpperKey string

func (k *ptrUpperKey) MarshalText() ([]byte, error) { return []byte(strings.ToUpper(string(*k))), nil }

type jsonKey string

func (k jsonKey) MarshalJSON() ([]byte, error) { return jsonv2.Marshal(strings.ToUpper(string(k))) }

func TestMarshal_TruncatedMapKeyNames(t *testing.T) {
	ka, kb, kc := upperKey("a"), upperKey("b"), upperKey("c")
	tests := []struct {
		name string
		data any
	}{
		{name: "value receiver MarshalText", data: map[upperKey]int{"a": 1, "b": 2, "c": 3}},
		{name: "pointer receiver MarshalText", data: map[ptrUpperKey]int{"a": 1, "b": 2, "c": 3}},
		{name: "pointer key", data: map[*upperKey]int{&ka: 1, &kb: 2, &kc: 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := prettyjson.Marshal(tt.data, prettyjson.WithEncOptsTruncateMap(1))
			if err != nil {
				t.Fatalf("Marshal(%v) error: %v", tt.data, err)
			}
			if want := `{"A":1,"B...3 pairs":2}`; string(got) != want {
				t.Errorf("Marshal(%v) = `%s`, want `%s`", tt.data, got, want)
			}
		})
	}
}

type namedString string

type textString string

func (s textString) MarshalText() ([]byte, error) { return []byte(s), nil }

type blob []byte

func TestMarshal_NamedTypes(t *testing.T) {
	long := strings.Repeat("x", 40)
	tests := []struct {
		name string
		data any
		opts []prettyjson.EncOptsOption
		want string
	}{
		{
			name: "named string",
			data: namedString(long),
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateString(5)},
			want: `"xxxxx...40 chars"`,
		},
		{
			name: "named string map key is not truncated",
			data: map[namedString]int{namedString(long): 1},
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateString(5)},
			want: `{"` + long + `":1}`,
		},
		{
			name: "named string with MarshalText is not truncated",
			data: textString(long),
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateString(5)},
			want: `"` + long + `"`,
		},
		{
			name: "named bytes",
			data: blob(long),
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateBytes(3)},
			want: `"eHh4...40 bytes"`,
		},
		{
			name: "raw JSON value is not truncated",
			data: jsontext.Value(`{"a":"` + long + `"}`),
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateBytes(3)},
			want: `{"a":"` + long + `"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := prettyjson.Marshal(tt.data, tt.opts...)
			if err != nil {
				t.Fatalf("Marshal(%v) error: %v", tt.data, err)
			}
			if string(got) != tt.want {
				t.Errorf("Marshal(%v) = `%s`, want `%s`", tt.data, got, tt.want)
			}
		})
	}
}

func TestAdvancedTruncateMarshalers_JoinsOtherMarshalers(t *testing.T) {
	data := struct {
		T time.Time
		S string
	}{S: strings.Repeat("x", 20)}
	marshalers := jsonv2.JoinMarshalers(
		jsonv2.MarshalFunc(func(time.Time) ([]byte, error) { return []byte(`"T"`), nil }),
		prettyjson.AdvancedTruncateMarshalers(prettyjson.WithEncOptsTruncateString(3)),
	)
	got, err := jsonv2.Marshal(data, jsonv2.WithMarshalers(marshalers))
	if err != nil {
		t.Fatalf("Marshal(%v) error: %v", data, err)
	}
	if want := `{"T":"T","S":"xxx...20 chars"}`; string(got) != want {
		t.Errorf("Marshal(%v) = `%s`, want `%s`", data, got, want)
	}
}

func TestMarshal_OmitZero(t *testing.T) {
	type inner struct{ X int }
	data := struct {
		Int   int
		Nil   []string
		Empty []string
		Arr   [2]int
		Inner inner
		Time  time.Time
		Str   string
	}{Empty: []string{}, Str: "s"}
	tests := []struct {
		name string
		opt  prettyjson.EncOptsOption
		want string
	}{
		{name: "OmitZero", opt: prettyjson.WithEncOptsOmitZero(true), want: `{"Empty":[],"Str":"s"}`},
		{name: "OmitEmpty keeps untagged non-nil empty slices", opt: prettyjson.WithEncOptsOmitEmpty(true), want: `{"Empty":[],"Str":"s"}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := prettyjson.Marshal(data, tt.opt)
			if err != nil {
				t.Fatalf("Marshal(%v) error: %v", data, err)
			}
			if string(got) != tt.want {
				t.Errorf("Marshal(%v) = `%s`, want `%s`", data, got, tt.want)
			}
		})
	}
}

func TestMarshalIndent(t *testing.T) {
	data := map[string]any{"a": []int{1, 2, 3}, "b": "<x>"}
	got, err := prettyjson.MarshalIndent(data, ">", "  ", prettyjson.WithEncOptsTruncateSliceOrArray(1))
	if err != nil {
		t.Fatalf("MarshalIndent(%v) error: %v", data, err)
	}
	want := `{
>  "a": [
>    1,
>    "...3 elems"
>  ],
>  "b": "<x>"
>}`
	if string(got) != want {
		t.Errorf("MarshalIndent(%v) = `%s`, want `%s`", data, got, want)
	}
}

func TestMarshal_MapWithUnsupportedKeyFallsBack(t *testing.T) {
	tests := []struct {
		name string
		data any
	}{
		{name: "float64 key", data: map[float64]int{1.5: 1, 2.5: 2, 3.5: 3}},
		{name: "bool key", data: map[bool]int{true: 1, false: 0}},
		{name: "MarshalJSON key", data: map[jsonKey]int{"a": 1, "b": 2, "c": 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want, wantErr := jsonv2.Marshal(tt.data, jsonv2.Deterministic(true))
			got, err := prettyjson.Marshal(tt.data, prettyjson.WithEncOptsTruncateMap(1))
			if (err != nil) != (wantErr != nil) {
				t.Fatalf("Marshal(%v) error = %v, want error %v", tt.data, err, wantErr)
			}
			if string(got) != string(want) {
				t.Errorf("Marshal(%v) = `%s`, want `%s`", tt.data, got, want)
			}
		})
	}
}

func TestMarshal_StringTruncation(t *testing.T) {
	longWord := strings.Repeat("x", 40)
	tests := []struct {
		name string
		data any
		opts []prettyjson.EncOptsOption
		want string
	}{
		{
			name: "statistics count runes",
			data: strings.Repeat("中", 10),
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateString(5)},
			want: `"中中中中中...10 chars"`,
		},
		{
			name: "url statistics count runes",
			data: "https://example.com/a?q=" + strings.Repeat("中", 20),
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateString(5), prettyjson.WithEncOptsTruncateUrl(true)},
			want: `"https://example.com/a...44 chars,1Q"`,
		},
		{
			name: "url with fragment only",
			data: "https://example.com/a#" + longWord,
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateString(5), prettyjson.WithEncOptsTruncateUrl(true)},
			want: `"https://example.com/a...62 chars,40F"`,
		},
		{
			name: "url statistics count runes of the original string",
			data: "https://example.com/a#" + strings.Repeat("中", 20),
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateString(5), prettyjson.WithEncOptsTruncateUrl(true)},
			want: `"https://example.com/a...42 chars,20F"`,
		},
		{
			name: "url without query or fragment is kept as is",
			data: "https://example.com/" + strings.Repeat("中", 20),
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateString(5), prettyjson.WithEncOptsTruncateUrl(true)},
			want: `"https://example.com/` + strings.Repeat("中", 20) + `"`,
		},
		{
			name: "colon prefixed text is not a url",
			data: "Error: " + longWord,
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateString(5), prettyjson.WithEncOptsTruncateUrl(true)},
			want: `"Error...47 chars"`,
		},
		{
			name: "colon prefixed text is not forced long",
			data: "Error: " + longWord,
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateString(5), prettyjson.WithEncOptsForceLongUrl(true)},
			want: `"Error...47 chars"`,
		},
		{
			name: "data uri is truncated as string",
			data: "data:image/png;base64," + longWord,
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateString(5), prettyjson.WithEncOptsTruncateUrl(true)},
			want: `"data:...62 chars"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := prettyjson.Marshal(tt.data, tt.opts...)
			if err != nil {
				t.Fatalf("Marshal(%v) error: %v", tt.data, err)
			}
			if string(got) != tt.want {
				t.Errorf("Marshal(%v) = `%s`, want `%s`", tt.data, got, tt.want)
			}
		})
	}
}
