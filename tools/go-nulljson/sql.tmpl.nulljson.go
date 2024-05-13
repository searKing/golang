package main

// Arguments to format are:
//
//	SqlJsonType: NullJson type name
//	ValueType: value type name
//	NilValue: nil value of map type
const tmplNullJson = `

// {{.SqlJsonType}} represents an interface that may be null.
// {{.SqlJsonType}} implements the Scanner interface so
// it can be used as a scan destination, similar to sql.NullString.
type {{.SqlJsonType}} struct {
	Data {{.ValueType}}

	Valid bool // Valid is true if Data is not NULL
}

// Scan implements the sql.Scanner interface.
func (nj *{{.SqlJsonType}}) Scan(src interface{}) error {
	if src == nil {
		nj.Data, nj.Valid = {{.NilValue}}, false
		return nil
	}
	nj.Valid = true

	var err error
	switch src := src.(type) {
	case string:
		if len(src) > 0 {
			var v interface{} = &nj.Data
			switch v := v.(type) {
{{- if .ProtoJson }}
			case proto.Message:
				err = protojson.UnmarshalOptions{
					AllowPartial:   true,
					DiscardUnknown: true,
				}.Unmarshal([]byte(src), v)
{{- end}}
			default:
				err = json.Unmarshal([]byte(src), v)
			}
		}
	case []byte:
		if len(src) > 0 {
			var v interface{} = &nj.Data
			switch v := v.(type) {
{{- if .ProtoJson }}
			case proto.Message:
				err = protojson.UnmarshalOptions{
					AllowPartial:   true,
					DiscardUnknown: true,
				}.Unmarshal(src, v)
{{- end}}
			default:
				err = json.Unmarshal(src, v)
			}
		}
	case time.Time:
		srcBytes, _ := json.Marshal(src)
		var v interface{} = &nj.Data
		switch v := v.(type) {
		case proto.Message:
{{- if .ProtoJson }}
			err = protojson.UnmarshalOptions{
				AllowPartial:   true,
				DiscardUnknown: true,
			}.Unmarshal(srcBytes, v)
{{- end}}
		default:
			err = json.Unmarshal(srcBytes, v)
		}
	case nil:
		nj.Data = {{.NilValue}}
		err = nil
	default:
		srcBytes, _ := json.Marshal(src)
		var v interface{} = &nj.Data
		switch v := v.(type) {
{{- if .ProtoJson }}
		case proto.Message:
			err = protojson.UnmarshalOptions{
				AllowPartial:   true,
				DiscardUnknown: true,
			}.Unmarshal(srcBytes, v)
{{- end}}
		default:
			err = json.Unmarshal(srcBytes, v)
		}
	}
	if err == nil {
		return nil
	}

	return fmt.Errorf("unsupported Scan, storing driver.Value type %T into type %T: %w", src, nj.Data, err)
}

// Value implements the driver.Valuer interface.
func (nj {{.SqlJsonType}}) Value() (driver.Value, error) {
	if !nj.Valid {
		return nil, nil
	}
	var v interface{} = &nj.Data
	switch v := v.(type) {
{{- if .ProtoJson }}
	case proto.Message:
		return protojson.MarshalOptions{
			AllowPartial:      true,
		}.Marshal(v)
{{- end}}
	default:
		return json.Marshal(v)
	}
}
`
