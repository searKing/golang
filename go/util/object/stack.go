package object

import (
	"runtime"
)

func GetFunc() *runtime.Func {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return nil
	}
	return runtime.FuncForPC(pc)
}
