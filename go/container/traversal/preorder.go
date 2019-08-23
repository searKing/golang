// Copyright 2019 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// https://en.wikipedia.org/wiki/Tree_traversal#Post-order_(LRN)
package traversal

// Pre-order (NLR)
// 1. Check if the current node is empty or null.
// 2. Display the data part of the root (or current node).
// 3. Traverse the left subtree by recursively calling the pre-order function.
// 4. Traverse the right subtree by recursively calling the pre-order function.
// Implement:
// 	postorder(node)
// 		if (node = null)
// 			return
// 		visit(node)
// 		postorder(node.left)
// 		postorder(node.right)

// TODO template in Go2.0 is expected
// Preorder traversals from node ele by Pre-order (NLR)
// ele is a node which may have some interfaces implemented:
// LeftNode|Middleer|Righter
func Preorder(node interface{}, filterFn func(ele interface{}, depth int) (gotoNextLayer bool), processFn func(ele interface{}, depth int) (gotoNextLayer bool)) {
	traversal([]levelNode{{node: node,}}, true, preorder, filterFn, processFn)
}

// isRoot root needs to be filtered first time
func preorder(current []levelNode, filterFn func(node levelNode) (gotoNextLayer bool), processFn func(node levelNode) (gotoNextLayer bool), isRoot bool) (gotoNextLayer bool) {
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
		// process root
		if !processFn(node) {
			return false
		}
		// filter children
		preorder(filterChildren(node, node.middleLevelNodes(), filterFn), filterFn, processFn, false)

		// filter children
		preorder(filterChildren(node, node.leftLevelNodes(), filterFn), filterFn, processFn, false)
		preorder(filterChildren(node, node.rightLevelNodes(), filterFn), filterFn, processFn, false)

	}
	return true
}
