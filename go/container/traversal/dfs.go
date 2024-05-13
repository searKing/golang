// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// https://en.wikipedia.org/wiki/Depth-first_search
package traversal

// Depth-first search (DFS)
// 1. Check if the current node is empty or null.
// 2. Display the data part of the root (or current node).
// 3. Traverse the middle subtree by recursively calling the dfs-order function.
// 4. Traverse the left subtree by recursively calling the dfs-order function.
// 5. Traverse the right subtree by recursively calling the dfs-order function.
// Implement:
//	procedure DFS(G,v):
//		label v as discovered
//		for all directed edges from v to w that are in G.adjacentEdges(v) do
//			if vertex w is not labeled as discovered then
//				recursively call DFS(G,w)

// TODO template in Go2.0 is expected
// DepthFirstSearchOrder traversals from node ele by Depth-first search (DFS)
// ele is a node which may have some interfaces implemented:
// LeftNodes|MiddleNodes|RightNodes
func DepthFirstSearchOrder(node any, handler Handler) {
	traversal(node, traversalerFunc(dfs), handler)
}

func dfs(currents []levelNode, handler levelNodeHandler) (goon bool) {
	if len(currents) == 0 {
		return true
	}
	// Step 1: brothers
	for _, node := range currents {
		if node.visited {
			continue
		}
		if !handler.Handle(node) {
			return false
		}
		if !dfs(node.middleLevelNodes(), handler) {
			return false
		}
		if !dfs(node.leftLevelNodes(), handler) {
			return false
		}
		if !dfs(node.rightLevelNodes(), handler) {
			return false
		}
	}
	return true
}
