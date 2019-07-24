// Copyright 2019 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// TypeInfo for Parsing.
// Also includes a Lexical Analysis and Syntactic Analysis.

package main

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

type tokenType int

const (
	tokenTypeParen tokenType = iota
	tokenTypeName
)

type _token struct {
	typ   tokenType
	value string
}

type typeInfo struct {
	// These fields are reset for each type being generated.
	eleName     string // Name of the atomic.Value type.
	eleImport   string // import path of the atomic.Value type.
	valueType   string // The type of the value in atomic.Value.
	valueImport string // import path of the atomic.Value's value.
}

// type <value>, type <value>
func tokenizer(inputs []rune) []_token {
	// type <key, value>
	current := 0
	var tokens []_token
	for current < len(inputs) {
		char := inputs[current]
		if char == '<' {
			tokens = append(tokens, _token{
				typ:   tokenTypeParen,
				value: "<",
			})
			current++
			continue
		}
		if char == '>' {
			tokens = append(tokens, _token{
				typ:   tokenTypeParen,
				value: ">",
			})
			current++
			continue
		}

		if unicode.IsSpace(char) {
			current++
			continue
		}

		var value bytes.Buffer

		// identifier = letter { letter | unicode_digit } .
		// letter        = unicode_letter | "_" .
		// decimal_digit = "0" … "9" .
		// octal_digit   = "0" … "7" .
		// hex_digit     = "0" … "9" | "A" … "F" | "a" … "f" .
		// newline        = /* the Unicode code point U+000A */ .
		// unicode_char   = /* an arbitrary Unicode code point except newline */ .
		// unicode_letter = /* a Unicode code point classified as "Letter" */ .
		// unicode_digit  = /* a Unicode code point classified as "Number, decimal digit" */ .
		if unicode.IsLetter(char) || char == '_' {
			for unicode.IsLetter(char) || char == '_' || unicode.IsNumber(char) || char == '.' {
				value.WriteRune(char)
				current++
				if current >= len(inputs) {
					break
				}
				char = inputs[current]
			}

			tokens = append(tokens, _token{
				typ:   tokenTypeName,
				value: value.String(),
			})
			continue
		}

		if char == ',' {
			current++
			tokens = append(tokens, _token{
				typ:   tokenTypeParen,
				value: ",",
			})
			continue
		}

		// 最后如果我们没有匹配上任何类型的 token，那么我们抛出一个错误。
		panic(fmt.Sprintf("I dont know what this character is: %s", string(char)))
	}

	return tokens
}

func splitImport(value string) (_import, _type string) {
	// a.b.c
	// a.b c
	extPos := strings.LastIndexByte(value, '.')
	if extPos < 0 {
		extPos = len(value) - 1
		return "", value
	}
	pkg := value[:extPos]
	name := value[extPos+1:]

	namPos := strings.LastIndexByte(pkg, '.')
	if namPos < 0 {
		return pkg, fmt.Sprintf("%s.%s", pkg, name)
	}
	return pkg, fmt.Sprintf("%s.%s", pkg[namPos+1:], name)
}

func walk(tokens []_token, current int, tokenInfos []typeInfo) []typeInfo {
	if len(tokens) <= current {
		return tokenInfos
	}

	token := tokens[current]
	if token.typ == tokenTypeParen && token.value == "," {
		current++
		return walk(tokens, current, tokenInfos)
	}

	if token.typ == tokenTypeName {
		mapImport, mapName := splitImport(token.value)
		node := typeInfo{
			eleImport: mapImport,
			eleName:   mapName,
		}
		current++
		if current >= len(tokens) {
			tokenInfos = append(tokenInfos, node)
			return tokenInfos
		}
		token = tokens[current]

		if token.typ == tokenTypeParen && token.value == "<" {
			current++
			if current >= len(tokens) {
				panic(fmt.Sprintf("missing token: %s after %s", ">", token.value))
			}
			token = tokens[current]

			if token.typ == tokenTypeName {
				keyImport, keyType := splitImport(token.value)
				node.valueImport = keyImport
				node.valueType = keyType
				current++
			}

			if current >= len(tokens) {
				panic(fmt.Sprintf("missing token: %s after %s", ">", token.value))
			}
			token = tokens[current]
			if token.typ == tokenTypeParen && token.value == ">" {
				current++
			} else {
				// 最后如果我们没有匹配上任何类型的 token，那么我们抛出一个错误。
				panic(fmt.Sprintf("unexpected token: %s", token.value))
			}

		}
		tokenInfos = append(tokenInfos, node)
	}
	return walk(tokens, current, tokenInfos)
}

func parser(tokens []_token) []typeInfo {
	// type <key, value>
	return walk(tokens, 0, nil)
}

func newTypeInfo(input string) []typeInfo {
	return parser(tokenizer([]rune(input)))
}
