// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"errors"
	"time"
)

// ErrBadResolverState may be returned by UpdateState to indicate a
// problem with the provided name resolver data.
var ErrBadResolverState = errors.New("bad resolver state")

// Build includes additional information for the builder to create
// the resolver.
//
//go:generate go-option -type "Build"
type Build struct {
	ClientConn ClientConn
}

// Builder is the interface that must be implemented by a database
// driver.
//
// Database drivers may implement DriverContext for access
// to contexts and to parse the name only once for a pool of connections,
// instead of once per connection.
type Builder interface {
	// Build creates a new resolver for the given target.
	//
	// gRPC dial calls Build synchronously, and fails if the returned error is
	// not nil.
	Build(ctx context.Context, target Target, opts ...BuildOption) (Resolver, error)

	// Scheme returns the scheme supported by this resolver.
	// Scheme is defined at https://github.com/grpc/grpc/blob/master/doc/naming.md.
	Scheme() string
}

// resolveOneAddr includes additional information for ResolveOneAddr.
//
//go:generate go-option -type "resolveOneAddr"
type resolveOneAddr struct {
	Picker []PickOption
}

// resolveAddr includes additional information for ResolveAddr.
//
//go:generate go-option -type "resolveAddr"
type resolveAddr struct{}

// resolveNow includes additional information for ResolveNow.
//
//go:generate go-option -type "resolveNow"
type resolveNow struct{}

// resolveDone includes additional information for ResolveDone.
//
//go:generate go-option -type "resolveDone"
type resolveDone struct{}

// Resolver watches for the updates on the specified target.
// Updates include address updates and service config updates.
type Resolver interface {
	// ResolveOneAddr will be called to try to resolve the target name directly.
	// resolver can not ignore this if it's not necessary.
	// ResolveOneAddr may trigger and wait for ResolveNow if no addr in resolver cache
	ResolveOneAddr(ctx context.Context, opts ...ResolveOneAddrOption) (Address, error)
	// ResolveAddr will be called to try to resolve the target name directly.
	// resolver can not ignore this if it's not necessary.
	// ResolveAddr may trigger and wait for ResolveNow if no addr in resolver cache
	ResolveAddr(ctx context.Context, opts ...ResolveAddrOption) ([]Address, error)
	// ResolveNow will be called to try to resolve the target name
	// again. It's just a hint, resolver can ignore this if it's not necessary.
	//
	// It could be called multiple times concurrently.
	// It may trigger ClientConn to UpdateState or ReportError if failed.
	// It may update cache used by ResolveOneAddr or ResolveAddr.
	ResolveNow(ctx context.Context, opts ...ResolveNowOption)
	// Close closes the resolver.
	Close()
}

// ResolveDoneResolver extends Resolver with ResolveDone
type ResolveDoneResolver interface {
	Resolver
	// ResolveDone will be called after the RPC finished.
	// resolver can ignore this if it's not necessary.
	ResolveDone(ctx context.Context, doneInfo DoneInfo, opts ...ResolveDoneOption)
}

// State contains the current Resolver state relevant to the ClientConn.
type State struct {
	// Addresses is the latest set of resolved addresses for the target.
	Addresses []Address
}

// ClientConn contains the callbacks for resolver to notify any updates
// to the gRPC ClientConn.
//
// This interface is to be implemented by gRPC. Users should not need a
// brand new implementation of this interface. For the situations like
// testing, the new implementation should embed this interface. This allows
// gRPC to add new methods to this interface.
type ClientConn interface {
	// UpdateState updates the state of the ClientConn appropriately.
	UpdateState(State) error
	// ReportError notifies the ClientConn that the Resolver encountered an
	// error.  The ClientConn will notify the load balancer and begin calling
	// ResolveNow on the Resolver with exponential backoff.
	ReportError(error)
}

// Address represents a server the client connects to.
type Address struct {
	// Addr is the server address on which a connection will be established.
	Addr string

	// ResolveLoad is the load received from resolver.
	ResolveLoad any
}

// DoneInfo contains additional information for done.
type DoneInfo struct {
	// Err is the rpc error the RPC finished with. It could be nil.
	// usually io.EOF represents server addr is resolved but unacceptable.
	Err error

	// Addr represents a server the client connects to.
	Addr Address

	Duration time.Duration
}
