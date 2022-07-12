// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

var textImportPackages = []string{`encoding`}

// Arguments to format are:
//	[1]: type name
const textTemplate = `
func _() {
	var _nil_%[1]s_value = func() (val %[1]s) { return }()

	// An "cannot convert %[1]s literal (type %[1]s) to type encoding.TextMarshaler" compiler error signifies that the base type have changed.
	// Re-run the go-enum command to generate them again.
	var _ encoding.TextMarshaler = _nil_%[1]s_value

	// An "cannot convert %[1]s literal (type %[1]s) to type encoding.TextUnmarshaler" compiler error signifies that the base type have changed.
	// Re-run the go-enum command to generate them again.
	var _ encoding.TextUnmarshaler = &_nil_%[1]s_value
}

// MarshalText implements the encoding.TextMarshaler interface for %[1]s
func (i %[1]s) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for %[1]s
func (i *%[1]s) UnmarshalText(text []byte) error {
	var err error
	*i, err = Parse%[1]sString(string(text))
	return err
}
`
