package traversal

// TODO template in Go2.0 is expected
func traversal(node interface{},
	traversalOrder traversaler,
	handler Handler) {
	traversalOrder.traversal([]levelNode{{
		node: node,
	}}, levelNodeHandlerFunc(func(ln levelNode) (gotoNextLayer bool) {
		if handler == nil {
			// traversal no node
			return false
		}
		return handler.Handle(ln.node, ln.depth)
	}))
}
