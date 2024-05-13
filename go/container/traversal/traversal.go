// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package traversal

// TODO template in Go2.0 is expected
func traversal(node any,
	traversalOrder traversaler,
	handler Handler) {
	traversalOrder.traversal([]levelNode{{
		node: node,
	}}, levelNodeHandlerFunc(func(ln levelNode) (goon bool) {
		if handler == nil {
			// traversal no node
			return false
		}
		return handler.Handle(ln.node, ln.depth)
	}))
}
