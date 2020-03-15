// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package object

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
)

// 获取一个对象的字段和方法
func DumpMethodInfo(i interface{}) []byte {
	var b bytes.Buffer

	// 获取 i 的类型信息
	t := reflect.TypeOf(i)

	for {
		// 进一步获取 i 的类别信息
		if t.Kind() == reflect.Struct {
			// 只有结构体可以获取其字段信息
			_, _ = fmt.Fprintf(&b, "\n%-8v %v 个字段:\n", t, t.NumField())
			// 进一步获取 i 的字段信息
			for i := 0; i < t.NumField(); i++ {
				fmt.Println(t.Field(i).Name)
			}
		}
		// 任何类型都可以获取其方法信息
		_, _ = fmt.Fprintf(&b, "\n%-8v %v 个方法:\n", t, t.NumMethod())
		// 进一步获取 i 的方法信息
		for i := 0; i < t.NumMethod(); i++ {
			_, _ = fmt.Fprintln(&b, t.Method(i).Name)
		}
		if t.Kind() == reflect.Ptr {
			// 如果是指针，则获取其所指向的元素
			t = t.Elem()
		} else {
			// 否则上面已经处理过了，直接退出循环
			break
		}
	}
	return b.Bytes()
}

// DumpFuncInfo returns m's function info
func DumpFuncInfo(m interface{}) []byte {
	//Reflection type of the underlying data of the interface
	x := reflect.TypeOf(m)

	numIn := x.NumIn()   //Count inbound parameters
	numOut := x.NumOut() //Count outbounding parameters
	var b bytes.Buffer

	_, _ = fmt.Fprintln(&b, "Method:", x.String())
	_, _ = fmt.Fprintln(&b, "Variadic:", x.IsVariadic()) // Used (<type> ...) ?
	_, _ = fmt.Fprintln(&b, "Package:", x.PkgPath())

	for i := 0; i < numIn; i++ {

		inV := x.In(i)
		inKind := inV.Kind() //func
		_, _ = fmt.Fprintf(&b, "\nParameter IN: "+strconv.Itoa(i)+"\nKind: %v\nName: %v\n-----------", inKind, inV.Name())
	}
	for o := 0; o < numOut; o++ {

		returnV := x.Out(0)
		returnKind := returnV.Kind()
		_, _ = fmt.Fprintf(&b, "\nParameter OUT: "+strconv.Itoa(o)+"\nKind: %v\nName: %v\n", returnKind, returnV.Name())
	}
	return b.Bytes()
}
