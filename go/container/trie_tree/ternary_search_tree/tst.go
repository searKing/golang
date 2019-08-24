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

import (
	"github.com/searKing/golang/go/container/traversal"
)

type TernarySearchTree interface {
	// Depth return max len of all prefixs
	Depth() int
	// Count returns the number of elements of list l, excluding (this) sentinel node.
	Count() int
	// Store stores value in prefix
	Store(prefix string, value interface{})
	// Follow returns node info of longest subPrefix of prefix
	Follow(prefix string) (subPrefix string, value interface{}, ok bool)
	Load(prefix string) (value interface{}, ok bool)
	// return true if prefix with key and value
	Contains(prefix string) bool
	// return true if prefix with key
	ContainsPrefix(prefix string) bool
	// remove value with prefix
	Remove(prefix string, shrinkToFit bool) (value interface{}, old bool)
	// remove node with prefix
	RemoveAll(prefix string) (value interface{}, ok bool)
	String() string
	Traversal(order traversal.Order, handler Handler)
}

// ternarySearchTree represents a Ternary Search Tree.
// The zero value for List is an empty list ready to use.
type ternarySearchTree struct {
	root node // sentinel list node, only &root, root.prev, and root.next are used
}

// Init initializes or clears tree l.
func (l *ternarySearchTree) init() *ternarySearchTree {
	l.root.left = &l.root
	l.root.middle = &l.root
	l.root.right = &l.root
	l.root.tree = l
	return l
}

// Init initializes or clears Tree l.
func New(prefixes ...string) TernarySearchTree {
	tree := (&ternarySearchTree{}).init()
	for _, prefix := range prefixes {
		tree.Store(prefix, nil)
	}
	return tree
}

func NewWithBytes(prefixes ...[]byte) TernarySearchTree {
	tree := (&ternarySearchTree{}).init()
	for _, prefix := range prefixes {
		tree.Store(string(prefix), nil)
	}
	return tree
}

// Depth return max len of all prefixs
func (l *ternarySearchTree) Depth() int {
	return l.root.Depth()
}

// Count returns the number of elements of list l, excluding (this) sentinel node.
func (l *ternarySearchTree) Count() int {
	var len int
	l.Traversal(traversal.Preorder, HandlerFunc(func(prefix []byte, value interface{}) (goon bool) {
		len++
		return true
	}))
	return len
}

// Front returns the first node of list l or nil if the list is empty.
func (l *ternarySearchTree) Left() *node {
	return l.root.left
}

// Middle returns the first node of list l or nil if the list is empty.
func (l *ternarySearchTree) Middle() *node {
	return l.root.middle
}

// Right returns the first node of list l or nil if the list is empty.
func (l *ternarySearchTree) Right() *node {
	return l.root.right
}

// lazyInit lazily initializes a zero List value.
func (l *ternarySearchTree) lazyInit() {
	if l.root.right == nil {
		l.init()
	}
}

func (l *ternarySearchTree) Traversal(order traversal.Order, handler Handler) {
	l.root.Traversal(order, handler)
}

func (l *ternarySearchTree) Follow(prefix string) (subPrefix string, value interface{}, ok bool) {
	pre, val, ok := l.root.Follow([]byte(prefix))
	return string(pre), val, ok
}

func (l *ternarySearchTree) Load(prefix string) (value interface{}, ok bool) {
	return l.root.Load([]byte(prefix))
}

func (l *ternarySearchTree) Store(prefix string, value interface{}) {
	l.root.Store([]byte(prefix), value)
}

func (l *ternarySearchTree) Contains(prefix string) bool {
	return l.root.Contains([]byte(prefix))
}

func (l *ternarySearchTree) ContainsPrefix(prefix string) bool {
	return l.root.ContainsPrefix([]byte(prefix))
}

func (l *ternarySearchTree) Remove(prefix string, shrinkToFit bool) (old interface{}, ok bool) {
	return l.root.Remove([]byte(prefix), shrinkToFit)
}

func (l *ternarySearchTree) RemoveAll(prefix string) (old interface{}, ok bool) {
	return l.root.RemoveAll([]byte(prefix))
}

func (l *ternarySearchTree) String() string {
	return l.root.String()
}
