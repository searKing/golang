// Copyright 2020 The searKing Author. All rights reserved.
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
// LeftNodes|MiddleNodes|RightNodes
func Preorder(node any, handler Handler) {
	traversal(node, traversalerFunc(preorder), handler)
}

func preorder(currents []levelNode, handler levelNodeHandler) (goon bool) {
	if len(currents) == 0 {
		return true
	}
	// Step 1: brothers
	for _, node := range currents {
		if node.visited {
			continue
		}
		// process root
		if !handler.Handle(node) {
			return false
		}
		// traversal children
		if !preorder(node.middleLevelNodes(), handler) {
			return false
		}
		if !preorder(node.leftLevelNodes(), handler) {
			return false
		}
		if !preorder(node.rightLevelNodes(), handler) {
			return false
		}
	}
	return true
}
