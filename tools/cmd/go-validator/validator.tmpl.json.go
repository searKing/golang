// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// Arguments to format are:
//	StructType: struct type name to validate
const tmplValidator = `
{{- range .Structs}}
func (v *{{.StructType}}) Validate(validate *validator.Validate) error {
	if validate == nil {
		validate = validator.New()
	}
	return validate.Struct(v)
}

{{- end}}
`
