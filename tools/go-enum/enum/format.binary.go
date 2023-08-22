// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package enum

var binaryImportPackages = []string{`encoding`}

// Arguments to format are:
//
//	[1]: type name
const binaryTemplate = `
func _() {
	var _nil_%[1]s_value = func() (val %[1]s) { return }()

	// An "cannot convert %[1]s literal (type %[1]s) to type encoding.BinaryMarshaler" compiler error signifies that the base type have changed.
	// Re-run the go-enum command to generate them again.
	var _ encoding.BinaryMarshaler = &_nil_%[1]s_value


	// An "cannot convert %[1]s literal (type %[1]s) to type encoding.BinaryUnmarshaler" compiler error signifies that the base type have changed.
	// Re-run the go-enum command to generate them again.
	var _ encoding.BinaryUnmarshaler = &_nil_%[1]s_value
}

// MarshalBinary implements the encoding.BinaryMarshaler interface for %[1]s
func (i %[1]s) MarshalBinary() (data []byte, err error) {
	return []byte(i.String()), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface for %[1]s
func (i *%[1]s) UnmarshalBinary(data []byte) error {
	var err error
	*i, err = Parse%[1]sString(string(data))
	return err
}
`
