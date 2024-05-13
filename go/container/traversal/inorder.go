// Copyright 2020 The searKing Author. All rights reserved.
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
// LeftNodes|MiddleNodes|RightNodes
func Inorder(node any, handler Handler) {
	traversal(node, traversalerFunc(inorder), handler)
}

func inorder(currents []levelNode, handler levelNodeHandler) (goon bool) {
	if len(currents) == 0 {
		return true
	}
	// Step 1: brothers
	for _, node := range currents {
		if node.visited {
			continue
		}
		// process children
		if !inorder(node.leftLevelNodes(), handler) {
			return false
		}

		// process root
		if !handler.Handle(node) {
			return false
		}
		// process children
		if !inorder(node.middleLevelNodes(), handler) {
			return false
		}
		if !inorder(node.rightLevelNodes(), handler) {
			return false
		}
	}
	return true
}
