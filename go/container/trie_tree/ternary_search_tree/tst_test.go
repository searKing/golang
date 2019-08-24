package ternary_search_tree_test

import (
	"github.com/searKing/golang/go/container/traversal"
	"github.com/searKing/golang/go/container/trie_tree/ternary_search_tree"
	"testing"
)

func TestTernarySearchTree(t *testing.T) {
	tree := ternary_search_tree.New()
	tree.Store("test", 1)
	if tree.Len() != 1 {
		t.Errorf("expecting len 1, actual = %v", tree.Len())
	}
	if !tree.Contains("test") {
		t.Errorf("expecting to find key=test")
	}
	if !tree.ContainsPrefix("tes") {
		t.Errorf("expecting to find key=tes")
	}

	val, ok := tree.Load("test")
	if !ok {
		t.Errorf("expecting to find key=test")
	}
	if val.(int) != 1 {
		t.Errorf("expecting test's value=1")
	}
	tree.Store("test", 11)
	val, ok = tree.Load("test")
	if !ok {
		t.Errorf("expecting to find key=test")
	}
	if val.(int) != 11 {
		t.Errorf("expecting test's value=11, actual = %v", val)
	}

	tree.Store("testing", 2)
	tree.Store("abcd", 0)
	if tree.Len() != 3 {
		t.Errorf("expecting len 3, actual = %v", tree.Len())
	}

	found := false
	tree.Traversal(traversal.Preorder, ternary_search_tree.HandlerFunc(
		func(key []byte, val interface{}) bool {
			if string(key) == "test" && val.(int) == 11 {
				found = true
				return false
			}
			return true
		}))
	if !found {
		t.Errorf("expecting iterator to find test")
	}

	val, ok = tree.Load("testing")
	if !ok {
		t.Errorf("expecting to find key=testing")
	}
	if val.(int) != 2 {
		t.Errorf("expecting testing's value=2")
	}

	val, ok = tree.Load("abcd")
	if !ok {
		t.Errorf("expecting to find key=abcd")
	}
	if val.(int) != 0 {
		t.Errorf("expecting abcd's value=0")
	}

	tree.Remove("testing", true)
	tree.Remove("abcd", false)

	v, ok := tree.Remove("test",false)
	if !ok {
		t.Errorf("expecting test can be found to be removed")
	}

	if tree.Len() != 0 {
		t.Errorf("expecting len 3, actual = %v", tree.Len())
	}

	if tree.Contains("test") {
		t.Errorf("expecting not to find key=test")
	}
	if v.(int) != 11 {
		t.Errorf("expecting test's value=11, actual = %v", val)
	}
}

func TestTernarySearchTree_String1(t *testing.T) {

	tree := ternary_search_tree.New()
	tree.Store("abcd", 0)
	tree.Store("abcd1234ABCD", 2)
	tree.Store("abcd1234", 1)
	s := tree.String()
	expect := `abcd:0
abcd1234:1
abcd1234ABCD:2
`
	if s != expect {
		t.Errorf("actual:\n%s\nexpect:\n%s", s, expect)
	}
}

func TestTernarySearchTree_String2(t *testing.T) {

	tree := ternary_search_tree.New()
	tree.Store("abcd", 0)
	tree.Store("1234", 1)
	s := tree.String()
	expect := `1234:1
abcd:0
`
	if s != expect {
		t.Errorf("actual:\n%s\nexpect:\n%s", s, expect)
	}
}
