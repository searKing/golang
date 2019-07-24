package reflect

import "reflect"

func IsNilObject(obtained interface{}) (result bool) {
	if obtained == nil {
		result = true
	} else {
		return IsNilValue(reflect.ValueOf(obtained))
	}
	return
}
