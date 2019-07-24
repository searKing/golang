package object

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
		panic(fmt.Errorf("The number of %s's params is %d, not adapted to %d.", name, len(args), fn.Type().NumIn()))
		return nil
	}
	return transOut(fn.Call(transIn(args...)))
}
func (c *Class) Invoke(method interface{}, args ...interface{}) []interface{} {
	name := GetFunctionName(method)
	return c.InvokeByName(name, args...)
}

// prototype - 将被创建对象的原型
// properties (可选) - 将被创建对象需要添加的属性
//
// throws TypeError 如果 'prototype' 参数不是一个对象，也不是null
//
// returns 新创建的对象
var classType = reflect.TypeOf((*properties)(nil)).Elem()

var TypeError = func(prototype interface{}) error {
	return fmt.Errorf("%v is not a Class", reflect.TypeOf(prototype).String())
}

type Father struct {
	Class
	Name string
}

func (*Father) Say() {
	fmt.Printf("Say Father\n")
}
func (Father) Speak() {
	fmt.Printf("Speak Father\n")
}

// public Son : public Father
// Son需要继承自Father，而需要对Father做访问权限控制，则想到增加代理类，如同智能指针一样， class就是代理类
// MakeVirtual 为son的构造函数
type Son struct {
	Father
}

func (*Son) Say(name string) {
	fmt.Printf("Say Son :%s\n", name)
}
func (*Son) Hello(string) {
	fmt.Printf("Say Son\n")
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
