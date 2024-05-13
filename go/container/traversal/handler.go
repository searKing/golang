// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package traversal

type Order func(node any, handler Handler)

type Handler interface {
	Handle(node any, depth int) (goon bool)
}

// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as traversal handlers. If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler that calls f.
type HandlerFunc func(node any, depth int) (goon bool)

// ServeHTTP calls f(w, r).
func (f HandlerFunc) Handle(node any, depth int) (goon bool) {
	return f(node, depth)
}

type traversaler interface {
	traversal(currents []levelNode, handler levelNodeHandler) (goon bool)
}

// The traversalerFunc type is an adapter to allow the use of
// ordinary functions as traversal handlers. If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler that calls f.
type traversalerFunc func(currents []levelNode, handler levelNodeHandler) (goon bool)

// ServeHTTP calls f(w, r).
func (f traversalerFunc) traversal(currents []levelNode, handler levelNodeHandler) (goon bool) {
	return f(currents, handler)
}

type levelNodeHandler interface {
	Handle(node levelNode) (goon bool)
}

// The levelNodeHandlerFunc type is an adapter to allow the use of
// ordinary functions as traversal handlers. If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler that calls f.
type levelNodeHandlerFunc func(node levelNode) (goon bool)

// ServeHTTP calls f(w, r).
func (f levelNodeHandlerFunc) Handle(node levelNode) (goon bool) {
	return f(node)
}
