package class

import (
	"fmt"
	"reflect"
)

// 类的属性
type prototype interface {
	SetProto(__proto__ prototype)
	SetProperities(properties properties)
	Invoke(method interface{}, args ...interface{}) []interface{}
}
type properties interface {
}

// 抽象类，因为interface没有成员变量
// 通过class，将外部结构体转为接口
type Class struct {
	__proto__ prototype // 虚函数表
	properties
	prototype
}

func (c *Class) SetProto(__proto__ prototype) {
	c.__proto__ = __proto__
}
func (c *Class) SetProperities(properties properties) {
	c.properties = properties
}

// 虚函数表递归
func (c *Class) findVirtualMethodByName(name string) reflect.Value {
	if c.__proto__ == nil {
		return reflect.Value{}
	}
	return reflect.ValueOf(c.__proto__).MethodByName(name)
}
func transIn(args ...interface{}) []reflect.Value {
	var values []reflect.Value
	for _, arg := range args {
		values = append(values, reflect.ValueOf(arg))
	}
	return values
}
func transOut(args []reflect.Value) []interface{} {
	var values []interface{}
	for _, arg := range args {
		values = append(values, arg.Interface())
	}
	return values
}

// call virtual method by name
// 后续增加虚函数表缓存策略
func (c *Class) InvokeByName(name string, args ...interface{}) []interface{} {
	fn := c.findVirtualMethodByName(name)

	if len(args) != fn.Type().NumIn() {
		panic(fmt.Errorf("the number of %s's params is %d, not adapted to %d", name, len(args), fn.Type().NumIn()))
		return nil
	}
	return transOut(fn.Call(transIn(args...)))
}
func (c *Class) Invoke(method interface{}, args ...interface{}) []interface{} {
	name := GetFunctionName(method)
	return c.InvokeByName(name, args...)
}

var classType = reflect.TypeOf((*properties)(nil)).Elem()

var TypeError = func(prototype interface{}) error {
	return fmt.Errorf("%v is not a Class", reflect.TypeOf(prototype).String())
}

var MakeVirtual = func(__proto__ interface{}, properties properties) {
	makeVirtual(__proto__.(prototype), properties)
}

// 使用interface来规避虚函数问题
// 传入一个旧的，生成一个新的
var makeVirtual = func(__proto__ prototype, properties properties) {
	if !reflect.ValueOf(__proto__).Type().Implements(classType) {
		panic(TypeError(__proto__))
	}
	__proto__.SetProperities(properties)
	////设置原型
	__proto__.SetProto(__proto__)
}

//// 虚函数表递归
//func resetAllClass(__proto__ prototype, properties properties, i interface{}) {
//	// 获取 i 的类型信息
//	t := reflect.TypeOf(i)
//	// 进一步获取 i 的类别信息
//	if t.Kind() == reflect.Struct {
//		// 只有结构体可以获取其字段信息
//		fmt.Printf("\n%-8v %v 个字段:\n", t, t.NumField())
//		// 进一步获取 i 的字段信息
//		if subT, ok := t.FieldByName("Class"); ok {
//			__proto__.SetProperities(properties)
//			////设置原型
//			__proto__.SetProto(__proto__)
//			resetAllClass(subT.Type)
//			return
//		}
//	}
//	return
//}
