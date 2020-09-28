// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package generic

import (
	"bytes"
	"fmt"

	"github.com/searKing/golang/tools/common/ast"
)

// map[[]string]map[]int
func consumeMap(tokens []ast.Token, current int) (int, string) {
	if len(tokens) <= current {
		return current, ""
	}

	token := tokens[current]

	var prefix bytes.Buffer
	if token.Type == ast.TokenTypeName && token.Value == "map" {
		prefix.WriteString(token.Value)
		current++
		if current >= len(tokens) {
			panic(fmt.Sprintf("missing token: %s after %s", "[", token.Value))
		}
		token = tokens[current]

		// for key type
		var expect int // count of "]" expected to receive
		for {
			if token.Type == ast.TokenTypeParen && token.Value == "[" {
				prefix.WriteString("[")
				current++
				expect++
				if current >= len(tokens) {
					panic(fmt.Sprintf("missing token: %s after %s", "]", token.Value))
				}
				token = tokens[current]
			} else if token.Type == ast.TokenTypeParen && token.Value == "]" {
				prefix.WriteString("]")
				current++
				expect--
				if current >= len(tokens) {
					panic(fmt.Sprintf("missing token: after %s", token.Value))
				}
				token = tokens[current]
			} else if token.Type == ast.TokenTypeParen &&
				(token.Value == "," || token.Value == "<" || token.Value == ">") {
				panic(fmt.Sprintf("unexpected token: %s, expect a %q", token.Value, ']'))
			} else {
				prefix.WriteString(token.Value)
				current++
				if current >= len(tokens) {
					panic(fmt.Sprintf("missing token: after %s", token.Value))
				}
				token = tokens[current]
			}

			// matched [xx[yyy]x] found, key type finished
			if expect == 0 {
				break
			}
		}
		// for value type
		for {
			// for slice
			if token.Type == ast.TokenTypeParen && token.Value == "[" {
				current++
				if current >= len(tokens) {
					panic(fmt.Sprintf("missing token: %s after %s", "]", token.Value))
				}
				token = tokens[current]

				if token.Type == ast.TokenTypeParen && token.Value == "]" {
					prefix.WriteString("[]")
					current++
					if current >= len(tokens) {
						panic(fmt.Sprintf("missing token: after %s", token.Value))
					}
					token = tokens[current]
				} else {
					// 最后如果我们没有匹配上任何类型的 token，那么我们抛出一个错误。
					panic(fmt.Sprintf("unexpected token: %s, expect a %q", token.Value, ']'))
				}
				continue
			}
			// for map
			if token.Type == ast.TokenTypeName && token.Value == "map" {
				next, mapType := consumeMap(tokens, current)
				current = next
				prefix.WriteString(mapType)
				if current >= len(tokens) {
					break
				}
				token = tokens[current]
				continue
			}
			if token.Type == ast.TokenTypeParen &&
				(token.Value == "," || token.Value == "<" || token.Value == ">") {
				break
			}

			prefix.WriteString(token.Value)
			current++
			if current >= len(tokens) {
				panic(fmt.Sprintf("missing token: after %s", token.Value))
			}
			token = tokens[current]
		}
	}
	return current, prefix.String()
}
