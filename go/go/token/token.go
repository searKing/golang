// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package token

// Token is the set of lexical tokens of the Go programming language.
type Type int

const (
	// Special tokens
	TypeILLEGAL Type = iota
	TypeEOF
	TypeCOMMENT
	TypeIgnored

	TypeLiteral
	TypeOperator
	TypeKeyword
)

type Token struct {
	Typ   Type
	Value string
}
