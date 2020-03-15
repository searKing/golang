// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package defaults

import "reflect"

// Convert v
// apply custom interface on v along with tag
func userDefinedConvertFunc(v reflect.Value, tag reflect.StructTag) (isUserDefined bool, err error) {
	isUserDefined = true
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return
	}
	m, ok := v.Interface().(Converter)
	if !ok {
		return
	}
	return isUserDefined, m.ConvertDefault(v, tag)
}

// Convert &v
// apply custom interface on &v along with tag
func addrUserDefinedConvertFunc(v reflect.Value, tag reflect.StructTag) (isUserDefined bool, err error) {
	isUserDefined = true
	vaddr := v.Addr()
	if vaddr.IsNil() {
		return
	}
	m := vaddr.Interface().(Converter)
	return isUserDefined, m.ConvertDefault(v, tag)
}

// newTypeConverter constructs an convertFunc for a type.
// The returned encoder only checks CanAddr when allowAddr is true.
func newTypeConverter(convFn convertFunc, t reflect.Type, allowAddr bool) convertFunc {
	// Handle UserDefined Case
	// t.ConvertDefault(value, tag)
	if t.Implements(converterType) {
		return userDefinedConvertFunc
	}

	// Handle UserDefined Case
	// Convert &v, iterate only once
	// (&t).ConvertDefault(value, tag) || convFn(value, tag)
	if t.Kind() != reflect.Ptr && allowAddr {
		if reflect.PtrTo(t).Implements(converterType) {
			return newCondAddrConvertFunc(addrUserDefinedConvertFunc, newTypeConverter(convFn, t, false))
		}
	}
	return convFn
}

// If CanAddr then get addr and handle else handle directly
type condAddrConvertFunc struct {
	canAddrConvert, elseConvert convertFunc
}

func (ce *condAddrConvertFunc) handle(v reflect.Value, tag reflect.StructTag) (isUserDefined bool, err error) {
	if v.CanAddr() {
		return ce.canAddrConvert(v, tag)
	}
	return ce.elseConvert(v, tag)
}

// newCondAddrConverter returns an encoder that checks whether its structTag
// CanAddr and delegates to canAddrConvert if so, else to elseConvert.
func newCondAddrConvertFunc(canAddrConvert, elseConvert convertFunc) convertFunc {
	convFn := &condAddrConvertFunc{canAddrConvert: canAddrConvert, elseConvert: elseConvert}
	return convFn.handle
}
