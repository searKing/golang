// Copyright 2020 The searKing Author. All rights reserved.
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

// TernarySearchTree represents a Ternary Search Tree.
// The zero value for List is an empty list ready to use.
type TernarySearchTree struct {
	root node // sentinel list node, only &root, root.prev, and root.next are used
}

// Init initializes or clears tree l.
func (l *TernarySearchTree) init() *TernarySearchTree {
	l.root.left = &l.root
	l.root.middle = &l.root
	l.root.right = &l.root
	l.root.tree = l
	return l
}

// New creates a tenary search tree, which is a trie tree.
func New(prefixes ...string) *TernarySearchTree {
	tree := (&TernarySearchTree{}).init()
	for _, prefix := range prefixes {
		tree.Store(prefix, nil)
	}
	return tree
}

// NewWithBytes behaves like New, but receive ...[]byte
func NewWithBytes(prefixes ...[]byte) *TernarySearchTree {
	tree := (&TernarySearchTree{}).init()
	for _, prefix := range prefixes {
		tree.Store(string(prefix), nil)
	}
	return tree
}

// Depth return max len of all prefixs
func (l *TernarySearchTree) Depth() int {
	return l.root.Depth()
}

// Count returns the number of elements of list l, excluding (this) sentinel node.
func (l *TernarySearchTree) Count() int {
	var len int
	l.Traversal(traversal.Preorder, HandlerFunc(func(prefix []byte, value interface{}) (goon bool) {
		len++
		return true
	}))
	return len
}

// Front returns the first node of list l or nil if the list is empty.
func (l *TernarySearchTree) Left() *node {
	return l.root.left
}

// Middle returns the first node of list l or nil if the list is empty.
func (l *TernarySearchTree) Middle() *node {
	return l.root.middle
}

// Right returns the first node of list l or nil if the list is empty.
func (l *TernarySearchTree) Right() *node {
	return l.root.right
}

// lazyInit lazily initializes a zero List value.
func (l *TernarySearchTree) lazyInit() {
	if l.root.right == nil {
		l.init()
	}
}

func (l *TernarySearchTree) Traversal(order traversal.Order, handler Handler) {
	l.root.Traversal(order, handler)
}

// Follow returns node info of longest subPrefix of prefix
func (l *TernarySearchTree) Follow(prefix string) (key string, value interface{}, ok bool) {
	pre, val, ok := l.root.Follow([]byte(prefix))
	return string(pre), val, ok
}

// Load loads value by key
func (l *TernarySearchTree) Load(key string) (value interface{}, ok bool) {
	return l.root.Load([]byte(key))
}

// Store stores value by key
func (l *TernarySearchTree) Store(key string, value interface{}) {
	l.root.Store([]byte(key), value)
}

// return true if the node matched with key, with value nonempty
func (l *TernarySearchTree) Contains(key string) bool {
	return l.root.Contains([]byte(key))
}

// return true if any node exists with key prefix, no matter value is empty or not
func (l *TernarySearchTree) ContainsPrefix(prefix string) bool {
	return l.root.ContainsPrefix([]byte(prefix))
}

// remove the node matched with key
func (l *TernarySearchTree) Remove(key string, shrinkToFit bool) (old interface{}, ok bool) {
	return l.root.Remove([]byte(key), shrinkToFit)
}

// remove all nodes started with key prefix
func (l *TernarySearchTree) RemoveAll(prefix string) (old interface{}, ok bool) {
	return l.root.RemoveAll([]byte(prefix))
}

func (l *TernarySearchTree) String() string {
	return l.root.String()
}
