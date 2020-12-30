package pragma_test

import (
	"fmt"
	"github.com/searKing/golang/go/pragma"
)

type NoUnkeyedLiterals struct {
	pragma.NoUnkeyedLiterals
	Name string
}

func ExampleNoUnkeyedLiterals() {
	//var a = NoUnkeyedLiterals{"Name"}
	//_ = a
	// compile error will print if uncomment codes above
	fmt.Println(`cannot convert "Name" (untyped string constant) to pragma.NoUnkeyedLiterals`)
	// Output: cannot convert "Name" (untyped string constant) to pragma.NoUnkeyedLiterals
}

type DoNotImplement interface {
	pragma.DoNotImplement
	String() string
}

type DoNotImplementStruct struct {
}

func (*DoNotImplementStruct) ProtoInternal(pragma.DoNotImplement) {

}

func (*DoNotImplementStruct) String() string {
	return "whoops"
}

func ExampleDoNotImplement() {
	var a DoNotImplementStruct
	_ = a

	// You can never implement this interface, with pragma.DoNotImplement embed
	//var b DoNotImplement = &a
	// go vet error will print if uncomment codes above
	fmt.Println("cannot use &a (type *DoNotImplementStruct) as type DoNotImplement in assignment:")

	// Output: cannot use &a (type *DoNotImplementStruct) as type DoNotImplement in assignment:
}

type DoNotCopy struct {
	pragma.DoNotCopy
}

func ExampleDoNotCopy() {
	//var a DoNotCopy
	//b := a
	//_ = b
	// go vet error will print if uncomment codes above
	fmt.Println("Assignment copies lock value to '_': type 'DoNotCopy' contains 'sync.Mutex' which is 'sync.Locker'")

	// Output: Assignment copies lock value to '_': type 'DoNotCopy' contains 'sync.Mutex' which is 'sync.Locker'
}

type DoNotCompare struct {
	age int

	pragma.DoNotCompare
}

func ExampleDoNotCompare() {
	// var a, b DoNotCompare
	// a == b
	// compile error will print if uncomment codes above
	fmt.Println("Invalid operation: a == b (operator == is not defined on DoNotCompare)")
	// Output: Invalid operation: a == b (operator == is not defined on DoNotCompare)
}

type CopyChecker struct {
	pragma.CopyChecker
}

func ExampleCopyChecker() {
	var a, b CopyChecker
	fmt.Printf("a copied : %t\n", a.Copied())
	fmt.Printf("a copied : %t\n", a.Copied())
	b = a
	fmt.Printf("a copied : %t\n", a.Copied())
	fmt.Printf("b copied : %t\n", b.Copied())

	// Output:
	// a copied : false
	// a copied : false
	// a copied : false
	// b copied : true
}
