package ternary_search_tree

import (
	"fmt"
	"github.com/searKing/golang/go/container/traversal"
	"github.com/searKing/golang/go/util"
)

const (
	NilKey = 0
)

type node struct {
	prefix   []byte
	key      byte
	hasKey   bool
	value    interface{}
	hasValue bool

	left, middle, right *node
	tree                *ternarySearchTree
}

func (n *node) LeftNodes() []interface{} {
	left := n.Left()
	if left == nil {
		return nil
	}
	return []interface{}{left}
}

func (n *node) MiddleNodes() []interface{} {
	middle := n.Middle()
	if middle == nil {
		return nil
	}
	return []interface{}{middle}
}

func (n *node) RightNodes() []interface{} {
	right := n.Right()
	if right == nil {
		return nil
	}
	return []interface{}{right}
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

// Right returns the middle list node or nil.
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
	order(n, traversal.HandlerFunc(func(ele interface{}, depth int) (goon bool) {
		currentNode := ele.(*node)
		if !currentNode.hasKey || !currentNode.hasValue {
			return true
		}
		return handler.Handle(currentNode.prefix, currentNode.value)
	}))
	return
}

func (n *node) Load(prefix []byte) (value interface{}, ok bool) {
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

const (
	posLeft = iota
	posMiddle
	posRight
	posButt
)

func (n *node) Store(prefix []byte, value interface{}) {
	// force update
	n.CAS(prefix, nil, value, util.AlwaysEqualComparator())
}

func (n *node) Remove(prefix []byte, shrinkToFit bool) (old interface{}, ok bool) {
	cur, last, lastPos, has := n.search(prefix)
	if !has {
		return nil, false
	}
	if !cur.hasValue {
		return nil, false
	}
	cur.hasValue = false
	// shrinkToFit if cur's children are empty
	if shrinkToFit {
		cur.shrinkToFit(last, lastPos)
	}
	// all matched, goto remove the old
	return cur.value, true
}

func (n *node) RemoveAll(prefix []byte) (value interface{}, ok bool) {
	cur, last, lastPos, has := n.search(prefix)
	if !has {
		return nil, false
	}

	// match
	switch lastPos {
	case posLeft:
		last.left = &n.tree.root
	case posMiddle:
		last.middle = &n.tree.root
	case posRight:
		last.right = &n.tree.root
	case posButt:
		last.left = &n.tree.root
		last.middle = &n.tree.root
		last.right = &n.tree.root
	}

	if !cur.hasValue {
		return nil, false
	}
	cur.hasValue = false
	// all matched, goto remove the value
	return cur.value, true
}

func (n *node) String() string {
	s := ""
	n.Traversal(traversal.Inorder, HandlerFunc(func(prefix []byte, value interface{}) (goon bool) {
		s += fmt.Sprintf("%s:%v\n", string(prefix), value)
		return true
	}))

	return s
}

func (n *node) CAS(prefix []byte, old, new interface{}, cmps ...util.Comparator) bool {
	newElement := func(prefix []byte, hasKey bool, key byte, hasValue bool, value interface{}) *node {
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
			var cmp util.Comparator
			if len(cmps) > 0 {
				cmp = cmps[0]
			}
			if cmp == nil {
				cmp = util.DefaultComparator()
			}
			if cmp.Compare(cur.value, old) == 0 {
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
	n.Traversal(traversal.Preorder, HandlerFunc(func(prefix []byte, value interface{}) (goon bool) {
		if depth < len(prefix) {
			depth = len(prefix)
		}
		return true
	}))
	return depth
}

// shrinkToFit cutoff lastnode's children nodes if all children nodes are empty
func (n *node) shrinkToFit(last *node, lastPos int) {
	var has bool
	n.Traversal(traversal.Preorder, HandlerFunc(func(prefix []byte, value interface{}) (goon bool) {
		has = true
		return false
	}))
	if !has {
		return
	}

	// match
	switch lastPos {
	case posLeft:
		last.left = &n.tree.root
	case posMiddle:
		last.middle = &n.tree.root
	case posRight:
		last.right = &n.tree.root
	case posButt:
		last.left = &n.tree.root
		last.middle = &n.tree.root
		last.right = &n.tree.root
	}
}

// return true if prefix matches, no matter value exists
func (n *node) search(prefix []byte) (cur, last *node, lastPos int, has bool) {
	cur = n
	last = n

	lastPos = posButt
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
