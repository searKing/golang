// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reflect

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	bytes_ "github.com/searKing/golang/go/bytes"
	strings_ "github.com/searKing/golang/go/strings"
)

// Truncate reset bytes and string, useful for dump into log if some field is huge
func Truncate(v interface{}, n int) {
	WalkValueBFS(reflect.ValueOf(v), FieldValueInfoHandlerFunc(func(info FieldValueInfo) (goon bool) {
		if !info.Value().CanSet() || !info.Value().CanInterface() || !info.Value().CanSet() {
			return true
		}
		vv := info.Value().Interface()

		switch vv := vv.(type) {
		case []byte:
			if len(vv) <= n {
				break
			}
			var buf bytes.Buffer
			buf.WriteString(fmt.Sprintf("size: %d, bytes: ", len(vv)))
			buf.Write(bytes_.Truncate(vv, n))
			info.Value().SetBytes(buf.Bytes())
			return true
		case string:
			if len(vv) <= n {
				break
			}
			var buf strings.Builder
			buf.WriteString(fmt.Sprintf("size: %d, string: ", len(vv)))
			buf.WriteString(strings_.Truncate(vv, n))
			info.Value().SetString(buf.String())
			return true
		}
		return true
	}))

}
