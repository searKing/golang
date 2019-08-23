// Copyright 2019 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package traversal

type LeftNode interface {
	// Left returns the left node list or nil.
	Lefts() []interface{}
}

type MiddleNode interface {
	// Middle returns the middle node list or nil.
	Middles() []interface{}
}

type RightNode interface {
	// Right returns the middle node list or nil.
	Rights() []interface{}
}

// levelNode represents a single node with depth found in a structure.
type levelNode struct {
	node  interface{}
	depth int
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
	left, ok := n.node.(LeftNode)
	if ok {
		return left.Lefts()
	}
	return nil
}

func (n *levelNode) middleNodes() []interface{} {
	if n.node == nil {
		return nil
	}
	middle, ok := n.node.(MiddleNode)
	if ok {
		return middle.Middles()
	}
	return nil

}

func (n *levelNode) rightNodes() []interface{} {
	if n.node == nil {
		return nil
	}
	right, ok := n.node.(RightNode)
	if ok {
		return right.Rights()
	}
	return nil

}
