// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package prettyjson

import (
	"encoding"
	jsonv2 "encoding/json/v2"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"unicode/utf8"
)

//go:generate go-option -type "encOpts"
type encOpts struct {
	// escapeHTML causes '<', '>', and '&' to be escaped in JSON strings.
	// It defaults to false, unlike encoding/json.Marshal.
	escapeHTML bool

	truncateBytes                  int  // truncate bytes to this length
	truncateBytesIfMoreThan        int  // truncate bytes to this length if more than this length
	truncateString                 int  // truncate string to this length
	truncateStringIfMoreThan       int  // truncate string to this length if more than this length
	truncateMap                    int  // truncate map to this length
	truncateMapIfMoreThan          int  // truncate map to this length if more than this length
	truncateSliceOrArray           int  // truncate slice or array to this length
	truncateSliceOrArrayIfMoreThan int  // truncate slice or array to this length if more than this length
	truncateUrl                    bool // truncate query and fragment in url
	forceLongUrl                   bool // force long url
	// truncateMarshalerOutput truncates strings and arrays in the output of
	// slices, arrays and maps marshaled by their own methods, such as MarshalJSON.
	// Objects in the output keep all their members.
	truncateMarshalerOutput bool

	// omitEmpty omits struct fields with zero values as omitZero does, and
	// applies the encoding/json omitempty rules to fields tagged with omitempty.
	// Empty but non-nil slices and maps of untagged fields are kept.
	omitEmpty bool
	// omitZero omits struct fields with zero values, as if every field were
	// tagged with omitzero. A value is zero if its IsZero method reports true,
	// or else if it is the zero value of its type.
	omitZero       bool
	omitStatistics bool // omit statistics info
}

func truncateTo(limit, limitIfMoreThan int) (int, int) {
	if limit == 0 {
		limit = limitIfMoreThan
	}
	if limitIfMoreThan == 0 {
		limitIfMoreThan = limit
	}
	return limit, limitIfMoreThan
}

// truncateUrl drops the query and fragment of s, which parses as u, and
// appends the number of chars of s, query parameters and fragment chars.
func truncateUrl(s string, u *url.URL, opts encOpts) string {
	q := u.Query()
	f := u.Fragment
	if len(q) == 0 && len(f) == 0 {
		return s
	}
	u.RawQuery = ""
	u.Fragment = ""
	u.RawFragment = ""
	st := u.String()
	if !opts.omitStatistics {
		st += fmt.Sprintf("...%d chars,", utf8.RuneCountInString(s))
		if len(q) > 0 {
			st += fmt.Sprintf("%dQ", len(q))
		}
		if len(f) > 0 {
			st += fmt.Sprintf("%dF", utf8.RuneCountInString(f))
		}
	}
	if len(st) < len(s) {
		return st
	}
	return s
}

// resolveKeyName returns the object member name json/v2 encodes the map key k
// as, or an error if k must be left to json/v2.
func resolveKeyName(k reflect.Value) (string, error) {
	if t := k.Type(); hasMarshalMethod(t) {
		if k.Kind() == reflect.Pointer && k.IsNil() {
			return "", fmt.Errorf("nil map key of type %s", t)
		}
		m := k.Interface()
		if k.Kind() != reflect.Pointer {
			// json/v2 map keys are addressable, so pointer receivers count.
			p := reflect.New(t)
			p.Elem().Set(k)
			m = p.Interface()
		}
		// Same precedence as json/v2.
		switch m := m.(type) {
		case jsonv2.MarshalerTo, jsonv2.Marshaler:
			return "", fmt.Errorf("map key of type %s marshals itself as JSON", t)
		case encoding.TextAppender:
			buf, err := m.AppendText(nil)
			return string(buf), err
		case encoding.TextMarshaler:
			buf, err := m.MarshalText()
			return string(buf), err
		}
	}
	switch k.Kind() {
	case reflect.String:
		return k.String(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(k.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(k.Uint(), 10), nil
	}
	return "", fmt.Errorf("unsupported map key type %s", k.Type())
}
