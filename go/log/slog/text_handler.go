// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slog

import (
	"encoding"
	"encoding/json"
	"fmt"
	"log/slog"
	"reflect"
	"strconv"
	"unicode"
	"unicode/utf8"

	"github.com/searKing/golang/go/unsafe"
)

func appendTextValue(s *handleState, v slog.Value) error {
	switch v.Kind() {
	case slog.KindString:
		s.appendStringMayQuote(v.String())
	case slog.KindTime:
		s.appendTime(v.Time())
	case slog.KindAny:
		a := v.Any()
		if tm, ok := a.(encoding.TextMarshaler); ok {
			data, err := tm.MarshalText()
			if err != nil {
				return err
			}
			s.appendBytesMayQuote(data)
			return nil
		}
		if err, ok := a.(error); ok && err != nil {
			s.appendStringMayQuote(err.Error())
			return nil
		}
		if tm, ok := a.(fmt.Stringer); ok {
			s.appendStringMayQuote(tm.String())
			return nil
		}
		if tm, ok := a.(json.Marshaler); ok {
			data, err := tm.MarshalJSON()
			if err != nil {
				return err
			}
			s.appendBytesMayQuote(data)
			return nil
		}
		if bs, ok := byteSlice(a); ok {
			// As of Go 1.19, this only allocates for strings longer than 32 bytes.
			s.buf.WriteString(strconv.Quote(unsafe.BytesToString(bs)))
			return nil
		}
		s.appendStringMayQuote(fmt.Sprintf("%+v", v.Any()))
	default:
		s.appendStringMayQuote(v.String())
	}
	return nil
}

// byteSlice returns its argument as a []byte if the argument's
// underlying type is []byte, along with a second return value of true.
// Otherwise it returns nil, false.
func byteSlice(a any) ([]byte, bool) {
	if bs, ok := a.([]byte); ok {
		return bs, true
	}
	// Like Printf's %s, we allow both the slice type and the byte element type to be named.
	t := reflect.TypeOf(a)
	if t != nil && t.Kind() == reflect.Slice && t.Elem().Kind() == reflect.Uint8 {
		return reflect.ValueOf(a).Bytes(), true
	}
	return nil, false
}

func needsQuoting(s string, unprintableOnly bool) bool {
	if len(s) == 0 {
		return true
	}
	for i := 0; i < len(s); {
		b := s[i]
		if b < utf8.RuneSelf {
			// Quote anything except a backslash that would need quoting in a
			// JSON string, as well as space and '='
			if unprintableOnly {
				if !safeSet[b] {
					return true
				}
			} else {
				if b != '\\' && (b == ' ' || b == '=' || !safeSet[b]) {
					return true
				}
			}
			i++
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		if unprintableOnly {
			if r == utf8.RuneError || !unicode.IsPrint(r) {
				return true
			}
		} else {
			if r == utf8.RuneError || unicode.IsSpace(r) || !unicode.IsPrint(r) {
				return true
			}
		}
		i += size
	}
	return false
}
