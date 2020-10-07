// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// Arguments to format are:
//	SqlxJsonType: NullJson type name
//	ValueType: value type name
//	NilValue: nil value of map type
const tmplJson = `

// {{.SqlxJsonType}} represents an interface that may be null.
// {{.SqlxJsonType}} implements the Scanner interface so
// it can be used as a scan destination, similar to sql.NullString.
{{- if ne .SqlxJsonType .ValueType}}
{{- if .CanAlias}}
type {{.SqlxJsonType}} = {{.ValueType}}
{{- else}}
type {{.SqlxJsonType}} {{.ValueType}}
{{- end}}
{{- end}}


func (_ {{.SqlxJsonType}}) TableName() string {
	return "ops_mall"
}

func (_ {{.SqlxJsonType}}) Column(col {{.SqlxJsonType}}Field) string {
	return strings.SnakeCase(col.String())
}

func (m {{.SqlxJsonType}}) TableColumn(col {{.SqlxJsonType}}Field) string {
	return fmt.Sprintf("%s.%s", m.TableName(), m.Column(col))
}










// {{.SqlxJsonType}} represents an interface that may be null.
// {{.SqlxJsonType}} implements the Scanner interface so
// it can be used as a scan destination, similar to sql.NullString.
{{- if ne .SqlxJsonType .ValueType}}
{{- if .CanAlias}}
type {{.SqlxJsonType}} = {{.ValueType}}
{{- else}}
type {{.SqlxJsonType}} {{.ValueType}}
{{- end}}
{{- end}}

// Scan implements the sql.Scanner interface.
func (nj *{{.SqlxJsonType}}) Scan(src interface{}) error {
	if src == nil {
{{- if .CanAlias}}
		*nj = {{.NilValue}}
{{- else}}
		*nj = {{.SqlxJsonType}}({{.NilValue}})
{{- end}}
		return nil
	}

	var err error
	switch src := src.(type) {
	case string:
		err = json.Unmarshal([]byte(src), nj)
	case []byte:
		err = json.Unmarshal(src, nj)
	case time.Time:
		srcBytes, _ := json.Marshal(src)
		err = json.Unmarshal(srcBytes, nj)
	case nil:
{{- if .CanAlias}}
		*nj = {{.NilValue}}
{{- else}}
		*nj = {{.SqlxJsonType}}({{.NilValue}})
{{- end}}
		err = nil
	default:
		srcBytes, _ := json.Marshal(src)
		err = json.Unmarshal(srcBytes, nj)
	}
	if err == nil {
		return nil
	}

	return fmt.Errorf("unsupported Scan, storing driver.Value type %%T into type %%T : %%w", src, nj, err)
}

// Value implements the driver.Valuer interface.
func (nj {{.SqlxJsonType}}) Value() (driver.Value, error) {
	return json.Marshal(nj)
}
`
