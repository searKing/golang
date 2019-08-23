// Copyright 2019 The searKing Author. All rights reserved.
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
// LeftNode|Middleer|Righter
func DepthFirstSearchOrder(node interface{}, filterFn func(ele interface{}, depth int) (gotoNextLayer bool), processFn func(ele interface{}, depth int) (gotoNextLayer bool)) {
	traversal([]levelNode{{node: node,}}, true, dfs, filterFn, processFn)
}

// isRoot root needs to be filtered first time
func dfs(current []levelNode, filterFn func(node levelNode) (gotoNextLayer bool), processFn func(node levelNode) (gotoNextLayer bool), isRoot bool) (gotoNextLayer bool) {
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
		if !processFn(node) {
			return false
		}
		// filter children
		dfs(filterChildren(node, node.middleLevelNodes(), filterFn), filterFn, processFn, false)
		dfs(filterChildren(node, node.leftLevelNodes(), filterFn), filterFn, processFn, false)
		dfs(filterChildren(node, node.rightLevelNodes(), filterFn), filterFn, processFn, false)
	}
	return true
}
