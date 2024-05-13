// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reflect

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	bytes_ "github.com/searKing/golang/go/bytes"
	strings_ "github.com/searKing/golang/go/strings"
)

// TruncateString reset string, useful for dump into log if some field is huge
// v is truncated in place
// return interface{} same as truncated v for stream-like api
func TruncateString(v any, n int) any {
	return Truncate(v, func(v any) bool {
		_, ok := v.(string)
		return ok
	}, n)
}

// TruncateBytes reset bytes, useful for dump into log if some field is huge
// v is truncated in place
// return interface{} same as truncated v for stream-like api
func TruncateBytes(v any, n int) any {
	return Truncate(v, func(v any) bool {
		_, ok := v.([]byte)
		return ok
	}, n)
}

// Truncate reset bytes and string at each run of value c satisfying f(c), useful for dump into log if some field is huge
// v is truncated in place
// return interface{} same as truncated v for stream-like api
func Truncate(v any, f func(v any) bool, n int) any {
	truncate(reflect.ValueOf(v), f, n)
	return v
}

func truncate(v reflect.Value, f func(v any) bool, n int) {
	if !v.IsValid() {
		return
	}
	if IsNilType(v.Type()) {
		return
	}

	if v.CanInterface() {
		vv := v.Interface()
		if f(vv) {
			// handle v in place, stop visit sons
			if v.CanSet() {
				switch vv := vv.(type) {
				case []byte:
					if len(vv) <= n {
						break
					}
					var buf bytes.Buffer
					buf.WriteString(fmt.Sprintf("size: %d, bytes: ", len(vv)))
					buf.Write(bytes_.Truncate(vv, n))
					v.SetBytes(buf.Bytes())
					return
				case string:
					if len(vv) <= n {
						break
					}
					var buf strings.Builder
					buf.WriteString(fmt.Sprintf("size: %d, string: ", len(vv)))
					buf.WriteString(strings_.Truncate(vv, n))
					v.SetString(buf.String())
					return
				}
			}
			return
		}
	}

	// handle v's sons
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			truncate(v.Index(i), f, n)
		}
	case reflect.Struct:
		// Scan typ for fields to include.
		for i := 0; i < v.NumField(); i++ {
			truncate(v.Field(i), f, n)
		}
	case reflect.Map:
		iter := v.MapRange()
		for iter.Next() {
			//truncate(iter.Key(), f, n) // Key of Map is not addressable
			truncate(iter.Value(), f, n)
		}
	case reflect.Ptr:
		truncate(reflect.Indirect(v), f, n)
	}
	return
}
