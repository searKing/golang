// Copyright 2019 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// https://en.wikipedia.org/wiki/Breadth-first_search
package traversal

// Breadth-first search (BFS)
// 1. Check if the current depth level is empty or null.
// 2. Display the data part of all of the neighbor nodes at the present depth.
// 2. Traverse the next depth level by recursively calling the bfs-order function.
// Implement:
//	procedure BreadthFirstSearch(G,start_v):
//		let Q be a queue
//		label start_v as discovered
//		Q.enqueue(start_v)
//		while Q is not empty
//			v = Q.dequeue()
//			if v is the goal:
//				return v
//			for all edges from v to w in G.adjacentEdges(v) do
//				if w is not labeled as discovered:
//					label w as discovered
//					w.parent = v
//					Q.enqueue(w)

// TODO template in Go2.0 is expected
// BreadthFirstSearchOrder traversals from node ele by Breadth-first search (BFS)
// ele is a node which may have some interfaces implemented:
// LeftNode|Middleer|Righter
func BreadthFirstSearchOrder(node interface{}, filterFn func(node interface{}, depth int) (gotoNextLayer bool), processFn func(node interface{}, depth int) (gotoNextLayer bool)) {
	traversal([]levelNode{{node: node,}}, true, bfs, filterFn, processFn)
}

// isRoot root needs to be filtered first time
func bfs(current []levelNode, filterFn func(node levelNode) (gotoNextLayer bool), processFn func(node levelNode) (gotoNextLayer bool), isRoot bool) (gotoNextLayer bool) {
	if len(current) == 0 {
		return false
	}
	// Step 1: brothers layer
	var nextBrothers []levelNode
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
		// filter brothers
		nextBrothers = append(nextBrothers, node)
	}

	// Step 2: children layer
	var nextChildren []levelNode
	// filter children
	for _, node := range nextBrothers {
		// Scan node for nodes to include.
		nextChildren = append(nextChildren, filterChildren(node, node.leftLevelNodes(), filterFn)...)
		nextChildren = append(nextChildren, filterChildren(node, node.middleLevelNodes(), filterFn)...)
		nextChildren = append(nextChildren, filterChildren(node, node.rightLevelNodes(), filterFn)...)
	}
	bfs(nextChildren, filterFn, processFn, false)
	return true
}
