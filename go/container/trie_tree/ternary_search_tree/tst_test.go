package ternary_search_tree_test

import (
	"github.com/searKing/golang/go/container/trie_tree/ternary_search_tree"
	"testing"
)

func TestTernarySearchTree(t *testing.T) {
	tree := ternary_search_tree.New()
	tree.Insert("test", 1)
	if tree.Len() != 1 {
		t.Errorf("expecting len 1")
	}
	if !tree.Contains("test") {
		t.Errorf("expecting to find key=test")
	}

	tree.Insert("testing", 2)
	tree.Insert("abcd", 0)

	found := false
	tree.TraversalInOrderFunc(func(key string, val interface{}) bool {
		if key == "test" && val.(int) == 1 {
			found = true
		}
		return true
	})
	if !found {
		t.Errorf("expecting iterator to find test")
	}
	tree.Remove("testing")
	tree.Remove("abcd")

	v, ok := tree.Remove("test")
	if !ok {
		t.Errorf("expecting test can be found to be removed")
	}
	if tree.Len() != 0 {
		t.Errorf("expecting len 0")
	}
	if tree.Contains("test") {
		t.Errorf("expecting not to find key=test")
	}
	if v.(int) != 1 {
		t.Errorf("expecting value=1")
	}
}

func TestTernarySearchTree_String1(t *testing.T) {

	tree := ternary_search_tree.New()
	tree.Insert("abcd", 0)
	tree.Insert("abcd1234ABCD", 2)
	tree.Insert("abcd1234", 1)
	s := tree.String()
	expect := `a:<nil>
ab:<nil>
abc:<nil>
abcd:0
abcd1:<nil>
abcd12:<nil>
abcd123:<nil>
abcd1234:1
abcd1234A:<nil>
abcd1234AB:<nil>
abcd1234ABC:<nil>
abcd1234ABCD:2
`
	if s != expect {
		t.Errorf("expect %s", expect)
	}
}

func TestTernarySearchTree_String2(t *testing.T) {

	tree := ternary_search_tree.New()
	tree.Insert("abcd", 0)
	tree.Insert("1234", 1)
	s := tree.String()
	expect := `1:<nil>
12:<nil>
123:<nil>
1234:1
a:<nil>
ab:<nil>
abc:<nil>
abcd:0
`
	if s != expect {
		t.Errorf("expect %s", expect)
	}
}
