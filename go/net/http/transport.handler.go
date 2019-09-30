package http

import "net/http"

// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as HTTP handlers. If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler that calls f.
type RoundTripFunc func(req *http.Request) (resp *http.Response, err error)

func (f RoundTripFunc) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	return f(req)
}

type Handler interface {
	RoundTrip(req *http.Request) (resp *http.Response, err error)
}

// HandlersChain defines a HandlerFunc array.
type HandlersChain []Handler

// Last returns the last handler in the chain. ie. the last handler is the main own.
func (c HandlersChain) Last() Handler {
	if length := len(c); length > 0 {
		return c[length-1]
	}
	return nil
}
