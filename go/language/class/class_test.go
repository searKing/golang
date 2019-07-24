package class

import (
	"fmt"
	"testing"
)

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
func TestClass(t *testing.T) {

}
