// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

var yamlImportPackages []string

//var yamlImportPackages = []string{`gopkg.in/yaml.v2`}

// Arguments to format are:
//	[1]: type name
const yamlTemplate = `
//func _() {
//	var _nil_%[1]s_value = func() (val %[1]s) { return }()
//
//	// An "cannot convert %[1]s literal (type %[1]s) to type yaml.Marshaler" compiler error signifies that the base type have changed.
//	// Re-run the go-enum command to generate them again.
//	var _ yaml.Marshaler = _nil_%[1]s_value
//
//	// An "cannot convert %[1]s literal (type %[1]s) to type yaml.Unmarshaler" compiler error signifies that the base type have changed.
//	// Re-run the go-enum command to generate them again.
//	var _ yaml.Unmarshaler = &_nil_%[1]s_value
//}

// MarshalYAML implements a YAML Marshaler for %[1]s
func (i %[1]s) MarshalYAML() (interface{}, error) {
	return i.String(), nil
}

// UnmarshalYAML implements a YAML Unmarshaler for %[1]s
func (i *%[1]s) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}

	var err error
	*i, err = Parse%[1]sString(s)
	return err
}
`
