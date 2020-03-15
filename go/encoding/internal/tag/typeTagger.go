// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tag

import (
	"reflect"
	"sync"
)

var tagFuncs tagFuncMap // map[reflect.Type]convertFunc

var taggerType = reflect.TypeOf(new(Tagger)).Elem()

func typeTagFunc(t reflect.Type) tagFunc {
	if fi, ok := tagFuncs.Load(t); ok {
		return fi
	}

	// To deal with recursive types, populate the map with an
	// indirect func before we build it. This type waits on the
	// real func (f) to be ready and then calls it. This indirect
	// func is only used for recursive types.
	var (
		wg sync.WaitGroup
		f  tagFunc
	)
	wg.Add(1)
	fi, loaded := tagFuncs.LoadOrStore(t, tagFunc(func(e *tagState, v reflect.Value, opts tagOpts) (isUserDefined bool) {
		// wait until f is assigned elsewhere
		wg.Wait()
		return f(e, v, opts)
	}))
	if loaded {
		return fi
	}

	// Compute the real encoder and replace the indirect func with it.
	f = newTypeTagger(t, true)
	wg.Done()
	tagFuncs.Store(t, f)
	return f
}

// newTypeTagger constructs an tagFunc for a type.
// The returned encoder only checks CanAddr when allowAddr is true.
func newTypeTagger(t reflect.Type, allowAddr bool) tagFunc {
	// Handle UserDefined Case
	// Convert v
	if t.Implements(taggerType) {
		return userDefinedTagFunc
	}

	// Handle UserDefined Case
	// Convert &v, iterate only once
	if t.Kind() != reflect.Ptr && allowAddr {
		if reflect.PtrTo(t).Implements(taggerType) {
			return newCondAddrTagFunc(addrUserDefinedTagFunc, newTypeTagger(t, false))
		}
	}

	// Handle BuiltinDefault Case
	switch t.Kind() {
	case reflect.Struct:
		return newStructTagFunc(t)
	case reflect.Ptr:
		return newPtrTagFunc(t)
	case reflect.Bool:
		fallthrough
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fallthrough
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		fallthrough
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		fallthrough
	case reflect.String:
		fallthrough
	case reflect.Interface:
		fallthrough
	case reflect.Map:
		fallthrough
	case reflect.Slice:
		fallthrough
	case reflect.Array:
		fallthrough
	default:
		return newNopConverter(t)
		//return unsupportedTypeConverter
	}
}
