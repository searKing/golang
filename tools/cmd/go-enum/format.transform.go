package main

import (
	"fmt"
	strings_ "github.com/searKing/golang/go/strings"
	"strings"
	"unicode"
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
	case "camel":
		mapfunc = func(s string) string {
			return strings_.CamelCase(strings.Map(preFunc, s), split)
		}
	case "small_camel":
		mapfunc = func(s string) string {
			return strings_.SmallCamelCase(strings.Map(preFunc, s), split)
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
		return
	}

	for i := range values {
		values[i].valueInfo.str = mapfunc(values[i].valueInfo.str)
	}
}
