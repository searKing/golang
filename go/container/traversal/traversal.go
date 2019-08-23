package traversal

// TODO template in Go2.0 is expected
func traversal(node interface{}, isRoot bool,
	traversalOrder func(current []levelNode, filterFn func(node levelNode) (gotoNextLayer bool), processFn func(node levelNode) (gotoNextLayer bool), isRoot bool) (gotoNextLayer bool),
	filterFn func(node interface{}, depth int) (gotoNextLayer bool),
	processFn func(node interface{}, depth int) (gotoNextLayer bool)) {
	traversalOrder([]levelNode{{
		node: node,
	}}, func(ln levelNode) (gotoNextLayer bool) {
		if filterFn == nil {
			// traversal every node
			return true
		}
		return filterFn(ln.node, ln.depth)
	}, func(ln levelNode) (gotoNextLayer bool) {
		if processFn == nil {
			// traversal no node
			return false
		}
		return processFn(ln.node, ln.depth)
	}, isRoot)
}
