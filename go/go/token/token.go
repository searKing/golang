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
