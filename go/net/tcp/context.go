package tcp

var (
	// ServerContextKey is a context key. It can be used in TCP
	// handlers with context.WithValue to access the server that
	// started the handler. The associated value will be of
	// type *Server.
	ServerContextKey = &contextKey{"tcp-server"}
	// ClientContextKey is a context key. It can be used in TCP
	// handlers with context.WithValue to access the client that
	// started the handler. The associated value will be of
	// type *Client.
	ClientContextKey = &contextKey{"tcp-client"}

	// LocalAddrContextKey is a context key. It can be used in
	// TCP handlers with context.WithValue to access the local
	// address the connection arrived on.
	// The associated value will be of type net.Addr.
	LocalAddrContextKey = &contextKey{"local-addr"}
)

type contextKey struct {
	name string
}

func (k *contextKey) String() string { return "net/http context value " + k.name }
