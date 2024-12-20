// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package prettyjson_test

import (
	"fmt"
	"os"

	"github.com/searKing/golang/go/encoding/prettyjson"
)

func ExampleMarshal() {
	type ColorGroup struct {
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
	}
	group := ColorGroup{
		ID:            1,
		Name:          "Reds",
		LongName:      "The quick brown fox jumps over the lazy dog",
		Colors:        []string{"The quick brown fox jumps over the lazy dog", "Crimson", "Red", "Ruby", "Maroon"},
		ColorById:     map[string]string{"0": "red", "1": "green", "2": "blue", "3": "white"},
		ColorEnumById: map[string]int{"0": 0, "1": 1, "2": 2, "3": 3},
		Url:           "https://example.com/tests/1?foo=1&bar=baz",
		LongUrl:       "https://example.com/tests/1.html?foo=1&bar=baz&a=0&b=1&c=2&d=3#paragraph",
	}
	{
		b, err := prettyjson.Marshal(group,
			prettyjson.WithEncOptsTruncateString(10),
			prettyjson.WithEncOptsTruncateBytes(10),
			prettyjson.WithEncOptsTruncateSliceOrArray(2),
			prettyjson.WithEncOptsTruncateMap(2),
			prettyjson.WithEncOptsTruncateUrl(true),
			prettyjson.WithEncOptsEscapeHTML(false),
			prettyjson.WithEncOptsOmitEmpty(true))
		if err != nil {
			fmt.Println("error:", err)
		}
		_, _ = os.Stdout.Write(b)
		_, _ = os.Stdout.Write([]byte("\n"))
	}
	{
		b, err := prettyjson.Marshal(group,
			prettyjson.WithEncOptsTruncateString(20),
			prettyjson.WithEncOptsTruncateStringIfMoreThan(10),
			prettyjson.WithEncOptsTruncateBytes(20),
			prettyjson.WithEncOptsTruncateBytesIfMoreThan(10),
			prettyjson.WithEncOptsTruncateSliceOrArray(4),
			prettyjson.WithEncOptsTruncateSliceOrArrayIfMoreThan(2),
			prettyjson.WithEncOptsTruncateMap(4),
			prettyjson.WithEncOptsTruncateMapIfMoreThan(2),
			prettyjson.WithEncOptsTruncateUrl(true),
			prettyjson.WithEncOptsEscapeHTML(false),
			prettyjson.WithEncOptsForceLongUrl(false),
			prettyjson.WithEncOptsOmitEmpty(true))
		if err != nil {
			fmt.Println("error:", err)
		}
		_, _ = os.Stdout.Write(b)
		_, _ = os.Stdout.Write([]byte("\n"))
	}

	// Output:
	// {"ID":1,"Name":"Reds","LongName":"The quick ...43 chars","Colors":["The quick ...43 chars","Crimson","...5 elems"],"ColorById":{"0":"red","1":"green","2...4 pairs":"blue"},"ColorEnumById":{"0":0,"1":1,"2...4 pairs":2},"Url":"https://example.com/tests/1?foo=1&bar=baz","LongUrl":"https://example.com/tests/1.html...72 chars,6Q9F]"}
	// {"ID":1,"Name":"Reds","LongName":"The quick brown fox ...43 chars","Colors":["The quick brown fox ...43 chars","Crimson","Red","Ruby","...5 elems"],"ColorById":{"0":"red","1":"green","2":"blue","3":"white"},"ColorEnumById":{"0":0,"1":1,"2":2,"3":3},"Url":"https://example.com/tests/1?foo=1&bar=baz","LongUrl":"https://example.com/tests/1.html...72 chars,6Q9F]"}
}
