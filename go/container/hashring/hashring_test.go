// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hashring

import (
	"runtime"
	"sort"
	"strconv"
	"testing"
	"testing/quick"
)

func TestNew(t *testing.T) {
	numReps := 160
	x := New(WithNumberNodeRepetitions(numReps))
	if x == nil {
		t.Errorf("expected obj")
		return
	}

	if x.numReps != numReps {
		t.Errorf("got %d, want %d", x.numReps, numReps)
	}
}

func TestAdd(t *testing.T) {
	numReps := 160
	x := New(WithNumberNodeRepetitions(numReps))
	x.AddKetamaNodes(StringNode("abcdefg"))

	if len(x.ketamaNodes) != numReps {
		t.Errorf("got %d, want %d", len(x.ketamaNodes), numReps)
	}
	if len(x.sortedKetamaHashRing) != numReps {
		t.Errorf("got %d, want %d", len(x.sortedKetamaHashRing), numReps)
	}
	if sort.IsSorted(x.sortedKetamaHashRing) == false {
		t.Errorf("expected sorted hashes to be sorted")
	}
	x.AddKetamaNodes(StringNode("qwer"))

	if len(x.ketamaNodes) != 2*numReps {
		t.Errorf("got %d, want %d", len(x.ketamaNodes), 2*numReps)
	}
	if len(x.sortedKetamaHashRing) != 2*numReps {
		t.Errorf("got %d, want %d", len(x.ketamaNodes), 2*numReps)
	}
	if sort.IsSorted(x.sortedKetamaHashRing) == false {
		t.Errorf("expected sorted hashes to be sorted")
	}
}

func TestRemove(t *testing.T) {
	numReps := 160
	x := New(WithNumberNodeRepetitions(numReps))
	x.AddKetamaNodes(StringNode("abcdefg"))
	x.RemoveKetamaNodes(StringNode("abcdefg"))
	if len(x.ketamaNodes) != 0 {
		t.Errorf("got %d, want %d", len(x.ketamaNodes), 0)
	}
	if len(x.sortedKetamaHashRing) != 0 {
		t.Errorf("got %d, want %d", len(x.ketamaNodes), 0)
	}
}

func TestRemoveNonExisting(t *testing.T) {
	numReps := 160
	x := New(WithNumberNodeRepetitions(numReps))
	x.AddKetamaNodes(StringNode("abcdefg"))
	x.RemoveKetamaNodes(StringNode("abcdefghijk"))
	if len(x.ketamaNodes) != numReps {
		t.Errorf("got %d, want %d", len(x.ketamaNodes), numReps)
	}
}

func TestGetEmpty(t *testing.T) {
	numReps := 160
	x := New(WithNumberNodeRepetitions(numReps))
	_, has := x.Get("asdfsadfsadf")
	if has {
		t.Errorf("expected error")
	}
}

func TestGetSingle(t *testing.T) {
	numReps := 160
	x := New(WithNumberNodeRepetitions(numReps))
	x.AddKetamaNodes(StringNode("abcdefg"))
	f := func(s string) bool {
		y, has := x.Get(s)
		if !has {
			return false
		}
		t.Logf("s = %q, y = %q", s, y)
		return y.GetSocketAddress().String() == "abcdefg"
	}
	if err := quick.Check(f, nil); err != nil {
		t.Logf("missing nodes")
	}
}

type gtest struct {
	in  string
	out string
}

var gmtests = []gtest{
	{"ggg", "abcdefg"},
	{"hhh", "opqrstu"},
	{"iii", "hijklmn"},
}

func TestGetMultiple(t *testing.T) {
	x := New()
	x.AddKetamaNodes(StringNode("abcdefg"))
	x.AddKetamaNodes(StringNode("hijklmn"))
	x.AddKetamaNodes(StringNode("opqrstu"))
	for i, v := range gmtests {
		result, has := x.Get(v.in)
		if !has {
			t.Fatal()
		}
		if result.GetSocketAddress().String() != v.out {
			t.Errorf("%d. got %q, expected %q", i, result, v.out)
		}
	}
}

func TestGetMultipleQuick(t *testing.T) {
	x := New()
	x.AddKetamaNodes(StringNode("abcdefg"))
	x.AddKetamaNodes(StringNode("hijklmn"))
	x.AddKetamaNodes(StringNode("opqrstu"))
	f := func(s string) bool {
		y, has := x.Get(s)
		if !has {
			return false
		}
		t.Logf("s = %q, y = %q", s, y)
		return y.GetSocketAddress().String() == "abcdefg" ||
			y.GetSocketAddress().String() == "hijklmn" ||
			y.GetSocketAddress().String() == "opqrstu"
	}
	if err := quick.Check(f, nil); err != nil {
		t.Logf("missing nodes")
	}
}

var rtestsBefore = []gtest{
	{"ggg", "abcdefg"},
	{"hhh", "opqrstu"},
	{"iii", "hijklmn"},
}

var rtestsAfter = []gtest{
	{"ggg", "abcdefg"},
	{"hhh", "opqrstu"},
	{"iii", "abcdefg"},
}

func TestGetMultipleRemove(t *testing.T) {
	x := New()
	x.AddKetamaNodes(StringNode("abcdefg"))
	x.AddKetamaNodes(StringNode("hijklmn"))
	x.AddKetamaNodes(StringNode("opqrstu"))
	for i, v := range rtestsBefore {
		result, has := x.Get(v.in)
		if !has {
			t.Fatal()
		}
		if result.GetSocketAddress().String() != v.out {
			t.Errorf("%d. got %q, expected %q before rm", i, result, v.out)
		}
	}
	x.RemoveKetamaNodes(StringNode("hijklmn"))
	for i, v := range rtestsAfter {
		result, has := x.Get(v.in)
		if !has {
			t.Fatal()
		}
		if result.GetSocketAddress().String() != v.out {
			t.Errorf("%d. got %q, expected %q after rm", i, result, v.out)
		}
	}
}

func TestGetMultipleRemoveQuick(t *testing.T) {
	x := New()
	x.AddKetamaNodes(StringNode("abcdefg"))
	x.AddKetamaNodes(StringNode("hijklmn"))
	x.AddKetamaNodes(StringNode("opqrstu"))
	x.RemoveKetamaNodes(StringNode("opqrstu"))
	f := func(s string) bool {
		y, has := x.Get(s)
		if !has {
			t.Logf("missing node")
			return false
		}
		t.Logf("s = %q, y = %q", s, y)
		return y.GetSocketAddress().String() == "abcdefg" || y.GetSocketAddress().String() == "hijklmn"
	}
	if err := quick.Check(f, nil); err != nil {
		t.Logf("missing nodes")
	}
}

func TestGetTwo(t *testing.T) {
	x := New()
	x.AddKetamaNodes(StringNode("abcdefg"))
	x.AddKetamaNodes(StringNode("hijklmn"))
	x.AddKetamaNodes(StringNode("opqrstu"))
	a, b, has := x.GetTwo("99999999")
	if !has {
		t.Fatal("missing nodes")
	}
	if a == b {
		t.Errorf("a shouldn't equal b")
	}
	if a.GetSocketAddress().String() != "opqrstu" {
		t.Errorf("wrong a: %q", a)
	}
	if b.GetSocketAddress().String() != "hijklmn" {
		t.Errorf("wrong b: %q", b)
	}
}

func TestGetTwoQuick(t *testing.T) {
	x := New()
	x.AddKetamaNodes(StringNode("abcdefg"))
	x.AddKetamaNodes(StringNode("hijklmn"))
	x.AddKetamaNodes(StringNode("opqrstu"))
	f := func(s string) bool {
		a, b, has := x.GetTwo(s)
		if !has {
			t.Logf("missing nodes")
			return false
		}
		if a == b {
			t.Logf("a == b")
			return false
		}
		if a.GetSocketAddress().String() != "abcdefg" &&
			a.GetSocketAddress().String() != "hijklmn" &&
			a.GetSocketAddress().String() != "opqrstu" {
			t.Logf("invalid a: %q", a)
			return false
		}

		if b.GetSocketAddress().String() != "abcdefg" &&
			b.GetSocketAddress().String() != "hijklmn" &&
			b.GetSocketAddress().String() != "opqrstu" {
			t.Logf("invalid b: %q", b)
			return false
		}
		return true
	}
	if err := quick.Check(f, nil); err != nil {
		t.Logf("missing nodes")
	}
}

func TestGetTwoOnlyTwoQuick(t *testing.T) {
	x := New()
	x.AddKetamaNodes(StringNode("abcdefg"))
	x.AddKetamaNodes(StringNode("hijklmn"))
	f := func(s string) bool {
		a, b, has := x.GetTwo(s)
		if !has {
			t.Logf("missing nodes")
			return false
		}
		if a == b {
			t.Logf("a == b")
			return false
		}
		if a.GetSocketAddress().String() != "abcdefg" && a.GetSocketAddress().String() != "hijklmn" {
			t.Logf("invalid a: %q", a)
			return false
		}

		if b.GetSocketAddress().String() != "abcdefg" && b.GetSocketAddress().String() != "hijklmn" {
			t.Logf("invalid b: %q", b)
			return false
		}
		return true
	}
	if err := quick.Check(f, nil); err != nil {
		t.Logf("missing nodes")
	}
}

func TestGetTwoOnlyOneInCircle(t *testing.T) {
	x := New()

	x.AddKetamaNodes(StringNode("abcdefg"))
	a, b, has := x.GetTwo("99999999")
	if !has {
		t.Logf("missing nodes")
	}
	if a == b {
		t.Errorf("a shouldn't equal b")
	}
	if a.GetSocketAddress().String() != "abcdefg" {
		t.Errorf("wrong a: %q", a)
	}
	if b != nil {
		t.Errorf("wrong b: %q", b)
	}
}

func TestGetN(t *testing.T) {
	x := New()
	x.AddKetamaNodes(StringNode("abcdefg"))
	x.AddKetamaNodes(StringNode("hijklmn"))
	x.AddKetamaNodes(StringNode("opqrstu"))
	members, has := x.GetN("9999999", 3)
	if !has {
		t.Logf("missing nodes")
	}
	if len(members) != 3 {
		t.Errorf("expected 3 allNodes instead of %d", len(members))
	}
	if members[0].GetSocketAddress().String() != "abcdefg" {
		t.Errorf("wrong allNodes[0]: %q", members[0])
	}
	if members[1].GetSocketAddress().String() != "opqrstu" {
		t.Errorf("wrong allNodes[1]: %q", members[1])
	}
	if members[2].GetSocketAddress().String() != "hijklmn" {
		t.Errorf("wrong allNodes[2]: %q", members[2])
	}
}

func TestGetNLess(t *testing.T) {
	x := New()
	x.AddKetamaNodes(StringNode("abcdefg"))
	x.AddKetamaNodes(StringNode("hijklmn"))
	x.AddKetamaNodes(StringNode("opqrstu"))
	members, has := x.GetN("99999999", 2)
	if !has {
		t.Logf("missing nodes")
	}
	if len(members) != 2 {
		t.Errorf("expected 2 allNodes instead of %d", len(members))
	}
	if members[0].GetSocketAddress().String() != "opqrstu" {
		t.Errorf("wrong allNodes[0]: %q", members[0])
	}
	if members[1].GetSocketAddress().String() != "hijklmn" {
		t.Errorf("wrong allNodes[1]: %q", members[1])
	}
}

func TestGetNMore(t *testing.T) {
	x := New()
	x.AddKetamaNodes(StringNode("abcdefg"))
	x.AddKetamaNodes(StringNode("hijklmn"))
	x.AddKetamaNodes(StringNode("opqrstu"))
	members, has := x.GetN("9999999", 5)
	if !has {
		t.Logf("missing nodes")
	}
	if len(members) != 3 {
		t.Errorf("expected 3 allNodes instead of %d", len(members))
	}
	if members[0].GetSocketAddress().String() != "abcdefg" {
		t.Errorf("wrong allNodes[0]: %q", members[0])
	}
	if members[1].GetSocketAddress().String() != "opqrstu" {
		t.Errorf("wrong allNodes[1]: %q", members[1])
	}
	if members[2].GetSocketAddress().String() != "hijklmn" {
		t.Errorf("wrong allNodes[2]: %q", members[2])
	}
}

func TestGetNQuick(t *testing.T) {
	x := New()
	x.AddKetamaNodes(StringNode("abcdefg"))
	x.AddKetamaNodes(StringNode("hijklmn"))
	x.AddKetamaNodes(StringNode("opqrstu"))
	f := func(s string) bool {
		members, has := x.GetN(s, 3)
		if !has {
			t.Logf("missing nodes")
			return false
		}
		if len(members) != 3 {
			t.Logf("expected 3 allNodes instead of %d", len(members))
			return false
		}
		set := make(map[string]bool, 4)
		for _, member := range members {
			if set[member.GetSocketAddress().String()] {
				t.Logf("duplicate error")
				return false
			}
			set[member.GetSocketAddress().String()] = true
			if member.GetSocketAddress().String() != "abcdefg" &&
				member.GetSocketAddress().String() != "hijklmn" &&
				member.GetSocketAddress().String() != "opqrstu" {
				t.Logf("invalid member: %q", member)
				return false
			}
		}
		return true
	}
	if err := quick.Check(f, nil); err != nil {
		t.Logf("missing nodes")
	}
}

func TestGetNLessQuick(t *testing.T) {
	x := New()
	x.AddKetamaNodes(StringNode("abcdefg"))
	x.AddKetamaNodes(StringNode("hijklmn"))
	x.AddKetamaNodes(StringNode("opqrstu"))
	f := func(s string) bool {
		members, has := x.GetN(s, 2)
		if !has {
			t.Logf("missing nodes")
			return false
		}
		if len(members) != 2 {
			t.Logf("expected 2 allNodes instead of %d", len(members))
			return false
		}
		set := make(map[string]bool, 4)
		for _, member := range members {
			if set[member.GetSocketAddress().String()] {
				t.Logf("duplicate error")
				return false
			}
			set[member.GetSocketAddress().String()] = true
			if member.GetSocketAddress().String() != "abcdefg" &&
				member.GetSocketAddress().String() != "hijklmn" &&
				member.GetSocketAddress().String() != "opqrstu" {
				t.Logf("invalid member: %q", member)
				return false
			}
		}
		return true
	}
	if err := quick.Check(f, nil); err != nil {
		t.Logf("missing nodes")
	}
}

func TestGetNMoreQuick(t *testing.T) {
	x := New()
	x.AddKetamaNodes(StringNode("abcdefg"))
	x.AddKetamaNodes(StringNode("hijklmn"))
	x.AddKetamaNodes(StringNode("opqrstu"))
	f := func(s string) bool {
		members, has := x.GetN(s, 5)
		if !has {
			t.Logf("missing nodes")
			return false
		}
		if len(members) != 3 {
			t.Logf("expected 3 allNodes instead of %d", len(members))
			return false
		}
		set := make(map[string]bool, 4)
		for _, member := range members {
			if set[member.GetSocketAddress().String()] {
				t.Logf("duplicate error")
				return false
			}
			set[member.GetSocketAddress().String()] = true
			if member.GetSocketAddress().String() != "abcdefg" && member.GetSocketAddress().String() != "hijklmn" && member.GetSocketAddress().String() != "opqrstu" {
				t.Logf("invalid member: %q", member)
				return false
			}
		}
		return true
	}
	if err := quick.Check(f, nil); err != nil {
		t.Logf("missing nodes")
	}
}

func TestSet(t *testing.T) {
	x := New()
	x.AddKetamaNodes(StringNode("abc"))
	x.AddKetamaNodes(StringNode("def"))
	x.AddKetamaNodes(StringNode("ghi"))
	x.SetKetamaNodes(StringNode("jkl"), StringNode("mno"))
	if len(x.allNodes) != 2 {
		t.Errorf("expected 2 elts, got %d", len(x.allNodes))
	}
	a, b, has := x.GetTwo("qwerqwerwqer")
	if !has {
		t.Fatal()
	}
	if a.GetSocketAddress().String() != "jkl" && a.GetSocketAddress().String() != "mno" {
		t.Errorf("expected jkl or mno, got %s", a)
	}
	if b.GetSocketAddress().String() != "jkl" && b.GetSocketAddress().String() != "mno" {
		t.Errorf("expected jkl or mno, got %s", b)
	}
	if a == b {
		t.Errorf("expected a != b, they were both %s", a)
	}
	x.SetKetamaNodes(StringNode("jkl"), StringNode("mno"))
	if len(x.allNodes) != 2 {
		t.Errorf("expected 2 elts, got %d", len(x.allNodes))
	}
	a, b, has = x.GetTwo("qwerqwerwqer")
	if !has {
		t.Fatal()
	}
	if a.GetSocketAddress().String() != "pqr" && a.GetSocketAddress().String() != "mno" {
		t.Errorf("expected jkl or mno, got %s", a)
	}
	if b.GetSocketAddress().String() != "pqr" && b.GetSocketAddress().String() != "mno" {
		t.Errorf("expected jkl or mno, got %s", b)
	}
	if a == b {
		t.Errorf("expected a != b, they were both %s", a)
	}
	x.SetKetamaNodes(StringNode("pqr"), StringNode("mno"))
	if len(x.allNodes) != 2 {
		t.Errorf("expected 2 elts, got %d", len(x.allNodes))
	}
	a, b, has = x.GetTwo("qwerqwerwqer")
	if !has {
		t.Fatal()
	}
	if a.GetSocketAddress().String() != "pqr" && a.GetSocketAddress().String() != "mno" {
		t.Errorf("expected jkl or mno, got %s", a)
	}
	if b.GetSocketAddress().String() != "pqr" && b.GetSocketAddress().String() != "mno" {
		t.Errorf("expected jkl or mno, got %s", b)
	}
	if a == b {
		t.Errorf("expected a != b, they were both %s", a)
	}
}

// allocBytes returns the number of bytes allocated by invoking f.
func allocBytes(f func()) uint64 {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	t := stats.TotalAlloc
	f()
	runtime.ReadMemStats(&stats)
	return stats.TotalAlloc - t
}

func mallocNum(f func()) uint64 {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	t := stats.Mallocs
	f()
	runtime.ReadMemStats(&stats)
	return stats.Mallocs - t
}

func BenchmarkAllocations(b *testing.B) {
	x := New()
	x.AddKetamaNodes(StringNode("stays"))
	b.ResetTimer()
	allocSize := allocBytes(func() {
		for i := 0; i < b.N; i++ {
			x.AddKetamaNodes(StringNode("Foo"))
			x.RemoveKetamaNodes(StringNode("Foo"))
		}
	})
	b.Logf("%d: Allocated %d bytes (%.2fx)", b.N, allocSize, float64(allocSize)/float64(b.N))
}

func BenchmarkMalloc(b *testing.B) {
	x := New()
	x.AddKetamaNodes(StringNode("stays"))
	b.ResetTimer()
	mallocs := mallocNum(func() {
		for i := 0; i < b.N; i++ {
			x.AddKetamaNodes(StringNode("Foo"))
			x.RemoveKetamaNodes(StringNode("Foo"))
		}
	})
	b.Logf("%d: Mallocd %d times (%.2fx)", b.N, mallocs, float64(mallocs)/float64(b.N))
}

func BenchmarkCycle(b *testing.B) {
	x := New()
	x.AddKetamaNodes(StringNode("nothing"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x.AddKetamaNodes(StringNode("foo" + strconv.Itoa(i)))
		x.RemoveKetamaNodes(StringNode("foo" + strconv.Itoa(i)))
	}
}

func BenchmarkCycleLarge(b *testing.B) {
	x := New()
	for i := 0; i < 10; i++ {
		x.AddKetamaNodes(StringNode("start" + strconv.Itoa(i)))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x.AddKetamaNodes(StringNode("foo" + strconv.Itoa(i)))
		x.RemoveKetamaNodes(StringNode("foo" + strconv.Itoa(i)))
	}
}

func BenchmarkGet(b *testing.B) {
	x := New()
	x.AddKetamaNodes(StringNode("nothing"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x.Get("nothing")
	}
}

func BenchmarkGetLarge(b *testing.B) {
	x := New()
	for i := 0; i < 10; i++ {
		x.AddKetamaNodes(StringNode("start" + strconv.Itoa(i)))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x.Get("nothing")
	}
}

func BenchmarkGetN(b *testing.B) {
	x := New()
	x.AddKetamaNodes(StringNode("nothing"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x.GetN("nothing", 3)
	}
}

func BenchmarkGetNLarge(b *testing.B) {
	x := New()
	for i := 0; i < 10; i++ {
		x.AddKetamaNodes(StringNode("start" + strconv.Itoa(i)))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x.GetN("nothing", 3)
	}
}

func BenchmarkGetTwo(b *testing.B) {
	x := New()
	x.AddKetamaNodes(StringNode("nothing"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x.GetTwo("nothing")
	}
}

func BenchmarkGetTwoLarge(b *testing.B) {
	x := New()
	for i := 0; i < 10; i++ {
		x.AddKetamaNodes(StringNode("start" + strconv.Itoa(i)))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x.GetTwo("nothing")
	}
}

// from @edsrzf on github:
func TestAddCollision(t *testing.T) {
	// These two strings produce several crc32 collisions after "|i" is
	// appended added by KetamaNodeLocator.virtualNode.
	const s1 = "abear"
	const s2 = "solidiform"
	x := New()
	x.AddKetamaNodes(StringNode(s1))
	x.AddKetamaNodes(StringNode(s2))
	elt1, has := x.Get("abear")
	if !has {
		t.Fatal("missing node")
	}

	y := New()
	// add elements in opposite order
	y.AddKetamaNodes(StringNode(s2))
	y.AddKetamaNodes(StringNode(s1))
	elt2, has := y.Get(s1)
	if !has {
		t.Fatal("missing node")
	}

	if elt1 != elt2 {
		t.Error(elt1, "and", elt2, "should be equal")
	}
}
