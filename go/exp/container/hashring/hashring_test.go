// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hashring

import (
	"runtime"
	"slices"
	"strconv"
	"testing"
	"testing/quick"
)

func TestNew(t *testing.T) {
	numReps := 160
	x := New[string](WithNodeLocatorNumReps[string](numReps))
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
	x := New[string](WithNodeLocatorNumReps[string](numReps))
	x.AddNodes("abcdefg")

	if len(x.nodeByKey) != numReps {
		t.Errorf("got %d, want %d", len(x.nodeByKey), numReps)
	}
	if len(x.sortedKeys) != numReps {
		t.Errorf("got %d, want %d", len(x.sortedKeys), numReps)
	}
	if !slices.IsSorted(x.sortedKeys) {
		t.Errorf("expected sorted hashes to be sorted")
	}
	x.AddNodes("qwer")

	if len(x.nodeByKey) != 2*numReps {
		t.Errorf("got %d, want %d", len(x.nodeByKey), 2*numReps)
	}
	if len(x.sortedKeys) != 2*numReps {
		t.Errorf("got %d, want %d", len(x.nodeByKey), 2*numReps)
	}
	if !slices.IsSorted(x.sortedKeys) {
		t.Errorf("expected sorted hashes to be sorted")
	}
}

func TestRemove(t *testing.T) {
	numReps := 160
	x := New[string](WithNodeLocatorNumReps[string](numReps))
	x.AddNodes("abcdefg")
	x.RemoveNodes("abcdefg")
	if len(x.nodeByKey) != 0 {
		t.Errorf("got %d, want %d", len(x.nodeByKey), 0)
	}
	if len(x.sortedKeys) != 0 {
		t.Errorf("got %d, want %d", len(x.nodeByKey), 0)
	}
}

func TestRemoveNonExisting(t *testing.T) {
	numReps := 160
	x := New[string](WithNodeLocatorNumReps[string](numReps))
	x.AddNodes("abcdefg")
	x.RemoveNodes("abcdefghijk")
	if len(x.nodeByKey) != numReps {
		t.Errorf("got %d, want %d", len(x.nodeByKey), numReps)
	}
}

func TestGetEmpty(t *testing.T) {
	numReps := 160
	x := New[string](WithNodeLocatorNumReps[string](numReps))
	_, has := x.Get("asdfsadfsadf")
	if has {
		t.Errorf("expected error")
	}
}

func TestGetSingle(t *testing.T) {
	numReps := 160
	x := New[string](WithNodeLocatorNumReps[string](numReps))
	x.AddNodes("abcdefg")
	f := func(s string) bool {
		y, has := x.Get(s)
		if !has {
			return false
		}
		// t.Logf("s = %q, y = %q", s, y)
		return y == "abcdefg"
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
	x := New[string]()
	x.AddNodes("abcdefg")
	x.AddNodes("hijklmn")
	x.AddNodes("opqrstu")
	for i, v := range gmtests {
		result, has := x.Get(v.in)
		if !has {
			t.Fatal()
		}
		if result != v.out {
			t.Errorf("%d. got %q, expected %q", i, result, v.out)
		}
	}
}

func TestGetMultipleQuick(t *testing.T) {
	x := New[string]()
	x.AddNodes("abcdefg")
	x.AddNodes("hijklmn")
	x.AddNodes("opqrstu")
	f := func(s string) bool {
		y, has := x.Get(s)
		if !has {
			return false
		}
		// t.Logf("s = %q, y = %q", s, y)
		return y == "abcdefg" ||
			y == "hijklmn" ||
			y == "opqrstu"
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
	x := New[string]()
	x.AddNodes("abcdefg")
	x.AddNodes("hijklmn")
	x.AddNodes("opqrstu")
	for i, v := range rtestsBefore {
		result, has := x.Get(v.in)
		if !has {
			t.Fatal()
		}
		if result != v.out {
			t.Errorf("%d. got %q, expected %q before rm", i, result, v.out)
		}
	}
	x.RemoveNodes("hijklmn")
	for i, v := range rtestsAfter {
		result, has := x.Get(v.in)
		if !has {
			t.Fatal()
		}
		if result != v.out {
			t.Errorf("%d. got %q, expected %q after rm", i, result, v.out)
		}
	}
}

func TestGetMultipleRemoveQuick(t *testing.T) {
	x := New[string]()
	x.AddNodes("abcdefg")
	x.AddNodes("hijklmn")
	x.AddNodes("opqrstu")
	x.RemoveNodes("opqrstu")
	f := func(s string) bool {
		y, has := x.Get(s)
		if !has {
			t.Logf("missing node")
			return false
		}
		//t.Logf("s = %q, y = %q", s, y)
		return y == "abcdefg" || y == "hijklmn"
	}
	if err := quick.Check(f, nil); err != nil {
		t.Logf("missing nodes")
	}
}

func TestGetTwo(t *testing.T) {
	x := New[string]()
	x.AddNodes("abcdefg")
	x.AddNodes("hijklmn")
	x.AddNodes("opqrstu")
	nodes, has := x.GetN("99999999", 2)
	if !has {
		t.Fatal("missing nodes")
	}
	if len(nodes) != 2 {
		t.Errorf("expected 2 allNodes instead of %d", len(nodes))
	}
	a, b := nodes[0], nodes[1]
	if a == b {
		t.Errorf("a shouldn't equal b")
	}
	if a != "opqrstu" {
		t.Errorf("wrong a: %q", a)
	}
	if b != "hijklmn" {
		t.Errorf("wrong b: %q", b)
	}
}

func TestGetTwoQuick(t *testing.T) {
	x := New[string]()
	x.AddNodes("abcdefg")
	x.AddNodes("hijklmn")
	x.AddNodes("opqrstu")
	f := func(s string) bool {
		nodes, has := x.GetN(s, 2)
		if !has {
			t.Logf("missing nodes")
			return false
		}
		if len(nodes) != 2 {
			t.Logf("expected 2 allNodes instead of %d", len(nodes))
			return false
		}
		a, b := nodes[0], nodes[1]
		if a == b {
			t.Logf("a == b")
			return false
		}
		if a != "abcdefg" &&
			a != "hijklmn" &&
			a != "opqrstu" {
			t.Logf("invalid a: %q", a)
			return false
		}

		if b != "abcdefg" &&
			b != "hijklmn" &&
			b != "opqrstu" {
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
	x := New[string]()
	x.AddNodes("abcdefg")
	x.AddNodes("hijklmn")
	f := func(s string) bool {
		nodes, has := x.GetN(s, 2)
		if !has {
			t.Logf("missing nodes")
			return false
		}
		if len(nodes) != 2 {
			t.Logf("expected 2 allNodes instead of %d", len(nodes))
			return false
		}
		a, b := nodes[0], nodes[1]
		if a == b {
			t.Logf("a == b")
			return false
		}
		if a != "abcdefg" && a != "hijklmn" {
			t.Logf("invalid a: %q", a)
			return false
		}

		if b != "abcdefg" && b != "hijklmn" {
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
	x := New[string]()

	x.AddNodes("abcdefg")
	nodes, has := x.GetN("99999999", 2)
	if !has {
		t.Fatalf("missing nodes")
	}
	if len(nodes) != 1 {
		t.Fatalf("expected 2 allNodes instead of %d", len(nodes))
	}
	a, b := nodes[0], ""
	if a == b {
		t.Errorf("a shouldn't equal b")
	}
	if a != "abcdefg" {
		t.Errorf("wrong a: %q", a)
	}
}

func TestGetN(t *testing.T) {
	x := New[string]()
	x.AddNodes("abcdefg")
	x.AddNodes("hijklmn")
	x.AddNodes("opqrstu")
	members, has := x.GetN("9999999", 3)
	if !has {
		t.Logf("missing nodes")
	}
	if len(members) != 3 {
		t.Errorf("expected 3 allNodes instead of %d", len(members))
	}
	if members[0] != "abcdefg" {
		t.Errorf("wrong allNodes[0]: %q", members[0])
	}
	if members[1] != "opqrstu" {
		t.Errorf("wrong allNodes[1]: %q", members[1])
	}
	if members[2] != "hijklmn" {
		t.Errorf("wrong allNodes[2]: %q", members[2])
	}
}

func TestGetNLess(t *testing.T) {
	x := New[string]()
	x.AddNodes("abcdefg")
	x.AddNodes("hijklmn")
	x.AddNodes("opqrstu")
	members, has := x.GetN("99999999", 2)
	if !has {
		t.Logf("missing nodes")
	}
	if len(members) != 2 {
		t.Errorf("expected 2 allNodes instead of %d", len(members))
	}
	if members[0] != "opqrstu" {
		t.Errorf("wrong allNodes[0]: %q", members[0])
	}
	if members[1] != "hijklmn" {
		t.Errorf("wrong allNodes[1]: %q", members[1])
	}
}

func TestGetNMore(t *testing.T) {
	x := New[string]()
	x.AddNodes("abcdefg")
	x.AddNodes("hijklmn")
	x.AddNodes("opqrstu")
	members, has := x.GetN("9999999", 5)
	if !has {
		t.Logf("missing nodes")
	}
	if len(members) != 3 {
		t.Errorf("expected 3 allNodes instead of %d", len(members))
	}
	if members[0] != "abcdefg" {
		t.Errorf("wrong allNodes[0]: %q", members[0])
	}
	if members[1] != "opqrstu" {
		t.Errorf("wrong allNodes[1]: %q", members[1])
	}
	if members[2] != "hijklmn" {
		t.Errorf("wrong allNodes[2]: %q", members[2])
	}
}

func TestGetNQuick(t *testing.T) {
	x := New[string]()
	x.AddNodes("abcdefg")
	x.AddNodes("hijklmn")
	x.AddNodes("opqrstu")
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
			if set[member] {
				t.Logf("duplicate error")
				return false
			}
			set[member] = true
			if member != "abcdefg" &&
				member != "hijklmn" &&
				member != "opqrstu" {
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
	x := New[string]()
	x.AddNodes("abcdefg")
	x.AddNodes("hijklmn")
	x.AddNodes("opqrstu")
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
			if set[member] {
				t.Logf("duplicate error")
				return false
			}
			set[member] = true
			if member != "abcdefg" &&
				member != "hijklmn" &&
				member != "opqrstu" {
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
	x := New[string]()
	x.AddNodes("abcdefg")
	x.AddNodes("hijklmn")
	x.AddNodes("opqrstu")
	f := func(s string) bool {
		// t.Log("check", s)
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
			if set[member] {
				t.Logf("duplicate error")
				return false
			}
			set[member] = true
			if member != "abcdefg" && member != "hijklmn" && member != "opqrstu" {
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
	x := New[string]()
	x.AddNodes("abc")
	x.AddNodes("def")
	x.AddNodes("ghi")
	x.SetNodes("jkl", "mno")
	if len(x.allNodes) != 2 {
		t.Errorf("expected 2 elts, got %d", len(x.allNodes))
	}
	nodes, has := x.GetN("qwerqwerwqer", 2)
	if !has {
		t.Fatal()
	}
	if len(nodes) != 2 {
		t.Errorf("expected 2 allNodes instead of %d", len(nodes))
	}
	a, b := nodes[0], nodes[1]
	if a != "jkl" && a != "mno" {
		t.Errorf("expected jkl or mno, got %s", a)
	}
	if b != "jkl" && b != "mno" {
		t.Errorf("expected jkl or mno, got %s", b)
	}
	if a == b {
		t.Errorf("expected a != b, they were both %s", a)
	}
	x.SetNodes("jkl", "mno")
	if len(x.allNodes) != 2 {
		t.Errorf("expected 2 elts, got %d", len(x.allNodes))
	}
	nodes, has = x.GetN("qwerqwerwqer", 2)
	if !has {
		t.Fatal()
	}
	if len(nodes) != 2 {
		t.Errorf("expected 2 allNodes instead of %d", len(nodes))
	}
	a, b = nodes[0], nodes[1]
	if a != "jkl" && a != "mno" {
		t.Errorf("expected jkl or mno, got %s", a)
	}
	if b != "jkl" && b != "mno" {
		t.Errorf("expected jkl or mno, got %s", b)
	}
	if a == b {
		t.Errorf("expected a != b, they were both %s", a)
	}
	x.SetNodes("pqr", "mno")
	if len(x.allNodes) != 2 {
		t.Errorf("expected 2 elts, got %d", len(x.allNodes))
	}
	nodes, has = x.GetN("qwerqwerwqer", 2)
	if !has {
		t.Fatal()
	}
	if len(nodes) != 2 {
		t.Errorf("expected 2 allNodes instead of %d", len(nodes))
	}
	a, b = nodes[0], nodes[1]
	if a != "pqr" && a != "mno" {
		t.Errorf("expected jkl or mno, got %s", a)
	}
	if b != "pqr" && b != "mno" {
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
	x := New[string]()
	x.AddNodes("stays")
	b.ResetTimer()
	allocSize := allocBytes(func() {
		for i := 0; i < b.N; i++ {
			x.AddNodes("Foo")
			x.RemoveNodes("Foo")
		}
	})
	b.Logf("%d: Allocated %d bytes (%.2fx)", b.N, allocSize, float64(allocSize)/float64(b.N))
}

func BenchmarkMalloc(b *testing.B) {
	x := New[string]()
	x.AddNodes("stays")
	b.ResetTimer()
	mallocs := mallocNum(func() {
		for i := 0; i < b.N; i++ {
			x.AddNodes("Foo")
			x.RemoveNodes("Foo")
		}
	})
	b.Logf("%d: Mallocd %d times (%.2fx)", b.N, mallocs, float64(mallocs)/float64(b.N))
}

func BenchmarkCycle(b *testing.B) {
	x := New[string]()
	x.AddNodes("nothing")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x.AddNodes("foo" + strconv.Itoa(i))
		x.RemoveNodes("foo" + strconv.Itoa(i))
	}
}

func BenchmarkCycleLarge(b *testing.B) {
	x := New[string]()
	for i := 0; i < 10; i++ {
		x.AddNodes("start" + strconv.Itoa(i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x.AddNodes("foo" + strconv.Itoa(i))
		x.RemoveNodes("foo" + strconv.Itoa(i))
	}
}

func BenchmarkGet(b *testing.B) {
	x := New[string]()
	x.AddNodes("nothing")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x.Get("nothing")
	}
}

func BenchmarkGetLarge(b *testing.B) {
	x := New[string]()
	for i := 0; i < 10; i++ {
		x.AddNodes("start" + strconv.Itoa(i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x.Get("nothing")
	}
}

func BenchmarkGetN(b *testing.B) {
	x := New[string]()
	x.AddNodes("nothing")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x.GetN("nothing", 3)
	}
}

func BenchmarkGetNLarge(b *testing.B) {
	x := New[string]()
	for i := 0; i < 10; i++ {
		x.AddNodes("start" + strconv.Itoa(i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x.GetN("nothing", 3)
	}
}

func BenchmarkGetTwo(b *testing.B) {
	x := New[string]()
	x.AddNodes("nothing")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x.GetN("nothing", 2)
	}
}

func BenchmarkGetTwoLarge(b *testing.B) {
	x := New[string]()
	for i := 0; i < 10; i++ {
		x.AddNodes("start" + strconv.Itoa(i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x.GetN("nothing", 2)
	}
}

// from @edsrzf on github:
func TestAddCollision(t *testing.T) {
	// These two strings produce several crc32 collisions after "|i" is
	// appended added by NodeLocator.virtualNode.
	const s1 = "abear"
	const s2 = "solidiform"
	x := New[string]()
	x.AddNodes(s1)
	x.AddNodes(s2)
	elt1, has := x.Get("abear")
	if !has {
		t.Fatal("missing node")
	}

	y := New[string]()
	// add elements in opposite order
	y.AddNodes(s2)
	y.AddNodes(s1)
	elt2, has := y.Get(s1)
	if !has {
		t.Fatal("missing node")
	}

	if elt1 != elt2 {
		t.Error(elt1, "and", elt2, "should be equal")
	}
}
