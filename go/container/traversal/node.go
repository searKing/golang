// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package traversal

type LeftNodes interface {
	// Left returns the left node list or nil.
	LeftNodes() []interface{}
}

type MiddleNodes interface {
	// Middle returns the middle node list or nil.
	MiddleNodes() []interface{}
}

type RightNodes interface {
	// Right returns the middle node list or nil.
	RightNodes() []interface{}
}

// a uniform node
type Node interface {
	LeftNodes
	MiddleNodes
	RightNodes
}

// levelNode represents a single node with depth found in a structure.
type levelNode struct {
	node    interface{}
	depth   int
	visited bool
}

func (n *levelNode) leftLevelNodes() []levelNode {
	var lefts []levelNode
	for _, node := range n.leftNodes() {
		ln := levelNode{
			node:  node,
			depth: n.depth + 1,
		}
		lefts = append(lefts, ln)
	}
	return lefts
}

func (n *levelNode) middleLevelNodes() []levelNode {
	var middles []levelNode
	for _, node := range n.middleNodes() {
		ln := levelNode{
			node:  node,
			depth: n.depth + 1,
		}
		middles = append(middles, ln)
	}
	return middles
}

func (n *levelNode) rightLevelNodes() []levelNode {
	var rights []levelNode
	for _, node := range n.rightNodes() {
		ln := levelNode{
			node:  node,
			depth: n.depth + 1,
		}
		rights = append(rights, ln)
	}
	return rights
}

// children
func (n *levelNode) leftNodes() []interface{} {
	if n.node == nil {
		return nil
	}
	left, ok := n.node.(LeftNodes)
	if ok {
		return left.LeftNodes()
	}
	return nil
}

func (n *levelNode) middleNodes() []interface{} {
	if n.node == nil {
		return nil
	}
	middle, ok := n.node.(MiddleNodes)
	if ok {
		return middle.MiddleNodes()
	}
	return nil

}

func (n *levelNode) rightNodes() []interface{} {
	if n.node == nil {
		return nil
	}
	right, ok := n.node.(RightNodes)
	if ok {
		return right.RightNodes()
	}
	return nil

}
