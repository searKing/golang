package traversal

// TODO template in Go2.0 is expected
// Breadth First Search
func TraversalBFS(ele interface{}, filterFn func(ele interface{}, depth int) (gotoNextLayer bool), processFn func(ele interface{}, depth int) (gotoNextLayer bool)) {
	traversalBFS([]Node{{
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
func traversalBFS(current []Node, filterFn func(node Node) (gotoNextLayer bool), processFn func(node Node) (gotoNextLayer bool), isRoot bool) (gotoNextLayer bool) {
	if len(current) == 0 {
		return false
	}
	// Step 1: brothers layer
	nextBrothers := []Node{}
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
	nextChildren := []Node{}
	// filter children
	for _, node := range nextBrothers {
		// Scan node for nodes to include.
		nextChildren = append(nextChildren, filterChildren(node, node.LeftNodes(), filterFn)...)
		nextChildren = append(nextChildren, filterChildren(node, node.MiddleNodes(), filterFn)...)
		nextChildren = append(nextChildren, filterChildren(node, node.RightNodes(), filterFn)...)
	}
	traversalBFS(nextChildren, filterFn, processFn, false)
	return true
}
