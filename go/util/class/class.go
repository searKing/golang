// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package class

// Classer is an interface to enhance go's embed struct with virtual member function in C++|Java
// It's used like this:
//import "github.com/searKing/golang/go/util/class"
//type Pet struct {
//	class.Class
//}
//
//func (pet *Pet) Name() string {
//	return "pet"
//}
//
//func (pet *Pet) VirtualUpperName() string {
//	p := pet.GetDerivedElse(pet).(Peter)
//	return strings.ToUpper(p.Name())
//}
//
//func (pet *Pet) UpperName() string {
//	return strings.ToUpper(pet.Name())
//}
//
////Dog derived from Pet
//type Dog struct {
//	Pet
//}
//
//func NewDogEmbedded() *Dog {
//	return &Dog{}
//}
//
//func NewDogDerived() *Dog {
//	dog := &Dog{}
//	dog.SetDerived(dog)
//	return dog
//}
//
//func (dog *Dog) Name() string {
//	return "dog"
//}
type Classer interface {
	GetDerived() Classer
	GetDerivedElse(defaultClass Classer) Classer
	SetDerived(derived Classer)
}

type Class struct {
	derived Classer
}

func NewClass() *Class {
	return &Class{}
}

// GetDerived returns actual outermost struct
func (task Class) GetDerived() Classer {
	return task.derived
}

// GetDerivedElse returns actual outermost struct if it is non-{@code nil} and
// otherwise returns the non-{@code nil} argument.
func (task Class) GetDerivedElse(defaultClass Classer) Classer {
	if task.derived == nil {
		return defaultClass
	}
	return task.derived
}

// SetDerived updates actual outermost struct
func (task *Class) SetDerived(derived Classer) {
	task.derived = derived
}
