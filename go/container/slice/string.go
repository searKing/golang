package slice

import "reflect"

func normalizeSlice(s []interface{}, as interface{}) interface{} {
	kind := reflect.ValueOf(as).Kind()
	switch kind {
	case reflect.Map:
		return normalizeSliceAsMap(s)
	}
	return s
}
func normalizeElem(elem, as interface{}) interface{} {
	return elem
}

func normalizeSliceAsMap(s []interface{}) interface{} {
	bs := make(map[interface{}]interface{})
	for _, m := range s {
		pair := m.(MapPair)
		bs[pair.Key] = pair.Value
	}
	return bs
}
