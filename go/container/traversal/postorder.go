// Copyright 2019 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// https://en.wikipedia.org/wiki/Tree_traversal#Post-order_(LRN)
package traversal

// Post-order (LRN)
// 1. Check if the current node is empty or null.
// 2. Traverse the left subtree by recursively calling the post-order function.
// 3. Traverse the right subtree by recursively calling the post-order function.
// 4. Display the data part of the root (or current node).
// Implement:
// 	postorder(node)
// 		if (node = null)
// 			return
// 		postorder(node.left)
// 		postorder(node.right)
// 		visit(node)
// TODO template in Go2.0 is expected
func Postorder(ele interface{}, filterFn func(ele interface{}, depth int) (gotoNextLayer bool), processFn func(ele interface{}, depth int) (gotoNextLayer bool)) {
	postorder([]Node{{
		ele: ele,
	}}, func(node Node) (gotoNextLayer bool) {
		if filterFn == nil {
			// traversal every node
			return true
		}
		return filterFn(node.ele, node.depth)
	}, func(node Node) (gotoNextLayer bool) {
		if processFn == nil {
			// traversal no node
			return false
		}
		return processFn(node.ele, node.depth)
	}, true)
}

// isRoot root needs to be filtered first time
func postorder(current []Node, filterFn func(node Node) (gotoNextLayer bool), processFn func(node Node) (gotoNextLayer bool), isRoot bool) (gotoNextLayer bool) {
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
		postorder(filterChildren(node, node.LeftNodes(), filterFn), filterFn, processFn, false)
		postorder(filterChildren(node, node.RightNodes(), filterFn), filterFn, processFn, false)

		// process root
		if !processFn(node) {
			return false
		}
		// filter children
		postorder(filterChildren(node, node.MiddleNodes(), filterFn), filterFn, processFn, false)
	}
	return true
}
