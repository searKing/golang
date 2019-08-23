package traversal

func filterChildren(root levelNode, children []levelNode, filterFn func(node levelNode) (truth bool)) (nextChildren []levelNode) {
	// Scan node for nodes to include.
	for _, n := range children {
		if filterFn(n) {
			nextChildren = append(nextChildren, n)
		}
	}
	return nextChildren
}
