// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// TypeInfo for Parsing.
// Also includes a Lexical Analysis and Syntactic Analysis.

package generic

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/searKing/golang/tools/pkg/ast"
)

type TemplateType struct {
	Type       string // The type's name.
	Import     string // import path of the type.
	IsPointer  bool   // whether the value's type is ptr
	TypePrefix string // The type's prefix, such as [][]
}

func (typ *TemplateType) String() string {
	var value bytes.Buffer
	if typ.IsPointer {
		value.WriteByte('*')
	}
	value.WriteString(typ.TypePrefix)
	value.WriteString(joinImport(typ.Import, typ.Type))
	return value.String()
}

// GenericType<TemplateType0,TemplateType1,...>
// Map<github.com/searKing/golang/tools/common/ast/generic.TypeInfo, int>
type TypeInfo struct {
	// These fields are reset for each type being generated.
	Name   string // Name of the sync.Map type.
	Import string // import path of the sync.Map type.

	TemplateTypes []TemplateType
}

func splitImport(value string) (_import, _type string) {
	// github.com/searKing/golang/tools/common/ast/generic.TypeInfo
	// =>
	// github.com/searKing/golang/tools/common/ast/generic generic.TypeInfo
	extPos := strings.LastIndexByte(value, '.')
	if extPos < 0 {
		extPos = len(value) - 1
		return "", value
	}
	// github.com/searKing/golang/tools/common/ast/generic
	pkg := value[:extPos]
	// TypeInfo
	name := value[extPos+1:]

	// github.com/searKing/golang/tools/common/ast
	// generic
	refPos := strings.LastIndexByte(pkg, '/')
	if refPos < 0 {
		return pkg, fmt.Sprintf("%s.%s", pkg, name)
	}
	return pkg, fmt.Sprintf("%s.%s", pkg[refPos+1:], name)
}

func joinImport(_import, _type string) string {
	if strings.TrimSpace(_import) == "" {
		return _type
	}
	typ := strings.Split(_type, ".")
	return strings.Join([]string{_import, typ[len(typ)-1]}, ".")
}

func walk(tokens []ast.Token, current int, tokenInfos []TypeInfo) []TypeInfo {
	if len(tokens) <= current {
		return tokenInfos
	}

	token := tokens[current]
	if token.Type == ast.TokenTypeParen && token.Value == "," {
		current++
		return walk(tokens, current, tokenInfos)
	}

	if token.Type == ast.TokenTypeName {
		import_, name_ := splitImport(token.Value)
		node := TypeInfo{
			Import: import_,
			Name:   name_,
		}
		current++
		if current >= len(tokens) {
			tokenInfos = append(tokenInfos, node)
			return tokenInfos
		}
		token = tokens[current]

		if token.Type == ast.TokenTypeParen && token.Value == "<" {
			current++
			if current >= len(tokens) {
				panic(fmt.Sprintf("missing token: %s after %s", ">", token.Value))
			}
			token = tokens[current]

			for {
				var isPointer bool
				if token.Type == ast.TokenTypeParen && token.Value == "*" {
					isPointer = true
					current++
					if current >= len(tokens) {
						panic(fmt.Sprintf("missing token: %s after %s", ">", token.Value))
					}
					token = tokens[current]
				}
				templateNode := TemplateType{
					IsPointer: isPointer,
				}

				// prefix
				var prefix bytes.Buffer
				for {
					if token.Type == ast.TokenTypeParen && token.Value == "*" {
						prefix.WriteString("*")
						current++
						if current >= len(tokens) {
							panic(fmt.Sprintf("missing token: after %s", token.Value))
						}
						token = tokens[current]
						continue
					}
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
					break
				}
				templateNode.TypePrefix = prefix.String()

				if token.Type == ast.TokenTypeName {
					// map[[]string]int
					var typ bytes.Buffer
					if token.Value == "map" {
						next, mapType := consumeMap(tokens, current)
						current = next
						typ.WriteString(mapType)
						if current >= len(tokens) {
							panic(fmt.Sprintf("missing token: %s after %s", "]", token.Value))
						}
						token = tokens[current]
					} else {
						current++
						typ.WriteString(token.Value)
						if current >= len(tokens) {
							panic(fmt.Sprintf("missing token: %s after %s", ">", token.Value))
						}
					}
					import_, type_ := splitImport(typ.String())
					templateNode.Import = import_
					templateNode.Type = type_
					token = tokens[current]
				}
				node.TemplateTypes = append(node.TemplateTypes, templateNode)

				if token.Type == ast.TokenTypeParen && token.Value == "," {
					current++
					if current >= len(tokens) {
						panic(fmt.Sprintf("missing token: %s after %s", ">", token.Value))
					}
					token = tokens[current]
					continue
				}
				break
			}
			if current >= len(tokens) {
				panic(fmt.Sprintf("missing token: %s after %s", ">", token.Value))
			}

			token = tokens[current]
			if token.Type == ast.TokenTypeParen && token.Value == ">" {
				current++
			} else {
				// 最后如果我们没有匹配上任何类型的 token，那么我们抛出一个错误。
				panic(fmt.Sprintf("unexpected token: %s, expect a %q", token.Value, '>'))
			}

		}
		tokenInfos = append(tokenInfos, node)
	} else {
		panic(fmt.Sprintf("unexpected token: %s, expect a TokenTypeName", token.Value))
	}
	return walk(tokens, current, tokenInfos)
}

func Parser(tokens []ast.Token) []TypeInfo {
	// type <key, value>
	return walk(tokens, 0, nil)
}

func New(input string) []TypeInfo {
	return Parser(ast.Tokenizer([]rune(input)))
}
