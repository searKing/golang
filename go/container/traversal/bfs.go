// Copyright 2020 The searKing Author. All rights reserved.
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
// LeftNodes|MiddleNodes|RightNodes
func BreadthFirstSearchOrder(node any, handler Handler) {
	traversal(node, traversalerFunc(bfs), handler)
}

func bfs(currents []levelNode, handler levelNodeHandler) (goon bool) {
	if len(currents) == 0 {
		return true
	}
	// Step 1: brothers layer
	var nextBrothers []levelNode
	for _, node := range currents {
		if node.visited {
			continue
		}
		if !handler.Handle(node) {
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
		nextChildren = append(nextChildren, node.leftLevelNodes()...)
		nextChildren = append(nextChildren, node.middleLevelNodes()...)
		nextChildren = append(nextChildren, node.rightLevelNodes()...)
	}

	return bfs(nextChildren, handler)
}
