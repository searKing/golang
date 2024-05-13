// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tag

import (
	"reflect"
	"sync"

	reflect_ "github.com/searKing/golang/go/reflect"
)

// A field represents a single field found in a struct.
type field struct {
	name      string
	structTag reflect.StructTag

	index []int
	typ   reflect.Type
}

var fieldCache fieldMap

type fieldMap struct {
	fields sync.Map // map[reflect.Type][]field
}

func (m *fieldMap) Store(type_ reflect.Type, fields []field) {
	m.fields.Store(type_, fields)
}

func (m *fieldMap) LoadOrStore(type_ reflect.Type, fields []field) ([]field, bool) {
	actual, loaded := m.fields.LoadOrStore(type_, fields)
	if actual == nil {
		return nil, loaded
	}
	return actual.([]field), loaded
}

func (m *fieldMap) Load(type_ reflect.Type) ([]field, bool) {
	fields, ok := m.fields.Load(type_)
	if fields == nil {
		return nil, ok
	}
	return fields.([]field), ok
}

func (m *fieldMap) Delete(type_ reflect.Type) {
	m.fields.Delete(type_)
}

func (m *fieldMap) Range(f func(type_ reflect.Type, fields []field) bool) {
	m.fields.Range(func(type_, fields any) bool {
		return f(type_.(reflect.Type), fields.([]field))
	})
}

// cachedTypeFields is like typeFields but uses a cache to avoid repeated work.
func cachedTypeFields(t reflect.Type) []field {
	if f, ok := fieldCache.Load(t); ok {
		return f
	}
	var fields []field
	reflect_.WalkTypeDFS(t, reflect_.FieldTypeInfoHandlerFunc(
		func(info reflect_.FieldTypeInfo) (goon bool) {
			// ignore struct's root
			if info.Depth() == 0 {
				return true
			}

			sf, ok := info.StructField()
			if !ok {
				return true
			}

			fields = append(fields, field{
				name:      sf.Name,
				structTag: sf.Tag,
				index:     info.Index(),
				typ:       sf.Type,
			})
			return true
		}))
	f, _ := fieldCache.LoadOrStore(t, fields)
	return f
}
