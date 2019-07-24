package object

import (
	"fmt"
	"reflect"
	"testing"
)

//https://wizardforcel.gitbooks.io/gopl-zh/ch12/ch12-05.html
func TestClass(t *testing.T) {
	// 调用s的构造函数
	//var s = MakeVirtual(Son{}, Father{}, nil).(Son)
	//var f Father = s.Father
	son := &Son{}
	//GetMembers(son)
	//father := &Father{}
	//son.Say()
	//fmt.Printf("after set")
	//s := son.Say
	//son.Say()
	//_ = s
	//fmt.Printf("son.say = %v\n", son.Say)
	////father := &Father{}
	//GetMembers(father)
	//father.Say()
	son.Name = "searKing"

	//son.Father.Invoke(son.Hello)
	MakeVirtual(son, nil)
	son.Invoke(son.Say, "Bitch")
	son.Father.Invoke(son.Say, "Bitch")
	son.Father.Invoke(son.Speak)
}

// 该测试用例验证了Receiver函数是不可以被重写的，所以，我们只能采用寻找的方式来解决虚函数的问题
// 寻找的话，需要实现切片编程，这也是不可行的
// 所以解决方案只有一个，就是将子类转为基类时，提供一个转换操作，返回的其实是子类的代理，该代理就是之前想要更改的父类的影子迂回
// 该代理需要实现自动实现父类所有的虚函数，我们称之为父类切片函数，欺骗函数中，完成子类调用寻找，或者干脆就直接返回子类的调用，这样简洁方便
// 现在问题就转换为如何实现父类的代理类
func TestMethodVirtualNotAccessable(t *testing.T) {
	// 调用s的构造函数
	//var s = MakeVirtual(Son{}, Father{}, nil).(Son)
	//var f Father = s.Father
	son := &Son{}
	tv := reflect.ValueOf(son).Elem()
	structType := tv.Type()
	for i := 0; i < structType.NumField(); i++ {
		tf := tv.Field(i)
		hookType := tf.Type()
		fmt.Printf("Field %d: Value=%v Kind=%v\n", i, tf.String(), tf.Kind().String())
		if hookType.Kind() != reflect.Func {
			continue
		}
		//if tf.IsNil() {
		//	tf.Set(of)
		//	continue
		//}
		//
		//// Make a copy of tf for tf to call. (Otherwise it
		//// creates a recursive call cycle and stack overflows)
		//tfCopy := reflect.ValueOf(tf.Interface())
		//
		//// We need to call both tf and of in some order.
		//newFunc := reflect.MakeFunc(hookType, func(args []reflect.Value) []reflect.Value {
		//	tfCopy.Call(args)
		//	return of.Call(args)
		//})
		//tv.Field(i).Set(newFunc)
	}
	for i := 0; i < structType.NumMethod(); i++ {
		tm := tv.Method(i)
		fmt.Printf("Method %d: Value=%v Kind=%v\n", i, tm.String(), tm.Kind().String())
	}
	GetMembers(son)
	son.Say("Son")
	father := &Father{}
	GetMembers(father)
	father.Say()
	fmt.Printf("NumMethod of son is %v\n", reflect.ValueOf(son).Elem().NumMethod())
	//sayMethod := reflect.ValueOf(son).Elem().MethodByName("Say").Elem()
	_ = reflect.ValueOf(son).MethodByName("Say")
	//fmt.Printf("sayMethod can addr = %v\n", reflect.ValueOf(son.Say).Elem().CanAddr())

	//fmt.Printf("sayMethod can addr = %v\n", reflect.ValueOf(son).MethodByName("Say").Elem().CanAddr())

	//sayMethod.Set(reflect.ValueOf(father).Elem().Method(0))
	//son.Say()

}
