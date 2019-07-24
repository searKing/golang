package _default

import (
	"fmt"
	"reflect"
	"testing"
)

type inputType struct {
	Name        Name              `default:"Alice"`
	Age         int               `default:"10"`
	IntArray    []int             `default:"[1,2,3]"`
	StringArray []string          `default:"[\"stdout\",\"./logs\"]"`
	Map         map[string]string `default:"{\"name\": \"Alice\", \"age\": 18}"`
}
type Name string

func (thiz *Name) ConvertDefault(_ reflect.Value, _ reflect.StructTag) error {
	if *thiz == "" {
		*thiz = "Bob"
	}
	return nil
}

func (thiz *Name) Hello() error {
	fmt.Printf("Hello\n")
	return nil
}

func TestConvert(t *testing.T) {
	i := &inputType{}
	expect := &inputType{
		Name:        "Bob",
		Age:         10,
		IntArray:    []int{1, 2, 3},
		StringArray: []string{"stdout", "./logs"},
		Map:         map[string]string{"name": "Alice", "age": "18"},
	}
	err := Convert(i)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(i, expect) {
		t.Errorf("expect\n[\n%v\n]\nactual[\n%v\n]", expect, i)
	}
	//reflect.ValueOf(&i.Name).MethodByName("Hello").Call(nil)
	//reflect.ValueOf(&i.Name).MethodByName("ConvertDefault").Call(nil)

	//fmt.Printf("converterType.Name %v %v\n", converterType.Name(), converterType.Method(0).Name)

	//fmt.Printf("Name implement \n 1: %v 2: %v\n",
	//	reflect.ValueOf(&i.Name).MethodByName(converterType.Name()).
	//	Call([]reflect.Value{reflect.ValueOf(reflect.Value{}), reflect.ValueOf(reflect.StructField{})}),converterType.Name())

}
