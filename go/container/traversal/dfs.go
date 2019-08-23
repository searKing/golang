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
// Lefter|Middleer|Righter
// Lefters|Middleers|Righters
func DepthFirstSearchOrder(ele interface{}, filterFn func(ele interface{}, depth int) (gotoNextLayer bool), processFn func(ele interface{}, depth int) (gotoNextLayer bool)) {
	dfs([]Node{{
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
func dfs(current []Node, filterFn func(node Node) (gotoNextLayer bool), processFn func(node Node) (gotoNextLayer bool), isRoot bool) (gotoNextLayer bool) {
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
		dfs(filterChildren(node, node.MiddleNodes(), filterFn), filterFn, processFn, false)
		dfs(filterChildren(node, node.LeftNodes(), filterFn), filterFn, processFn, false)
		dfs(filterChildren(node, node.RightNodes(), filterFn), filterFn, processFn, false)
	}
	return true
}
