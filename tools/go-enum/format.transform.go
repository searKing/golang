// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"strings"
	"unicode"

	strings_ "github.com/searKing/golang/go/strings"
)

func (g *Generator) transformValueNames(values []Value, transformMethod string) {
	var split = '_'
	var preFunc = func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return r
		}
		return split
	}

	var mapfunc func(s string) string
	//[nop, upper, lower, snake, camel, small_camel, kebab, dotted]
	switch strings.ToLower(strings.TrimSpace(transformMethod)) {
	case "", "nop":
		return
	case "upper":
		mapfunc = strings.ToUpper
	case "lower":
		mapfunc = strings.ToLower
	case "snake":
		mapfunc = func(s string) string {
			return strings_.SnakeCase(strings.Map(preFunc, s), split)
		}
	case "upper_camel":
		mapfunc = func(s string) string {
			return strings_.UpperCamelCase(strings.Map(preFunc, s), split)
		}
	case "lower_camel":
		mapfunc = func(s string) string {
			return strings_.LowerCamelCase(strings.Map(preFunc, s), split)
		}
	case "kebab":
		mapfunc = func(s string) string {
			return strings_.KebabCase(strings.Map(preFunc, s), split)
		}
	case "dotted":
		mapfunc = func(s string) string {
			return strings_.DotCase(strings.Map(preFunc, s), split)
		}
	default:
		panic(fmt.Sprintf("unknown transform method %s, only [nop, upper, lower, snake, camel, small_camel, kebab, dotted] is supported.", transformMethod))
	}

	for i := range values {
		values[i].nameInfo.trimmedName = mapfunc(values[i].nameInfo.trimmedName)
	}
}
