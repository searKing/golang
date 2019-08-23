// Copyright 2019 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// https://en.wikipedia.org/wiki/Tree_traversal#In-order_(LNR)
package traversal

// In-order (LNR)
// 1. Check if the current node is empty or null.
// 2. Traverse the left subtree by recursively calling the in-order function.
// 3. Display the data part of the root (or current node).
// 4. Traverse the right subtree by recursively calling the in-order function.
// Implement:
// 	inorder(node)
// 		if (node = null)
// 			return
// 		inorder(node.left)
// 		visit(node)
// 		inorder(node.right)

// TODO template in Go2.0 is expected
// Inorder traversals from node ele by In-order (LNR)
// ele is a node which may have some interfaces implemented:
// LeftNode|Middleer|Righter
func Inorder(node interface{}, filterFn func(ele interface{}, depth int) (gotoNextLayer bool), processFn func(ele interface{}, depth int) (gotoNextLayer bool)) {
	traversal([]levelNode{{node: node,}}, true, inorder, filterFn, processFn)
}

// isRoot root needs to be filtered first time
func inorder(current []levelNode, filterFn func(node levelNode) (gotoNextLayer bool), processFn func(node levelNode) (gotoNextLayer bool), isRoot bool) (gotoNextLayer bool) {
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
		inorder(filterChildren(node, node.leftLevelNodes(), filterFn), filterFn, processFn, false)

		// process root
		if !processFn(node) {
			return false
		}
		// filter children
		inorder(filterChildren(node, node.middleLevelNodes(), filterFn), filterFn, processFn, false)

		inorder(filterChildren(node, node.rightLevelNodes(), filterFn), filterFn, processFn, false)

	}
	return true
}
