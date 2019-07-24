package tag

import (
	"gopkg.in/yaml.v2"
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

func (thiz *Name) TagDefault() error {
	if *thiz == "" {
		*thiz = "Bob"
	}
	return nil
}
func TestTag(t *testing.T) {
	i := &inputType{}
	expect := &inputType{
		Name:        "Bob",
		Age:         10,
		IntArray:    []int{1, 2, 3},
		StringArray: []string{"stdout", "./logs"},
		Map:         map[string]string{"name": "Alice", "age": "18"},
	}
	err := Tag(i, func(val reflect.Value, tag reflect.StructTag) error {
		return yaml.Unmarshal([]byte(tag.Get("default")), val.Addr().Interface())
	})
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(i, expect) {
		t.Errorf("expect\n[\n%v\n]\nactual[\n%v\n]", expect, i)
	}
}
