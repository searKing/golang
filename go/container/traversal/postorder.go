// Copyright 2020 The searKing Author. All rights reserved.
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
// LeftNodes|MiddleNodes|RightNodes
func Postorder(node interface{}, handler Handler) {
	traversal(node, traversalerFunc(postorder), handler)
}

func postorder(currents []levelNode, handler levelNodeHandler) (goon bool) {
	if len(currents) == 0 {
		return true
	}
	// Step 1: brothers
	for _, node := range currents {
		if node.visited {
			continue
		}
		// process children
		if !outorder(node.leftLevelNodes(), handler) {
			return false
		}
		if !outorder(node.rightLevelNodes(), handler) {
			return false
		}
		// process root
		if !handler.Handle(node) {
			return false
		}
		// process children
		if !outorder(node.middleLevelNodes(), handler) {
			return false
		}
	}
	return true
}
