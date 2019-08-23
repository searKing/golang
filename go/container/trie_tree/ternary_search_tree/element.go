package ternary_search_tree

import "fmt"

const (
	NilKey = 0
)

type Element struct {
	// The Key stored with this element.
	Key byte
	// The value stored with this element.
	Value interface{}

	left, middle, right *Element
	tree                *TernarySearchTree
}

// Left returns the left list element or nil.
func (e *Element) Left() *Element {
	if p := e.left; e.tree != nil && p != &e.tree.root {
		return p
	}
	return nil
}

// Middle returns the middle list element or nil.
func (e *Element) Middle() *Element {
	if p := e.middle; e.tree != nil && p != &e.tree.root {
		return p
	}
	return nil
}

// Right returns the middle list element or nil.
func (e *Element) Right() *Element {
	if p := e.right; e.tree != nil && p != &e.tree.root {
		return p
	}
	return nil
}

// 前序遍历
func (e *Element) TraversalPreOrderFunc(f func(prefix []byte, value interface{}) (goon bool)) (goon bool) {
	return e.traversalPreOrderFunc([]byte{}, f)
}
func (e *Element) traversalPreOrderFunc(prefix []byte, f func(prefix []byte, value interface{}) (goon bool)) (goon bool) {
	traversal := func(e *Element) bool {
		if e != nil {
			if !e.traversalPreOrderFunc(prefix, f) {
				return false
			}
		}
		return true
	}

	curPrefix := append(prefix, e.Key)
	if !f(curPrefix, e.Value) {
		return false
	}

	if !traversal(e.Left()) {
		return false
	}
	if !traversal(e.Right()) {
		return false
	}

	prefix = curPrefix
	if !traversal(e.Middle()) {
		return false
	}

	return true
}

// 中序遍历
func (e *Element) TraversalInOrderFunc(f func(prefix []byte, value interface{}) (goon bool)) (goon bool) {
	return e.traversalInOrderFunc([]byte{}, f)
}
func (e *Element) traversalInOrderFunc(prefix []byte, f func(prefix []byte, value interface{}) (goon bool)) (goon bool) {
	traversal := func(e *Element) bool {
		if e != nil {
			if !e.traversalInOrderFunc(prefix, f) {
				return false
			}
		}
		return true
	}
	if !traversal(e.Left()) {
		return false
	}
	curPrefix := append(prefix, e.Key)
	if !f(curPrefix, e.Value) {
		return false
	}
	if !traversal(e.Right()) {
		return false
	}
	prefix = curPrefix
	if !traversal(e.Middle()) {
		return false
	}
	return true
}

// 后序遍历
func (e *Element) TraversalPostOrderFunc(f func(prefix []byte, value interface{}) (goon bool)) (goon bool) {
	return e.traversalPostOrderFunc([]byte{}, f)
}
func (e *Element) traversalPostOrderFunc(prefix []byte, f func(prefix []byte, value interface{}) (goon bool)) (goon bool) {
	traversal := func(e *Element) bool {
		if e != nil {
			if !e.traversalPostOrderFunc(prefix, f) {
				return false
			}
		}
		return true
	}
	if !traversal(e.Left()) {
		return false
	}
	if !traversal(e.Right()) {
		return false
	}
	prefix = append(prefix, e.Key)
	if !f(prefix, e.Value) {
		return false
	}
	if !traversal(e.Middle()) {
		return false
	}
	return true
}

func (e *Element) Get(prefix []byte) (value interface{}, ok bool) {
	cur := e
	for idx := 0; idx < len(prefix); {
		k := prefix[idx]
		if k < cur.Key {
			cur = cur.Left()
			if cur == nil {
				return nil, false
			}
		} else if k > cur.Key {
			cur = cur.Right()
			if cur == nil {
				return nil, false
			}
		} else {
			idx++
			if idx == len(prefix) {
				return cur.Value, true
			}
			v := cur.Value
			cur = cur.Middle()
			if cur == nil {
				return v, false
			}
		}

	}
	return nil, false
}
func (e *Element) Contains(prefix []byte) bool {
	_, ok := e.Get(prefix)
	return ok
}

const (
	posLeft = iota
	posMiddle
	posRight
	posButt
)

func (e *Element) Insert(prefix []byte, value interface{}) {
	newElement := func(key byte, value interface{}) *Element {
		return &Element{
			Key:    key,
			Value:  value,
			left:   &e.tree.root,
			middle: &e.tree.root,
			right:  &e.tree.root,
			tree:   e.tree,
		}
	}

	cur := e
	for idx := 0; idx < len(prefix); {
		// create the idx layer if not exist
		// otherwise, step to the next layer
		k := prefix[idx]
		if cur.Key == NilKey {
			cur.Key = k
		}
		// goto left
		if k < cur.Key {
			left := cur.Left()
			if left == nil {
				cur.left = newElement(k, nil)
			}
			cur = cur.left
			continue
		}
		// goto right
		if k > cur.Key {
			right := cur.Right()
			if right == nil {
				cur.right = newElement(k, nil)
			}
			cur = cur.right
			continue
		}
		// key match, goto match next layer
		idx++
		// all matched, goto set the value
		if idx == len(prefix) {
			cur.Value = value
			return
		}
		// partial matched, goto middle on next layer
		middle := cur.Middle()
		if middle == nil {
			cur.middle = newElement(NilKey, nil)
		}
		cur = cur.middle
	}
	return
}

func (e *Element) Remove(prefix []byte) (value interface{}, ok bool) {
	cur := e
	var last *Element
	var lastPos int = posButt
	for idx := 0; idx < len(prefix); {
		// create the idx layer if not exist
		// otherwise, step to the next layer
		k := prefix[idx]
		if k < cur.Key {
			left := cur.Left()
			if left == nil {
				return nil, false
			}
			last = cur
			lastPos = posLeft
			cur = left
			continue
		}
		if k > cur.Key {
			right := cur.Right()
			if right == nil {
				return nil, false
			}
			last = cur
			lastPos = posRight
			cur = right
			continue
		}
		idx++
		if idx == len(prefix) {
			// match
			switch lastPos {
			case posLeft:
				last.left = &e.tree.root
			case posMiddle:
				last.middle = &e.tree.root
			case posRight:
				last.right = &e.tree.root
			case posButt:
				last.left = &e.tree.root
				last.middle = &e.tree.root
				last.right = &e.tree.root
			}
			return cur.Value, true
		}
		middle := cur.Middle()
		if middle == nil {
			return nil, false
		}
		last = cur
		lastPos = posMiddle
		cur = middle
	}
	return nil, false
}
func (e *Element) String() string {
	s := ""
	e.TraversalInOrderFunc(func(prefix []byte, value interface{}) (goon bool) {
		s += fmt.Sprintf("%s:%v\n", string(prefix), value)
		return true
	})
	return s
}
