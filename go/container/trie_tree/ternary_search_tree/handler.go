package ternary_search_tree

type Handler interface {
	Handle(prefix []byte, value interface{}) (goon bool)
}

// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as traversal handlers. If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler that calls f.
type HandlerFunc func(prefix []byte, value interface{}) (goon bool)

// ServeHTTP calls f(w, r).
func (f HandlerFunc) Handle(prefix []byte, value interface{}) (goon bool) {
	return f(prefix, value)
}
