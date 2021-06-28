// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

var sqlImportPackages = []string{`database/sql`, `database/sql/driver`}

// Arguments to format are:
//	[1]: type name
const sqpTemplate = `
func _() {
	var _nil_%[1]s_value = func() (val %[1]s) { return }()

	// An "cannot convert %[1]s literal (type %[1]s) to type driver.Valuer" compiler error signifies that the base type have changed.
	// Re-run the go-enum command to generate them again.
	var _ driver.Valuer = _nil_%[1]s_value


	// An "cannot convert %[1]s literal (type %[1]s) to type sql.Scanner" compiler error signifies that the base type have changed.
	// Re-run the go-enum command to generate them again.
	var _ sql.Scanner = &_nil_%[1]s_value
}

func (i %[1]s) Value() (driver.Value, error) {
	return i.String(), nil
}

func (i *%[1]s) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	str, ok := value.(string)
	if !ok {
		bytes, ok := value.([]byte)
		if !ok {
			return fmt.Errorf("value is not a byte slice")
		}

		str = string(bytes[:])
	}

	val, err := Parse%[1]sString(str)
	if err != nil {
		return err
	}
	
	*i = val
	return nil
}
`
