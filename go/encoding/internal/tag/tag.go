// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tag

import "reflect"

type tagOpts struct {
	TagHandler func(val reflect.Value, tag reflect.StructTag) error
}

// Convert wrapper of convertState
func Tag(v interface{}, tagHandler func(val reflect.Value, tag reflect.StructTag) error) error {
	e := newTagState()
	err := e.handle(v, tagOpts{tagHandler})
	if err != nil {
		return err
	}

	e.Reset()
	tagStatePool.Put(e)
	return nil
}

// Tagger is the interface implemented by types that
// can marshal themselves into valid JSON.
type Tagger interface {
	TagDefault() error
}
