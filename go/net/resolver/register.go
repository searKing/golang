// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"sort"
	"sync"
)

var (
	resolversMu sync.RWMutex
	resolvers   = make(map[string]Builder)
	// defaultScheme is the default scheme to use.
	defaultScheme = "passthrough"
)

// Register makes a database driver available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(driver Builder) {
	resolversMu.Lock()
	defer resolversMu.Unlock()
	if driver == nil {
		panic("resolver: Register driver is nil")
	}
	if _, dup := resolvers[driver.Scheme()]; dup {
		panic("resolver: Register called twice for driver " + driver.Scheme())
	}
	resolvers[driver.Scheme()] = driver
}

// Get returns the resolver builder registered with the given scheme.
//
// If no builder is register with the scheme, nil will be returned.
func Get(scheme string) Builder {
	resolversMu.Lock()
	defer resolversMu.Unlock()
	if b, ok := resolvers[scheme]; ok {
		return b
	}
	return nil
}

// SetDefaultScheme sets the default scheme that will be used. The default
// default scheme is "passthrough".
//
// NOTE: this function must only be called during initialization time (i.e. in
// an init() function), and is not thread-safe. The scheme set last overrides
// previously set values.
func SetDefaultScheme(scheme string) {
	defaultScheme = scheme
}

// GetDefaultScheme gets the default scheme that will be used.
func GetDefaultScheme() string {
	return defaultScheme
}

func unregisterAllDrivers() {
	resolversMu.Lock()
	defer resolversMu.Unlock()
	// For tests.
	resolvers = make(map[string]Builder)
}

// Resolvers returns a sorted list of the names of the registered resolvers.
func Resolvers() []string {
	resolversMu.RLock()
	defer resolversMu.RUnlock()
	list := make([]string, 0, len(resolvers))
	for name := range resolvers {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}
