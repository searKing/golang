package main

import (
	"fmt"
	"strings"
)

// Value represents a declared constant.
type Value struct {
	originalName string // The name of the constant.
	name         string // The name with trimmed prefix.
	str          string // The string representation given by the "go/constant" package.

	mapImport string // import path of the sync.Map type.
	mapName   string // Name of the sync.Map type.

	keyImport     string // import path of the sync.Map's key.
	keyType       string // The type of the key in sync.Map.
	keyIsPointer  bool   // whether the value's type is ptr
	keyTypePrefix string // The type's prefix, such as []*[]

	valueImport     string // import path of the sync.Map's value.
	valueType       string // The type of the value in sync.Map.
	valueIsPointer  bool   // whether the value's type is ptr
	valueTypePrefix string // The type's prefix, such as []*[]
}

func (v *Value) String() string {
	return v.str
}

// Helpers

// createValAndNameDecl returns the pair of declarations for the run. The caller will add "var".
func (g *Generator) createValAndNameDecl(val Value) (string, string) {
	goRep := strings.NewReplacer(".", "_", "{", "_", "}", "_")

	nilValName := fmt.Sprintf("_nil_%s_%s_value", val.mapName, goRep.Replace(val.valueType))
	nilValDecl := fmt.Sprintf("%s = func() (val %s) { return }()", nilValName, val.valueType)

	return nilValName, nilValDecl
}
