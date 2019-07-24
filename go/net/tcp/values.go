package tcp

import "errors"

// ErrServerClosed is returned by the Server's Serve and ListenAndServe
// methods after a call to Shutdown or Close.
var ErrServerClosed = errors.New("tcp: Server closed")
var ErrNotFound = errors.New("tcp: Server not found")
var ErrClientClosed = errors.New("tcp: Client closed")
var ErrUnImplement = errors.New("UnImplement Method")
