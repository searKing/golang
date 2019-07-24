package class

import (
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

func getFunctionFullName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
func GetFunctionName(i interface{}) string {
	nameFull := getFunctionFullName(i)
	nameEnd := filepath.Ext(nameFull)        // .foo-fm
	name := strings.TrimPrefix(nameEnd, ".") // foo-fm
	name = strings.Split(name, "-")[0]       // foo
	return name
}
