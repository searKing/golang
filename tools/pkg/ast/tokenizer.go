// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ast

import (
	"bytes"
	"fmt"
	"unicode"
)

type TokenType int

const (
	TokenTypeParen TokenType = iota
	TokenTypeName
)

type Token struct {
	Type  TokenType
	Value string
}

// type <T,U>, type <interface{},*U>
// NumValue<int, string>, AnotherNumValue<int, interface{}>
// =>
// ["NumValue","<", "int", ",", "string", ">", ",", "AnotherNumValue", "<", "int", ",", "interface{}", ">"]
func Tokenizer(inputs []rune) []Token {
	// type <T, U>
	current := 0
	var tokens []Token
	for current < len(inputs) {
		char := inputs[current]
		if char == '<' {
			tokens = append(tokens, Token{
				Type:  TokenTypeParen,
				Value: "<",
			})
			current++
			continue
		}

		if char == '>' {
			tokens = append(tokens, Token{
				Type:  TokenTypeParen,
				Value: ">",
			})
			current++
			continue
		}

		// Special case: * [ ] Type
		if char == '*' || char == '[' || char == ']' {
			tokens = append(tokens, Token{
				Type:  TokenTypeParen,
				Value: string(char),
			})
			current++
			continue
		}

		if unicode.IsSpace(char) {
			current++
			if current >= len(inputs) {
				break
			}
			char = inputs[current]
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
			// '.' for import reference and website of import path
			// '/', '-' for import path
			for unicode.IsLetter(char) || char == '_' || unicode.IsNumber(char) || char == '.' || char == '/' || char == '-' {
				value.WriteRune(char)
				current++
				if current >= len(inputs) {
					break
				}
				char = inputs[current]
			}

			// Special case: interface{}|struct{}
			if value.String() == "interface" || value.String() == "struct" {
				for {
					if unicode.IsSpace(char) {
						current++
						continue
					}
					break
				}
				// expect {}
				if char == '{' {
					current++
					if current >= len(inputs) {
						break
					}
					char = inputs[current]

					for {
						if unicode.IsSpace(char) {
							current++
							continue
						}
						break
					}

					if char == '}' {
						current++
						if current >= len(inputs) {
							break
						}
						char = inputs[current]
					} else {
						panic(fmt.Sprintf("I dont know what this character at %d is: %q", current, string(char)))
					}
					value.WriteString("{}")
				} else {
					panic(fmt.Sprintf("I dont know what this character at %d is: %q", current, string(char)))
				}
			}

			tokens = append(tokens, Token{
				Type:  TokenTypeName,
				Value: value.String(),
			})
			continue
		}

		if char == ',' {
			current++
			tokens = append(tokens, Token{
				Type:  TokenTypeParen,
				Value: ",",
			})
			continue
		}

		// 最后如果我们没有匹配上任何类型的 token，那么我们抛出一个错误。
		panic(fmt.Sprintf("I dont know what this character is: %s", string(char)))
	}

	return tokens
}
