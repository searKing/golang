// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// Arguments to format are:
//	SqlJsonType: NullJson type name
//	ValueType: value type name
//	NilValue: nil value of map type
const tmplJson = `

// {{.SqlJsonType}} represents an interface that may be null.
// {{.SqlJsonType}} implements the Scanner interface so
// it can be used as a scan destination, similar to sql.NullString.
{{- if ne .SqlJsonType .ValueType}}
{{- if .CanAlias}}
type {{.SqlJsonType}} = {{.ValueType}}
{{- else}}
type {{.SqlJsonType}} {{.ValueType}}
{{- end}}
{{- end}}

// Scan implements the sql.Scanner interface.
func (nj *{{.SqlJsonType}}) Scan(src interface{}) error {
	if src == nil {
{{- if .CanAlias}}
		*nj = {{.NilValue}}
{{- else}}
		*nj = {{.SqlJsonType}}({{.NilValue}})
{{- end}}
		return nil
	}

	var err error
	switch src := src.(type) {
	case string:
		if len(src) > 0 {
			err = json.Unmarshal([]byte(src), nj)
		}
	case []byte:
		if len(src) > 0 {
			err = json.Unmarshal(src, nj)
		}
	case time.Time:
		srcBytes, _ := json.Marshal(src)
		err = json.Unmarshal(srcBytes, nj)
	case nil:
{{- if .CanAlias}}
		*nj = {{.NilValue}}
{{- else}}
		*nj = {{.SqlJsonType}}({{.NilValue}})
{{- end}}
		err = nil
	default:
		srcBytes, _ := json.Marshal(src)
		err = json.Unmarshal(srcBytes, nj)
	}
	if err == nil {
		return nil
	}

	return fmt.Errorf("unsupported Scan, storing driver.Value type %T into type %T : %w", src, nj, err)
}

// Value implements the driver.Valuer interface.
func (nj {{.SqlJsonType}}) Value() (driver.Value, error) {
	return json.Marshal(nj)
}
`
