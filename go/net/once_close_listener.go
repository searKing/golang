package net

import (
	"net"
	"sync"
)

type onceCloseListener struct {
	net.Listener
	once     sync.Once
	closeErr error
}

func (oc *onceCloseListener) Close() error {
	oc.once.Do(oc.close)
	return oc.closeErr
}

func (oc *onceCloseListener) close() { oc.closeErr = oc.Listener.Close() }

// OnceCloseListener wraps a net.Listener, protecting it from
// multiple Close calls.
func OnceCloseListener(l net.Listener) net.Listener {
	return &onceCloseListener{Listener: l}
}
