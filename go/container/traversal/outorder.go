// Copyright 2020 The searKing Author. All rights reserved.
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
// LeftNodes|MiddleNodes|RightNodes
func Outorder(node interface{}, handler Handler) {
	traversal(node, traversalerFunc(outorder), handler)
}

func outorder(currents []levelNode, handler levelNodeHandler) (goon bool) {
	if len(currents) == 0 {
		return true
	}
	// Step 1: brothers
	for _, node := range currents {
		if node.visited {
			continue
		}
		// process children
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
		if !outorder(node.leftLevelNodes(), handler) {
			return false
		}
	}
	return true
}
