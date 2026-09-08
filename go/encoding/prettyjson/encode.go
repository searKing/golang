// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package prettyjson

import (
	"encoding"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	_ "unsafe"
)

//go:generate go-option -type "encOpts"
type encOpts struct {
	// quoted causes primitive fields to be encoded inside JSON strings.
	quoted bool `option:"-"`

	// escapeHTML causes '<', '>', and '&' to be escaped in JSON strings.
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

	omitEmpty      bool // omit empty value
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

func truncateUrl(u *url.URL, opts encOpts) (st string) {
	if u == nil {
		return ""
	}
	s := u.String()
	q := u.Query()
	f := u.Fragment
	u.Fragment = ""
	if len(q) > 0 || len(f) > 0 {
		u.RawQuery = ""
		u.Fragment = ""
		u.RawFragment = ""
		st := u.String()
		if !opts.omitStatistics {
			st += fmt.Sprintf("...%d chars", len(s))
			if len(q) > 0 {
				st += fmt.Sprintf(",%dQ", len(q))
			}
			if len(f) > 0 {
				st += fmt.Sprintf("%dF", len(f))
			}
			st += "]"
		}
		if len(st) < len(s) {
			s = st
		}
	}
	return s
}

func resolveKeyName(k reflect.Value) (string, error) {
	if k.Kind() == reflect.String {
		return k.String(), nil
	}
	if tm, ok := k.Interface().(encoding.TextMarshaler); ok {
		if k.Kind() == reflect.Pointer && k.IsNil() {
			return "", nil
		}
		buf, err := tm.MarshalText()
		return string(buf), err
	}
	switch k.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(k.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(k.Uint(), 10), nil
	}
	panic("unexpected map key type")
}
