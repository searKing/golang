// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ternary_search_tree_test

import (
	"fmt"

	"github.com/searKing/golang/go/container/traversal"
	"github.com/searKing/golang/go/container/trie_tree/ternary_search_tree"
)

func ExampleNew() {
	tree := ternary_search_tree.New("abcdef")
	fmt.Println("count:	", tree.Count())
	fmt.Println("depth:	", tree.Depth())
	fmt.Printf("contains key %q:	%v\n", "abcdef", tree.Contains("abcdef"))
	fmt.Printf("contains key prefix %q:	%v\n", "ab", tree.ContainsPrefix("ab"))
	val, _ := tree.Load("test")
	fmt.Printf("load key %q's value:	%v\n", "test", val)

	// Output:
	// count:	 1
	// depth:	 6
	// contains key "abcdef":	true
	// contains key prefix "ab":	true
	// load key "test"'s value:	<nil>
}

func ExampleTernarySearchTree_Store() {
	tree := ternary_search_tree.New()
	tree.Store("key", "val")

	val, ok := tree.Load("key")
	fmt.Printf("%q:	%q, %v\n", "key", val, ok)
	// Output:
	// "key":	"val", true
}

func ExampleTernarySearchTree_Load() {
	tree := ternary_search_tree.New()
	tree.Store("key", "val")

	val, ok := tree.Load("key")
	fmt.Printf("%q:	%q, %v\n", "key", val, ok)
	val, ok = tree.Load("not exist key")
	fmt.Printf("%q:	%v, %v\n", "not exist key", val, ok)
	// Output:
	// "key":	"val", true
	// "not exist key":	<nil>, false
}

func ExampleTernarySearchTree_Count() {
	tree := ternary_search_tree.New()
	tree.Store("a", nil)
	fmt.Println(tree.Count())
	tree.Store("ab", nil)
	fmt.Println(tree.Count())
	tree.Store("x", nil)
	fmt.Println(tree.Count())
	tree.Store("y", nil)
	fmt.Println(tree.Count())

	tree.Remove("abc", true)
	fmt.Println(tree.Count())
	tree.Remove("x", true)
	fmt.Println(tree.Count())
	tree.Remove("a", true)
	fmt.Println(tree.Count())

	// Output:
	// 1
	// 2
	// 3
	// 4
	// 4
	// 3
	// 2
}

func ExampleNode_Depth() {
	tree := ternary_search_tree.New()
	tree.Store("a", nil)
	fmt.Println(tree.Depth())
	tree.Store("ab", nil)
	fmt.Println(tree.Depth())
	tree.Store("x", nil)
	fmt.Println(tree.Depth())
	tree.Store("y", nil)
	fmt.Println(tree.Depth())

	tree.Remove("abc", true)
	fmt.Println(tree.Depth())
	tree.Remove("x", true)
	fmt.Println(tree.Depth())
	tree.Remove("a", true)
	fmt.Println(tree.Depth())
	tree.Remove("ab", true)
	fmt.Println(tree.Depth())

	// Output:
	// 1
	// 2
	// 2
	// 2
	// 2
	// 2
	// 2
	// 1
}

func ExampleTernarySearchTree_RemoveAll() {
	tree := ternary_search_tree.New()
	tree.Store("a", nil)
	fmt.Println(tree.Count())
	tree.Store("ab", nil)
	fmt.Println(tree.Count())
	tree.Store("x", nil)
	fmt.Println(tree.Count())
	tree.Store("y", nil)
	fmt.Println(tree.Count())

	tree.RemoveAll("x")
	fmt.Println(tree.Count())
	tree.RemoveAll("a")
	fmt.Println(tree.Count())
	tree.RemoveAll("ab")
	fmt.Println(tree.Count())

	// Output:
	// 1
	// 2
	// 3
	// 4
	// 2
	// 1
	// 1
}

func ExampleTernarySearchTree_Traversal() {
	tree := ternary_search_tree.New()
	tree.Store("test", 1)
	tree.Store("test", 11)
	tree.Store("testing", 2)
	tree.Store("abcd", 0)

	found := false
	tree.Traversal(traversal.Preorder, ternary_search_tree.HandlerFunc(
		func(key []byte, val interface{}) bool {
			if string(key) == "test" && val.(int) == 11 {
				found = true
				return false
			}
			return true
		}))
	fmt.Printf("traversal for key %q, found:	%v\n", "test", found)

	// Output:
	// traversal for key "test", found:	true
}
