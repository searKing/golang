// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package boolean

import "testing"

type input struct {
	a      bool
	b      bool
	c      []bool
	expect bool
}

func TestAND(t *testing.T) {
	ins := []input{
		{
			a:      true,
			b:      true,
			expect: true,
		},
		{
			a:      true,
			b:      false,
			expect: false,
		},
		{
			a:      false,
			b:      true,
			expect: false,
		},
		{
			a:      false,
			b:      false,
			expect: false,
		},
		{
			a: true,
			b: true,
			c: []bool{
				true,
			},
			expect: true,
		},
		{
			a: true,
			b: true,
			c: []bool{
				false,
				false,
			},
			expect: false,
		},
	}
	for idx, in := range ins {
		if AND(in.a, in.b, in.c...) != in.expect {
			t.Errorf("#%d expect %t", idx, in.expect)
		}
	}
}

func TestOR(t *testing.T) {
	ins := []input{
		{
			a:      true,
			b:      true,
			expect: true,
		},
		{
			a:      true,
			b:      false,
			expect: true,
		},
		{
			a:      false,
			b:      true,
			expect: true,
		},
		{
			a:      false,
			b:      false,
			expect: false,
		},
		{
			a: false,
			b: false,
			c: []bool{
				true,
			},
			expect: true,
		},
		{
			a: false,
			b: false,
			c: []bool{
				false,
				true,
			},
			expect: true,
		},
	}
	for idx, in := range ins {
		if OR(in.a, in.b, in.c...) != in.expect {
			t.Errorf("#%d expect %t", idx, in.expect)
		}
	}
}

func TestXOR(t *testing.T) {
	ins := []input{
		{
			a:      true,
			b:      true,
			expect: false,
		},
		{
			a:      true,
			b:      false,
			expect: true,
		},
		{
			a:      false,
			b:      true,
			expect: true,
		},
		{
			a:      false,
			b:      false,
			expect: false,
		},
		{
			a: false,
			b: false,
			c: []bool{
				true,
			},
			expect: true,
		},
		{
			a: false,
			b: false,
			c: []bool{
				false,
				true,
			},
			expect: true,
		},
	}
	for idx, in := range ins {
		if XOR(in.a, in.b, in.c...) != in.expect {
			t.Errorf("#%d expect %t", idx, in.expect)
		}
	}
}

func TestXNOR(t *testing.T) {
	ins := []input{
		{
			a:      true,
			b:      true,
			expect: true,
		},
		{
			a:      true,
			b:      false,
			expect: false,
		},
		{
			a:      false,
			b:      true,
			expect: false,
		},
		{
			a:      false,
			b:      false,
			expect: true,
		},
		{
			a: false,
			b: false,
			c: []bool{
				true,
			},
			expect: true,
		},
		{
			a: false,
			b: false,
			c: []bool{
				false,
				true,
			},
			expect: false,
		},
	}
	for idx, in := range ins {
		if XNOR(in.a, in.b, in.c...) != in.expect {
			t.Errorf("#%d expect %t", idx, in.expect)
		}
	}
}
func TestBoolFunc(t *testing.T) {
	ins := []input{
		{
			a:      true,
			b:      true,
			expect: true,
		},
		{
			a:      true,
			b:      false,
			expect: false,
		},
		{
			a:      false,
			b:      true,
			expect: false,
		},
		{
			a:      false,
			b:      false,
			expect: false,
		},
		{
			a: true,
			b: true,
			c: []bool{
				false,
			},
			expect: false,
		},
		{
			a: true,
			b: true,
			c: []bool{
				true,
				false,
			},
			expect: false,
		},
	}
	allIsTrue := func(a, b bool) bool {
		if !a || !b {
			return false
		}
		return true
	}
	for idx, in := range ins {
		if BoolFunc(in.a, in.b, allIsTrue, in.c...) != in.expect {
			t.Errorf("#%d expect %t", idx, in.expect)
		}
	}
}
