// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tag

import (
	"reflect"
	"sync"
)

type tagFunc func(e *tagState, v reflect.Value, opts tagOpts) (isUserDefined bool)

// map[reflect.Type]tagFunc
type tagFuncMap struct {
	tagFns sync.Map
}

func (t *tagFuncMap) Store(reflectType reflect.Type, fn tagFunc) {
	t.tagFns.Store(reflectType, fn)
}

func (t *tagFuncMap) LoadOrStore(reflectType reflect.Type, fn tagFunc) (tagFunc, bool) {
	actual, loaded := t.tagFns.LoadOrStore(reflectType, fn)
	if actual == nil {
		return nil, loaded
	}
	return actual.(tagFunc), loaded
}

func (t *tagFuncMap) Load(reflectType reflect.Type) (tagFunc, bool) {
	fn, ok := t.tagFns.Load(reflectType)
	if fn == nil {
		return nil, ok
	}
	return fn.(tagFunc), ok
}

func (t *tagFuncMap) Delete(reflectType reflect.Type) {
	t.tagFns.Delete(reflectType)
}

func (t *tagFuncMap) Range(f func(reflectType reflect.Type, fn tagFunc) bool) {
	t.tagFns.Range(func(reflectType, fn interface{}) bool {
		return f(reflectType.(reflect.Type), fn.(tagFunc))
	})
}
