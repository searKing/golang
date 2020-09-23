package main

// Arguments to format are:
//	SqlJsonType: NullJson type name
//	ValueType: value type name
//	NilValue: nil value of map type
const tmplJson = `

// {{.SqlJsonType}} represents an interface that may be null.
// {{.SqlJsonType}} implements the Scanner interface so
// it can be used as a scan destination, similar to sql.NullString.
{{if ne .SqlJsonType .ValueType}}
type {{.SqlJsonType}} = {{.ValueType}}
{{end}}

// Scan implements the sql.Scanner interface.
func (nj *{{.SqlJsonType}}) Scan(src interface{}) error {
	if src == nil {
		*nj = {{.NilValue}}
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
		*nj = {{.NilValue}}
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
func (nj {{.SqlJsonType}}) Value() (driver.Value, error) {
	return json.Marshal(nj)
}
`
