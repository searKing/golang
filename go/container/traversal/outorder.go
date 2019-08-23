// Copyright 2019 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// https://en.wikipedia.org/wiki/Tree_traversal#Out-order_(RNL)
package traversal

// Out-order (RNL)
// 1. Check if the current node is empty or null.
// 2. Traverse the right subtree by recursively calling the out-order function.
// 3. Display the data part of the root (or current node).
// 4. Traverse the left subtree by recursively calling the out-order function.
// Implement:
// 	outorder(node)
// 		if (node = null)
// 			return
// 		outorder(node.right)
// 		visit(node)
// 		outorder(node.left)

// TODO template in Go2.0 is expected
// Outorder traversals from node ele by Out-order (RNL)
// ele is a node which may have some interfaces implemented:
// LeftNode|Middleer|Righter
func Outorder(ele interface{}, filterFn func(ele interface{}, depth int) (gotoNextLayer bool), processFn func(ele interface{}, depth int) (gotoNextLayer bool)) {
	outorder([]levelNode{{
		node: ele,
	}}, func(node levelNode) (gotoNextLayer bool) {
		if filterFn == nil {
			// traversal every node
			return true
		}
		return filterFn(node.node, node.depth)
	}, func(node levelNode) (gotoNextLayer bool) {
		if processFn == nil {
			// traversal no node
			return false
		}
		return processFn(node.node, node.depth)
	}, true)
}

// isRoot root needs to be filtered first time
func outorder(current []levelNode, filterFn func(node levelNode) (gotoNextLayer bool), processFn func(node levelNode) (gotoNextLayer bool), isRoot bool) (gotoNextLayer bool) {
	if len(current) == 0 {
		return false
	}
	// Step 1: brothers
	for _, node := range current {
		// filter root
		if isRoot {
			if !filterFn(node) {
				return false
			}
		}
		// filter children
		outorder(filterChildren(node, node.rightLevelNodes(), filterFn), filterFn, processFn, false)

		// process root
		if !processFn(node) {
			return false
		}
		// filter children
		outorder(filterChildren(node, node.middleLevelNodes(), filterFn), filterFn, processFn, false)

		outorder(filterChildren(node, node.leftLevelNodes(), filterFn), filterFn, processFn, false)

	}
	return true
}
