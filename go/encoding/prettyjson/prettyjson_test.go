// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package prettyjson_test

import (
	"testing"

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
			want: `"https://example.com/tests/1.html...72 chars,6Q9F]"`}, // take url as url
		{
			data: longUrl,
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateString(1), prettyjson.WithEncOptsTruncateStringIfMoreThan(4),
				prettyjson.WithEncOptsForceLongUrl(true)},
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
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateMap(1), prettyjson.WithEncOptsTruncateSliceOrArrayIfMoreThan(4)},
			want: `{"0":"red","1...5 pairs":"green"}`},
		{
			data: longMap,
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateMap(3), prettyjson.WithEncOptsTruncateSliceOrArrayIfMoreThan(1)},
			want: `{"0":"red","1":"green","2":"blue","3...5 pairs":"white"}`},
		{
			data: longMap,
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateMap(3), prettyjson.WithEncOptsTruncateSliceOrArrayIfMoreThan(1),
				prettyjson.WithEncOptsOmitStatistics(true)},
			want: `{"0":"red","1":"green","2":"blue"}`},
		{
			data: longMap,
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateMap(4), prettyjson.WithEncOptsTruncateSliceOrArrayIfMoreThan(1)},
			want: `{"0":"red","1":"green","2":"blue","3":"white","4":"black"}`},
		{
			data: longMap,
			opts: []prettyjson.EncOptsOption{prettyjson.WithEncOptsTruncateMap(4), prettyjson.WithEncOptsTruncateSliceOrArrayIfMoreThan(1),
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
			want: `{"ID":1,"Name":"Reds","LongName":"The quick ...43 chars","Colors":["The quick ...43 chars","Crimson","...5 elems"],"ColorById":{"0":"red","1":"green","2...5 pairs":"blue"},"ColorEnumById":{"0":0,"1":1,"2":2,"3":3},"Url":"https://example.com/tests/1?foo=1&bar=baz","LongUrl":"https://example.com/tests/1.html...72 chars,6Q9F]"}`},
	}
	for i, tt := range tests {
		got, err := prettyjson.Marshal(tt.data, tt.opts...)
		if err != nil {
			t.Errorf("#%d: Marshal(%v) error: %v", i, tt.data, err)
		}
		if tt.want != string(got) {
			t.Errorf("#%d: Marshal(%v) = `%v`, want `%s`", i, tt.data, string(got), tt.want)
		}
	}
}
