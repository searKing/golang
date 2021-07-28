// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"errors"
)

// ErrBadResolverState may be returned by UpdateState to indicate a
// problem with the provided name resolver data.
var ErrBadResolverState = errors.New("bad resolver state")

// Target represents a target for gRPC, as specified in:
// https://github.com/grpc/grpc/blob/master/doc/naming.md.
// It is parsed from the target string that gets passed into Dial or DialContext by the user. And
// grpc passes it to the resolver and the balancer.
//
// If the target follows the naming spec, and the parsed scheme is registered with grpc, we will
// parse the target string according to the spec. e.g. "dns://some_authority/foo.bar" will be parsed
// into &Target{Scheme: "dns", Authority: "some_authority", Endpoint: "foo.bar"}
//
// If the target does not contain a scheme, we will apply the default scheme, and set the Target to
// be the full target string. e.g. "foo.bar" will be parsed into
// &Target{Scheme: resolver.GetDefaultScheme(), Endpoint: "foo.bar"}.
//
// If the parsed scheme is not registered (i.e. no corresponding resolver available to resolve the
// endpoint), we set the Scheme to be the default scheme, and set the Endpoint to be the full target
// string. e.g. target string "unknown_scheme://authority/endpoint" will be parsed into
// &Target{Scheme: resolver.GetDefaultScheme(), Endpoint: "unknown_scheme://authority/endpoint"}.
type Target struct {
	Scheme    string
	Authority string
	Endpoint  string
}

// Build includes additional information for the builder to create
// the resolver.
//go:generate go-option -type "Build"
type Build struct{}

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
	Build(target Target, cc ClientConn, opts ...BuildOption) (Resolver, error)

	// Scheme returns the scheme supported by this resolver.
	// Scheme is defined at https://github.com/grpc/grpc/blob/master/doc/naming.md.
	Scheme() string
}

// ResolveAddr includes additional information for ResolveAddr.
//go:generate go-option -type "ResolveAddr"
type ResolveAddr struct{}

// ResolveNow includes additional information for ResolveNow.
//go:generate go-option -type "ResolveNow"
type ResolveNow struct{}

// Resolver watches for the updates on the specified target.
// Updates include address updates and service config updates.
type Resolver interface {
	ResolveAddr(ctx context.Context, opts ...ResolveAddrOption) ([]Address, error)
	// ResolveNow will be called to try to resolve the target name
	// again. It's just a hint, resolver can ignore this if it's not necessary.
	//
	// It could be called multiple times concurrently.
	ResolveNow(ctx context.Context, opts ...ResolveNowOption)
	// Close closes the resolver.
	Close()
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
}
