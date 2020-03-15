// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package structinfo

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/searKing/golang/go/container/lru"
)

// --------------------------------------------------------------------------
// Maintain a mapping of keys to structure field indexes

type StructInfo struct {
	fieldsLRU lru.LRU
	Zero      reflect.Value
}
type fieldInfo struct {
	Key  string
	Num  int
	Tags map[string]string
}

var structMap = make(map[reflect.Type]*StructInfo)
var structMapMutex sync.RWMutex

type externalPanic string

func (e externalPanic) String() string {
	return string(e)
}
func IsStructFieldPrivate(field reflect.StructField) bool {
	if field.PkgPath != "" && !field.Anonymous {
		return true // Private field
	}
	return false
}

// Single Instance
// Ignore the private field
func GetStructInfo(st reflect.Type) (*StructInfo, error) {
	// lazyFound
	structMapMutex.RLock()
	sinfo, found := structMap[st]
	structMapMutex.RUnlock()
	if found {
		return sinfo, nil
	}

	// Traversal all Fields of the Struct
	n := st.NumField()
	fieldsLRU := lru.LRU{}
	inlineMap := -1
	for i := 0; i != n; i++ {
		field := st.Field(i)
		if IsStructFieldPrivate(field) {
			continue // Private field
		}

		info := fieldInfo{Num: i}

		tag := field.Tag.Get("bson")
		if tag == "" && strings.Index(string(field.Tag), ":") < 0 {
			tag = string(field.Tag)
		}
		if tag == "-" {
			continue
		}

		inline := false
		fields := strings.Split(tag, ",")
		if len(fields) > 1 {
			for _, flag := range fields[1:] {
				switch flag {
				case "omitempty":
					info.Tags["OmitEmpty"] = "omitempty"
				case "inline":
					inline = true
				default:
					msg := fmt.Sprintf("Unsupported flag %q in tag %q of type %s", flag, tag, st)
					panic(externalPanic(msg))
				}
			}
			tag = fields[0]
		}

		if inline {
			switch field.Type.Kind() {
			case reflect.Map:
				if inlineMap >= 0 {
					return nil, errors.New("Multiple ,inline maps in struct " + st.String())
				}
				if field.Type.Key() != reflect.TypeOf("") {
					return nil, errors.New("Option ,inline needs a map with string keys in struct " + st.String())
				}
				inlineMap = info.Num
			case reflect.Struct:
				sinfo, err := GetStructInfo(field.Type)
				if err != nil {
					return nil, err
				}

				for _, finfoPair := range sinfo.fieldsLRU.Pairs() {
					if _, found := fieldsLRU.Find(finfoPair.Key); found {
						msg := "Duplicated key '" + finfoPair.Key.(string) + "' in struct " + st.String()
						return nil, errors.New(msg)
					}
					fieldsLRU.AddPair(finfoPair)
				}
			default:
				panic("Option ,inline needs a struct value or map field")
			}
			continue
		}

		if tag != "" {
			info.Key = tag
		} else {
			info.Key = strings.ToLower(field.Name)
		}

		if _, found = fieldsLRU.Find(info.Key); found {
			msg := "Duplicated key '" + info.Key + "' in struct " + st.String()
			return nil, errors.New(msg)
		}
		fieldsLRU.Add(info.Key, info)
	}
	sinfo = &StructInfo{
		fieldsLRU,
		reflect.New(st).Elem(),
	}
	structMapMutex.Lock()
	structMap[st] = sinfo
	structMapMutex.Unlock()
	return sinfo, nil
}
