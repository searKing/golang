// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package token

import (
	"fmt"
)

var (
	Literals = map[string]struct{}{"<": {}, ">": {}}
)

type TokenizerFunc interface {
	IsLiteral() bool
	IsOperator() bool
	IsKeyword() bool
}

// type <value>, type <value>
func Tokenizer(inputs []rune, f func(inputs []rune, current int) (token Token, next int), strict bool) []Token {
	// type <key, value>
	current := 0
	var tokens []Token
	for current < len(inputs) {
		var token Token
		token, current = f(inputs, current)
		tokens = append(tokens, token)

		if token.Typ == TypeILLEGAL && strict {
			panic(fmt.Sprintf("I dont know what this token is: %s", token.Value))
		}

		if token.Typ == TypeEOF {
			return tokens
		}
	}
	return tokens
}
