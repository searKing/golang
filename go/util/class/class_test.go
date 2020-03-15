// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package class_test

import (
	"strings"
	"testing"

	"github.com/searKing/golang/go/util/class"
)

type Peter interface {
	Name() string
	VirtualUpperName() string
	UpperName() string
}

type Pet struct {
	class.Class
}

func (pet *Pet) Name() string {
	return "pet"
}

func (pet *Pet) VirtualUpperName() string {
	p := pet.GetDerivedElse(pet).(Peter)
	return strings.ToUpper(p.Name())
}

func (pet *Pet) UpperName() string {
	return strings.ToUpper(pet.Name())
}

//Dog derived from Pet
type Dog struct {
	Pet
}

func NewDogEmbedded() *Dog {
	return &Dog{}
}

func NewDogDerived() *Dog {
	dog := &Dog{}
	dog.SetDerived(dog)
	return dog
}

func (dog *Dog) Name() string {
	return "dog"
}

type ClassTests struct {
	input  Peter
	output []string
}

var classTests = []ClassTests{
	{
		input:  NewDogEmbedded(),
		output: []string{"dog", "PET", "PET"},
	},
	{
		input:  NewDogDerived(),
		output: []string{"dog", "DOG", "PET"},
	},
}

func TestClass(t *testing.T) {
	for n, test := range classTests {
		gots := []string{test.input.Name(), test.input.VirtualUpperName(), test.input.UpperName()}
		for i, got := range gots {
			if got != test.output[i] {
				t.Errorf("#%d[%d]: %v: got %s runs; expected %s", n, i, test.input, got, test.output[i])
			}
		}
	}
}
