// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// Arguments to format are:
//	StructType: NullJson type trimmedStructName
//	TableName: value type trimmedStructName
//	NilValue: nil value of map type
const tmplJson = `

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	reflect_ "github.com/searKing/golang/go/reflect"
	"github.com/searKing/golang/go/strings"
)

// {{.StructType}} represents an orm of table {{.TableName}}.
// {{.StructType}} implements the Scanner interface so
// it can be used as a scan destination, similar to sql.NullString.

// TableName returns table's name
func (_ {{.StructType}}) TableName() string {
	return "{{.TableName}}"
}

// MarshalMap marshal themselves into or append a valid map
func (m {{.StructType}}) MarshalMap(valueByCol map[string]interface{}) map[string]interface{} {
	if valueByCol == nil {
		valueByCol = map[string]interface{}{}
	}
{{- range .Fields}}	
	valueByCol[m.MapColumn({{$.StructType}}Field{{.FieldType}})] = m.{{.FieldType}}
{{- end}}
	return valueByCol
}

// UnmarshalMap is the interface implemented by types
// that can unmarshal a map description of themselves.
func (m *{{.StructType}}) UnmarshalMap(valueByCol map[string]interface{}) error {
	for col, val := range valueByCol {
		switch col {
{{- range .Fields}}	
		case m.MapColumn({{$.StructType}}Field{{.FieldType}}):
			data, err := json.Marshal(val)
			if err != nil {
				return fmt.Errorf("marshal col %q, got %w", col, err)
			}
			err = json.Unmarshal(data, &m.{{.FieldType}})
			if err != nil {
				return fmt.Errorf("unmarshal col %q, got %w", col, err)
			}
{{- end}}
		}
	}
	return nil
}
// 列名
type {{.StructType}}Field int

const (
{{- range .Fields }}
	{{$.StructType}}Field{{.FieldType}}    {{$.StructType}}Field = iota
{{- end }}
)

func (f {{.StructType}}Field) String() string {
	switch f {
{{- range .Fields}}
	case {{$.StructType}}Field{{.FieldType}}:
		return "{{.DbName}}"
{{- end}}
	}
	return "{{.StructType}}Field(" + strconv.FormatInt(int64(f), 10) + ")"
}

func (a {{.StructType}}) ColumnEditor() *{{.StructType}}Columns {
	return &{{.StructType}}Columns{
		arg: a,
	}
}
func (a {{.StructType}}) Column(col {{.StructType}}Field) string {
	return strings.SnakeCase(col.String())
}

func (a {{.StructType}}) TableColumn(col {{.StructType}}Field) string {
	return fmt.Sprintf("%s.%s", a.TableName(), a.Column(col))
}

func (a {{.StructType}}) MapColumn(col {{.StructType}}Field) string {
	return fmt.Sprintf("%s_%s", a.TableName(), a.Column(col))
}

// columns

type {{.StructType}}Columns struct {
	arg  {{.StructType}}
	cols []string
}

func (c {{.StructType}}Columns) Columns(cols ...string) []string {
	return append(c.cols, cols...)
}

func (c *{{.StructType}}Columns) AppendColumn(col {{.StructType}}Field, forceAppend bool) *{{.StructType}}Columns {
	var zero = reflect_.IsZeroValue(reflect.ValueOf(c.arg).FieldByName(col.String()))

	if forceAppend || !zero {
		c.cols = append(c.cols, strings.SnakeCase(col.String()))
	}
	return c
}

func (c *{{.StructType}}Columns) AppendAll() *{{.StructType}}Columns {
	return c.
{{- range .Fields}}
		AppendColumn({{$.StructType}}Field{{.FieldType}}, false).
{{- end}}
		self()
}

func (c *{{.StructType}}Columns) self() *{{.StructType}}Columns {
	return c
}

`
