package traversal

func filterChildren(root Node, children []Node, filterFn func(node Node) (truth bool)) (nextChildren []Node) {
	// Scan node for nodes to include.
	for _, n := range children {
		if filterFn(n) {
			nextChildren = append(nextChildren, n)
		}
	}
	return nextChildren
}
