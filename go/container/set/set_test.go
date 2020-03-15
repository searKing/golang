// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package set

import (
	"testing"
)

func Test(t *testing.T) {
	s := New()

	if s.Count() != 0 {
		t.Errorf("Length should be 0")
	}
	s.Remove(0)
	s.Add(5)

	if s.Count() != 1 {
		t.Errorf("Length should be 1")
	}

	if !s.Contains(5) {
		t.Errorf("Membership test failed")
	}

	s.Remove(5)

	if s.Count() != 0 {
		t.Errorf("Length should be 0")
	}

	if s.Contains(5) {
		t.Errorf("The set should be empty")
	}

	// Difference
	s1 := Of(1, 2, 3, 4, 5, 6)
	s2 := Of(4, 5, 6)
	s3 := s1.RemoveAll(s2.ToStream())

	if s3.Length() != 3 {
		t.Errorf("Length should be 3")
	}

	if !(s3.Contains(1) && s3.Contains(2) && s3.Contains(3)) {
		t.Errorf("Set should only contain 1, 2, 3")
	}

	s1 = Of(1, 2, 3, 4, 5, 6)
	// Intersection
	s3 = s1.RetainAll(s2.ToStream())
	if s3.Count() != 3 {
		t.Errorf("Length should be 3 after RetainAll")
	}

	if !(s3.Contains(4) && s3.Contains(5) && s3.Contains(6)) {
		t.Errorf("Set should contain 4, 5, 6")
	}

	// Union
	s4 := Of(7, 8, 9)
	s3 = s2.AddAll(s4.ToStream())

	if s3.Length() != 6 {
		t.Errorf("Length should be 6 after union")
	}

	if !(s3.Contains(7)) {
		t.Errorf("Set should contain 4, 5, 6, 7, 8, 9")
	}

	// Subset
	if !s1.ContainsAll(s1.ToStream()) {
		t.Errorf("set should be a subset of itself")
	}

}
