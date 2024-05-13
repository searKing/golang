// Code generated by "go-enum -type Nums"; DO NOT EDIT.

// Install go-enum by `go get install github.com/searKing/golang/tools/go-enum`
package main

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"fmt"
	"strconv"
)

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[One-1]
	_ = x[Two-2]
	_ = x[Three-3]
}

const _Nums_name = "OneTwoThree"

var _Nums_index = [...]uint8{0, 3, 6, 11}

func _() {
	var _nil_Nums_value = func() (val Nums) { return }()

	// An "cannot convert Nums literal (type Nums) to type fmt.Stringer" compiler error signifies that the base type have changed.
	// Re-run the go-enum command to generate them again.
	var _ fmt.Stringer = _nil_Nums_value
}

func (i Nums) String() string {
	i -= 1
	if i < 0 || i >= Nums(len(_Nums_index)-1) {
		return "Nums(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _Nums_name[_Nums_index[i]:_Nums_index[i+1]]
}

// New returns a pointer to a new addr filled with the Nums value passed in.
func (i Nums) New() *Nums {
	clone := i
	return &clone
}

var _Nums_values = []Nums{1, 2, 3}

var _Nums_name_to_values = map[string]Nums{
	_Nums_name[0:3]:  1,
	_Nums_name[3:6]:  2,
	_Nums_name[6:11]: 3,
}

// ParseNumsString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func ParseNumsString(s string) (Nums, error) {
	if val, ok := _Nums_name_to_values[s]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to Nums values", s)
}

// NumsValues returns all values of the enum
func NumsValues() []Nums {
	return _Nums_values
}

// IsANums returns "true" if the value is listed in the enum definition. "false" otherwise
func (i Nums) Registered() bool {
	for _, v := range _Nums_values {
		if i == v {
			return true
		}
	}
	return false
}

func _() {
	var _nil_Nums_value = func() (val Nums) { return }()

	// An "cannot convert Nums literal (type Nums) to type encoding.BinaryMarshaler" compiler error signifies that the base type have changed.
	// Re-run the go-enum command to generate them again.
	var _ encoding.BinaryMarshaler = &_nil_Nums_value

	// An "cannot convert Nums literal (type Nums) to type encoding.BinaryUnmarshaler" compiler error signifies that the base type have changed.
	// Re-run the go-enum command to generate them again.
	var _ encoding.BinaryUnmarshaler = &_nil_Nums_value
}

// MarshalBinary implements the encoding.BinaryMarshaler interface for Nums
func (i Nums) MarshalBinary() (data []byte, err error) {
	return []byte(i.String()), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface for Nums
func (i *Nums) UnmarshalBinary(data []byte) error {
	var err error
	*i, err = ParseNumsString(string(data))
	return err
}

func _() {
	var _nil_Nums_value = func() (val Nums) { return }()

	// An "cannot convert Nums literal (type Nums) to type json.Marshaler" compiler error signifies that the base type have changed.
	// Re-run the go-enum command to generate them again.
	var _ json.Marshaler = _nil_Nums_value

	// An "cannot convert Nums literal (type Nums) to type encoding.Unmarshaler" compiler error signifies that the base type have changed.
	// Re-run the go-enum command to generate them again.
	var _ json.Unmarshaler = &_nil_Nums_value
}

// MarshalJSON implements the json.Marshaler interface for Nums
func (i Nums) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for Nums
func (i *Nums) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("Nums should be a string, got %s", data)
	}

	var err error
	*i, err = ParseNumsString(s)
	return err
}

func _() {
	var _nil_Nums_value = func() (val Nums) { return }()

	// An "cannot convert Nums literal (type Nums) to type encoding.TextMarshaler" compiler error signifies that the base type have changed.
	// Re-run the go-enum command to generate them again.
	var _ encoding.TextMarshaler = _nil_Nums_value

	// An "cannot convert Nums literal (type Nums) to type encoding.TextUnmarshaler" compiler error signifies that the base type have changed.
	// Re-run the go-enum command to generate them again.
	var _ encoding.TextUnmarshaler = &_nil_Nums_value
}

// MarshalText implements the encoding.TextMarshaler interface for Nums
func (i Nums) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for Nums
func (i *Nums) UnmarshalText(text []byte) error {
	var err error
	*i, err = ParseNumsString(string(text))
	return err
}

//func _() {
//	var _nil_Nums_value = func() (val Nums) { return }()
//
//	// An "cannot convert Nums literal (type Nums) to type yaml.Marshaler" compiler error signifies that the base type have changed.
//	// Re-run the go-enum command to generate them again.
//	var _ yaml.Marshaler = _nil_Nums_value
//
//	// An "cannot convert Nums literal (type Nums) to type yaml.Unmarshaler" compiler error signifies that the base type have changed.
//	// Re-run the go-enum command to generate them again.
//	var _ yaml.Unmarshaler = &_nil_Nums_value
//}

// MarshalYAML implements a YAML Marshaler for Nums
func (i Nums) MarshalYAML() (any, error) {
	return i.String(), nil
}

// UnmarshalYAML implements a YAML Unmarshaler for Nums
func (i *Nums) UnmarshalYAML(unmarshal func(any) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}

	var err error
	*i, err = ParseNumsString(s)
	return err
}

func _() {
	var _nil_Nums_value = func() (val Nums) { return }()

	// An "cannot convert Nums literal (type Nums) to type driver.Valuer" compiler error signifies that the base type have changed.
	// Re-run the go-enum command to generate them again.
	var _ driver.Valuer = _nil_Nums_value

	// An "cannot convert Nums literal (type Nums) to type sql.Scanner" compiler error signifies that the base type have changed.
	// Re-run the go-enum command to generate them again.
	var _ sql.Scanner = &_nil_Nums_value
}

func (i Nums) Value() (driver.Value, error) {
	return i.String(), nil
}

func (i *Nums) Scan(value any) error {
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

	val, err := ParseNumsString(str)
	if err != nil {
		return err
	}

	*i = val
	return nil
}

// NumsSliceContains reports whether sunEnums is within enums.
func NumsSliceContains(enums []Nums, sunEnums ...Nums) bool {
	var seenEnums = map[Nums]bool{}
	for _, e := range sunEnums {
		seenEnums[e] = false
	}

	for _, v := range enums {
		if _, has := seenEnums[v]; has {
			seenEnums[v] = true
		}
	}

	for _, seen := range seenEnums {
		if !seen {
			return false
		}
	}

	return true
}

// NumsSliceContainsAny reports whether any sunEnum is within enums.
func NumsSliceContainsAny(enums []Nums, sunEnums ...Nums) bool {
	var seenEnums = map[Nums]struct{}{}
	for _, e := range sunEnums {
		seenEnums[e] = struct{}{}
	}

	for _, v := range enums {
		if _, has := seenEnums[v]; has {
			return true
		}
	}

	return false
}
