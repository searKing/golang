// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ternary_search_tree

import (
	"fmt"
	"strings"

	"github.com/searKing/golang/go/container/traversal"
)

const (
	NilKey = 0
)

type node struct {
	prefix   []byte
	key      byte
	hasKey   bool
	value    any
	hasValue bool

	left, middle, right *node
	tree                *TernarySearchTree
}

func (n *node) LeftNodes() []any {
	left := n.Left()
	if left == nil {
		return nil
	}
	return []any{left}
}

func (n *node) MiddleNodes() []any {
	middle := n.Middle()
	if middle == nil {
		return nil
	}
	return []any{middle}
}

func (n *node) RightNodes() []any {
	right := n.Right()
	if right == nil {
		return nil
	}
	return []any{right}
}

// Left returns the left list node or nil.
func (n *node) Left() *node {
	if p := n.left; n.tree != nil && p != &n.tree.root {
		return p
	}
	return nil
}

// Middle returns the middle list node or nil.
func (n *node) Middle() *node {
	if p := n.middle; n.tree != nil && p != &n.tree.root {
		return p
	}
	return nil
}

// Right returns the right list node or nil.
func (n *node) Right() *node {
	if p := n.right; n.tree != nil && p != &n.tree.root {
		return p
	}
	return nil
}

func (n *node) Traversal(order traversal.Order, handler Handler) {
	if handler == nil {
		return
	}
	if !n.hasKey {
		return
	}
	order(n, traversal.HandlerFunc(func(ele any, depth int) (goon bool) {
		currentNode := ele.(*node)
		if !currentNode.hasKey || !currentNode.hasValue {
			return true
		}
		return handler.Handle(currentNode.prefix, currentNode.value)
	}))
	return
}

func (n *node) Follow(prefix []byte) (key []byte, value any, ok bool) {
	graph := n.follow(prefix)
	if len(graph) == 0 {
		return nil, nil, false
	}
	tail := graph[len(graph)-1]
	return tail.prefix, tail.value, tail.hasValue
}

func (n *node) Load(prefix []byte) (value any, ok bool) {
	cur, _, _, has := n.search(prefix)
	if !has {
		return nil, false
	}
	return cur.value, cur.hasValue
}

func (n *node) ContainsPrefix(prefix []byte) bool {
	cur, _, _, has := n.search(prefix)
	if !has {
		return false
	}
	return cur.hasKey
}

func (n *node) Contains(prefix []byte) bool {
	_, ok := n.Load(prefix)
	return ok
}

type pos int

const (
	posNotFound pos = iota // Root's middle, default
	posLeft
	posMiddle
	posRight
	posRoot
)

func (n *node) Store(prefix []byte, value any) {
	// force update
	n.CAS(prefix, nil, value, func(x, y any) int { return 0 })
}

func (n *node) remove(prefix []byte, shrinkToFit bool, omitMiddle bool) (old any, ok bool) {
	cur, last, lastPos, has := n.search(prefix)
	if !has {
		return nil, false
	}
	// shrinkToFit if cur's children are empty
	if shrinkToFit {
		cur.shrinkToFit(last, lastPos, omitMiddle)
	}

	if !cur.hasValue {
		return nil, false
	}
	cur.hasValue = false
	// all matched, goto remove the old
	return cur.value, true
}

func (n *node) Remove(prefix []byte, shrinkToFit bool) (old any, ok bool) {
	return n.remove(prefix, shrinkToFit, false)
}

func (n *node) RemoveAll(prefix []byte) (value any, ok bool) {
	return n.remove(prefix, true, true)
}

func (n *node) String() string {
	s := ""
	n.Traversal(traversal.Inorder, HandlerFunc(func(prefix []byte, value any) (goon bool) {
		s += fmt.Sprintf("%s:%v\n", string(prefix), value)
		return true
	}))

	return strings.TrimRight(s, "\n")
}

func (n *node) CAS(prefix []byte, old, new any, cmps ...func(x, y any) int) bool {
	newElement := func(prefix []byte, hasKey bool, key byte, hasValue bool, value any) *node {
		var p = make([]byte, len(prefix))
		copy(p, prefix)
		return &node{
			prefix:   p,
			key:      key,
			hasKey:   hasKey,
			value:    value,
			hasValue: hasValue,
			left:     &n.tree.root,
			middle:   &n.tree.root,
			right:    &n.tree.root,
			tree:     n.tree,
		}
	}

	cur := n
	for idx := 0; idx < len(prefix); {
		// create the idx layer if not exist
		// otherwise, step to the next layer
		k := prefix[idx]
		if !cur.hasKey {
			cur.key = k
			cur.hasKey = true
			cur.prefix = prefix[:idx+1]
		}
		// goto left
		if k < cur.key {
			left := cur.Left()
			if left == nil {
				cur.left = newElement(prefix[:idx+1], true, k, false, nil)
			}
			cur = cur.left
			continue
		}
		// goto right
		if k > cur.key {
			right := cur.Right()
			if right == nil {
				cur.right = newElement(prefix[:idx+1], true, k, false, nil)
			}
			cur = cur.right
			continue
		}
		// key match, goto match next layer
		idx++
		// all matched, goto set the value
		if idx == len(prefix) {
			// no old
			if !cur.hasValue {
				cur.value = new
				cur.hasValue = true
				return true
			}
			var cmp func(x, y any) int
			if len(cmps) > 0 {
				cmp = cmps[0]
			}
			if cmp == nil {
				cmp = func(x, y any) int {
					if x == y {
						return 0
					}
					return -1
				}
			}
			if cmp(cur.value, old) == 0 {
				cur.value = new
				cur.hasValue = true
				return true
			}
			return false
		}
		// partial matched, goto middle on next layer
		middle := cur.Middle()
		if middle == nil {
			cur.middle = newElement(nil, false, NilKey, false, nil)
		}
		cur = cur.middle
	}
	// never reach
	return false
}

// Depth return max len of all prefixs
func (n *node) Depth() int {
	var depth int
	n.Traversal(traversal.Preorder, HandlerFunc(func(prefix []byte, value any) (goon bool) {
		if depth < len(prefix) {
			depth = len(prefix)
		}
		return true
	}))
	return depth
}

// shrinkToFit cutoff last node's children nodes if all children nodes are empty
func (n *node) shrinkToFit(last *node, lastPos pos, omitMiddle bool) {
	var has bool
	n.Traversal(traversal.Preorder, HandlerFunc(func(prefix []byte, value any) (goon bool) {
		has = true
		return false
	}))
	if !has {
		return
	}

	// match
	switch lastPos {
	case posLeft:
		n.shrinkLeft(last, omitMiddle)
	case posMiddle:
		n.shrinkMiddle(last, omitMiddle)
	case posRight:
		n.shrinkRight(last, omitMiddle)
	case posRoot:
		n.shrinkRoot(last, omitMiddle)
	}
}

// return node graph until prefix matched
func (n *node) follow(prefix []byte) (graph []*node) {
	cur := n

	for idx := 0; idx < len(prefix); {
		// return if nilKey has been meet
		if !cur.hasKey {
			return
		}
		// create the idx layer if not exist
		// otherwise, step to the next layer
		k := prefix[idx]
		if k < cur.key {
			left := cur.Left()
			if left == nil {
				return
			}
			cur = left
			continue
		}
		if k > cur.key {
			right := cur.Right()
			if right == nil {
				return
			}
			cur = right
			continue
		}
		if cur.hasValue {
			graph = append(graph, cur)
		}
		// key match, goto match next layer
		idx++
		// all matched, goto remove the value
		if idx == len(prefix) {
			// match
			return
		}
		// partial matched, goto middle on next layer
		middle := cur.Middle()
		if middle == nil {
			return
		}
		cur = middle
	}
	return
}

// return true if prefix matches, no matter value exists
func (n *node) search(prefix []byte) (cur, last *node, lastPos pos, has bool) {
	cur = n
	last = n

	lastPos = posNotFound
	for idx := 0; idx < len(prefix); {
		// return if nilKey has been met
		if !cur.hasKey {
			return
		}
		// create the idx layer if not exist
		// otherwise, step to the next layer
		k := prefix[idx]
		if k < cur.key {
			left := cur.Left()
			if left == nil {
				return
			}
			last = cur
			lastPos = posLeft
			cur = left
			continue
		}
		if k > cur.key {
			right := cur.Right()
			if right == nil {
				return
			}
			last = cur
			lastPos = posRight
			cur = right
			continue
		}
		// key match, goto match next layer
		idx++
		// all matched, goto remove the value
		if idx == len(prefix) {
			// match
			// Special case: prefix is Root node
			if idx == 1 {
				lastPos = posRoot
			}
			return cur, last, lastPos, true
		}
		// partial matched, goto middle on next layer
		middle := cur.Middle()
		if middle == nil {
			return
		}
		last = cur
		lastPos = posMiddle
		cur = middle
	}
	return
}

func (n *node) IsLeaf() bool {
	if !n.hasKey {
		return false
	}
	if n.Left() != nil {
		return false
	}
	if n.Middle() != nil {
		return false
	}
	if n.Right() != nil {
		return false
	}
	return false
}

func (n *node) shrinkLeft(last *node, omitMiddle bool) {
	cur := last.left
	if omitMiddle {
		cur.middle = &n.tree.root
	}
	// L M R empty
	// remove current node if all children nodes of current node are empty
	if cur.Left() == nil && cur.Right() == nil && cur.Middle() == nil {
		last.right = &n.tree.root
		return
	}

	// M nonempty
	if cur.Middle() != nil {
		return
	}
	// M is empty below

	// reconnect only one nonempty child[L|R] to be Parent Node's left node
	if cur.Left() == nil {
		last.left = cur.right
		return
	}
	if cur.Right() == nil {
		last.left = cur.left
		return
	}

	// L|R both nonempty, M empty

	// Step1. merge L onto R's left most node
	node := cur.right
	for {
		if node.Left() == nil {
			node.left = cur.left
			break
		}
		node = node.left
	}
	// Step1. move R to be Parent Node's left node
	last.left = cur.right
}

func (n *node) shrinkRight(last *node, omitMiddle bool) {
	cur := last.middle
	if omitMiddle {
		cur.middle = &n.tree.root
	}
	// remove current node if all children nodes of current node are empty
	if cur.Left() == nil && cur.Right() == nil && cur.Middle() == nil {
		last.middle = &n.tree.root
		return
	}
	// M nonempty
	if cur.Middle() != nil {
		return
	}
	// M is empty below

	// reconnect only one nonempty child[L|R] to be Parent Node's middle node
	if cur.Left() == nil {
		last.middle = cur.right
		return
	}
	if cur.Right() == nil {
		last.middle = cur.left
		return
	}

	// L|R both nonempty, M empty

	// Step1. merge R onto L's right most node
	node := cur.left
	for {
		if node.Right() == nil {
			node.right = cur.right
			break
		}
		node = node.right
	}
	// Step1. move L to be Parent Node's middle node
	last.middle = cur.left
}

func (n *node) shrinkMiddle(last *node, omitMiddle bool) {
	cur := last.middle
	if omitMiddle {
		cur.middle = &n.tree.root
	}
	// remove current node if all children nodes of current node are empty
	if cur.Left() == nil && cur.Right() == nil && cur.Middle() == nil {
		last.middle = &n.tree.root
		return
	}
	// M nonempty
	if cur.Middle() != nil {
		return
	}
	// M is empty below

	// reconnect only one nonempty child[L|R] to be Parent Node's middle node
	if cur.Left() == nil {
		last.middle = cur.right
		return
	}
	if cur.Right() == nil {
		last.middle = cur.left
		return
	}

	// L|R both nonempty, M empty

	// Step1. merge R onto L's right most node
	node := cur.left
	for {
		if node.Right() == nil {
			node.right = cur.right
			break
		}
		node = node.right
	}
	// Step1. move L to be Parent Node's middle node
	last.middle = cur.left
}

func (n *node) shrinkRoot(last *node, omitMiddle bool) {
	cur := last
	if omitMiddle {
		cur.middle = &n.tree.root
	}
	// nop if all children nodes of current node are empty
	if cur.Left() == nil && cur.Right() == nil && cur.Middle() == nil {
		return
	}
	// M nonempty
	if cur.Middle() != nil {
		return
	}
	// M is empty below

	// reconnect only one nonempty child[L|R] to be Parent Node
	if cur.Left() == nil {
		last = cur.right
		return
	}
	if cur.Right() == nil {
		last.middle = cur.left
		return
	}

	// L|R both nonempty, M empty

	// Step1. merge R onto L's right most node
	node := cur.left
	for {
		if node.Right() == nil {
			node.right = cur.right
			break
		}
		node = node.right
	}
	// Step1. move L to be Parent Node
	last = cur.left
}
