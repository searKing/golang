// Copyright 2019 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// https://en.wikipedia.org/wiki/Ternary_search_tree
// In computer science, a ternary search tree is a type of trie (sometimes called a prefix tree)
// where nodes are arranged in a manner similar to a binary search tree,
// but with up to three children rather than the binary tree's limit of two.
// Like other prefix trees, a ternary search tree can be used as an associative map structure
// with the ability for incremental string search.
// However, ternary search trees are more space efficient compared to standard prefix trees,
// at the cost of speed. Common applications for ternary search trees include spell-checking and
// auto-completion.
package ternary_search_tree

// TernarySearchTree represents a Ternary Search Tree.
// The zero value for List is an empty list ready to use.
type TernarySearchTree struct {
	root Element // sentinel list element, only &root, root.prev, and root.next are used
	len  int     // current list length excluding (this) sentinel element
}

// Init initializes or clears tree l.
func (l *TernarySearchTree) Init() *TernarySearchTree {
	l.root.left = &l.root
	l.root.middle = &l.root
	l.root.right = &l.root
	l.root.tree = l
	l.len = 0
	return l
}

// Init initializes or clears Tree l.
func New() *TernarySearchTree {
	return (&TernarySearchTree{}).Init()
}

// Len returns the number of elements of list l.
// The complexity is O(1).
func (l *TernarySearchTree) Len() int { return l.len }

// Front returns the first element of list l or nil if the list is empty.
func (l *TernarySearchTree) Left() *Element {
	if l.len == 0 {
		return nil
	}
	return l.root.left
}

// Middle returns the first element of list l or nil if the list is empty.
func (l *TernarySearchTree) Middle() *Element {
	if l.len == 0 {
		return nil
	}
	return l.root.middle
}

// Right returns the first element of list l or nil if the list is empty.
func (l *TernarySearchTree) Right() *Element {
	if l.len == 0 {
		return nil
	}
	return l.root.right
}

// lazyInit lazily initializes a zero List value.
func (l *TernarySearchTree) lazyInit() {
	if l.root.right == nil {
		l.Init()
	}
}

func (l *TernarySearchTree) TraversalPreOrderFunc(f func(prefix string, value interface{}) (goon bool)) (goon bool) {
	return l.root.TraversalPreOrderFunc(func(pre []byte, v interface{}) (goon bool) {
		return f(string(pre), v)
	})
}
func (l *TernarySearchTree) TraversalInOrderFunc(f func(prefix string, value interface{}) (goon bool)) (goon bool) {
	return l.root.TraversalInOrderFunc(func(pre []byte, v interface{}) (goon bool) {
		return f(string(pre), v)
	})
}
func (l *TernarySearchTree) TraversalPostOrderFunc(f func(prefix string, value interface{}) (goon bool)) (goon bool) {
	return l.root.TraversalPostOrderFunc(func(pre []byte, v interface{}) (goon bool) {
		return f(string(pre), v)
	})
}
func (l *TernarySearchTree) Get(prefix string) (value interface{}, ok bool) {
	return l.root.Get([]byte(prefix))
}
func (l *TernarySearchTree) Contains(prefix string) bool {
	return l.root.Contains([]byte(prefix))
}

func (l *TernarySearchTree) Insert(prefix string, value interface{}) {
	l.root.Insert([]byte(prefix), value)
	l.len++
}

func (l *TernarySearchTree) Remove(prefix string) (value interface{}, ok bool) {
	value, ok = l.root.Remove([]byte(prefix))
	if ok {
		l.len--
	}
	return value, ok
}
func (l *TernarySearchTree) String() string {
	return l.root.String()
}
