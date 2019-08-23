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
// Postorder traversals from node ele by Post-order (LRN)
// ele is a node which may have some interfaces implemented:
// LeftNode|Middleer|Righter
func Postorder(node interface{}, filterFn func(ele interface{}, depth int) (gotoNextLayer bool), processFn func(ele interface{}, depth int) (gotoNextLayer bool)) {
	traversal([]levelNode{{node: node,}}, true, postorder, filterFn, processFn)
}

// isRoot root needs to be filtered first time
func postorder(current []levelNode, filterFn func(node levelNode) (gotoNextLayer bool), processFn func(node levelNode) (gotoNextLayer bool), isRoot bool) (gotoNextLayer bool) {
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
		postorder(filterChildren(node, node.leftLevelNodes(), filterFn), filterFn, processFn, false)
		postorder(filterChildren(node, node.rightLevelNodes(), filterFn), filterFn, processFn, false)

		// process root
		if !processFn(node) {
			return false
		}
		// filter children
		postorder(filterChildren(node, node.middleLevelNodes(), filterFn), filterFn, processFn, false)
	}
	return true
}
