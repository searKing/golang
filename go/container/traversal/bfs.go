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
func BreadthFirstSearchOrder(ele interface{}, filterFn func(ele interface{}, depth int) (gotoNextLayer bool), processFn func(ele interface{}, depth int) (gotoNextLayer bool)) {
	bfs([]Node{{
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
func bfs(current []Node, filterFn func(node Node) (gotoNextLayer bool), processFn func(node Node) (gotoNextLayer bool), isRoot bool) (gotoNextLayer bool) {
	if len(current) == 0 {
		return false
	}
	// Step 1: brothers layer
	var nextBrothers []Node
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
	var nextChildren []Node
	// filter children
	for _, node := range nextBrothers {
		// Scan node for nodes to include.
		nextChildren = append(nextChildren, filterChildren(node, node.LeftNodes(), filterFn)...)
		nextChildren = append(nextChildren, filterChildren(node, node.MiddleNodes(), filterFn)...)
		nextChildren = append(nextChildren, filterChildren(node, node.RightNodes(), filterFn)...)
	}
	bfs(nextChildren, filterFn, processFn, false)
	return true
}
