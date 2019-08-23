package defaults

import "reflect"

type convertFunc func(v reflect.Value, tag reflect.StructTag) (isUserDefined bool, err error)
