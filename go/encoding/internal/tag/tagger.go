package tag

import "reflect"

func valueTaggeFunc(v reflect.Value) tagFunc {
	if !v.IsValid() {
		return invalidValueTagFunc
	}
	return typeTagFunc(v.Type())
}
