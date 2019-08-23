package traversal

// TODO template in Go2.0 is expected
// Depth First Search
func DFS(ele interface{}, filterFn func(ele interface{}, depth int) (gotoNextLayer bool), processFn func(ele interface{}, depth int) (gotoNextLayer bool)) {
	traversalDFS([]Node{{
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
func traversalDFS(current []Node, filterFn func(node Node) (gotoNextLayer bool), processFn func(node Node) (gotoNextLayer bool), isRoot bool) (gotoNextLayer bool) {
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
		traversalDFS(filterChildren(node, node.MiddleNodes(), filterFn), filterFn, processFn, false)
		traversalDFS(filterChildren(node, node.LeftNodes(), filterFn), filterFn, processFn, false)
		traversalDFS(filterChildren(node, node.RightNodes(), filterFn), filterFn, processFn, false)
	}
	return true
}
