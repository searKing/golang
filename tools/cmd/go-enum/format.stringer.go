// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// Arguments to format are:
//	[1]: type name
//	[2]: size of index element (8 for uint8 etc.)
//	[3]: less than zero check (for signed types)
const stringOneRun = `
func _() {
	var _nil_%[1]s_value = func() (val %[1]s) { return }()

	// An "cannot convert %[1]s literal (type %[1]s) to type fmt.Stringer" compiler error signifies that the base type have changed.
	// Re-run the go-enum command to generate them again.
	var _ fmt.Stringer = _nil_%[1]s_value
}

func (i %[1]s) String() string {
	if %[3]si >= %[1]s(len(_%[1]s_index)-1) {
		return "%[1]s(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _%[1]s_name[_%[1]s_index[i]:_%[1]s_index[i+1]]
}
`

// Arguments to format are:
//	[1]: type name
//	[2]: lowest defined value for type, as a string
//	[3]: size of index element (8 for uint8 etc.)
//	[4]: less than zero check (for signed types)
/*
 */
const stringOneRunWithOffset = `
func _() {
	var _nil_%[1]s_value = func() (val %[1]s) { return }()

	// An "cannot convert %[1]s literal (type %[1]s) to type fmt.Stringer" compiler error signifies that the base type have changed.
	// Re-run the go-enum command to generate them again.
	var _ fmt.Stringer = _nil_%[1]s_value
}

func (i %[1]s) String() string {
	i -= %[2]s
	if %[4]si >= %[1]s(len(_%[1]s_index)-1) {
		return "%[1]s(" + strconv.FormatInt(int64(i + %[2]s), 10) + ")"
	}
	return _%[1]s_name[_%[1]s_index[i] : _%[1]s_index[i+1]]
}
`

// Argument to format is the type name.
const stringMap = `
func _() {
	var _nil_%[1]s_value = func() (val %[1]s) { return }()

	// An "cannot convert %[1]s literal (type %[1]s) to type fmt.Stringer" compiler error signifies that the base type have changed.
	// Re-run the go-enum command to generate them again.
	var _ fmt.Stringer = _nil_%[1]s_value
}

func (i %[1]s) String() string {
	if str, ok := _%[1]s_map[i]; ok {
		return str
	}
	return "%[1]s(" + strconv.FormatInt(int64(i), 10) + ")"
}
`
