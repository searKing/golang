// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package trace

import (
	"context"
	"net/url"
	"sort"
	"sync"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// ExporterURLOpener represents types that can open metric exporters based on a URL.
// The opener must not modify the URL argument. OpenExporterURL must be safe to
// call from multiple goroutines.
//
// This interface is generally implemented by types in driver packages.
type ExporterURLOpener interface {
	// OpenExporterURL creates a new exporter for the given target.
	OpenExporterURL(ctx context.Context, u *url.URL) (sdktrace.SpanExporter, error)

	// Scheme returns the scheme supported by this exporter.
	// Scheme is defined at https://github.com/grpc/grpc/blob/master/doc/naming.md.
	Scheme() string
}

var (
	driversMu sync.RWMutex
	drivers   = make(map[string]ExporterURLOpener)
	// defaultScheme is the default scheme to use.
	defaultScheme = "passthrough"
)

// Register makes a driver available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(driver ExporterURLOpener) {
	driversMu.Lock()
	defer driversMu.Unlock()
	if driver == nil {
		panic("metric: Register driver is nil")
	}
	if _, dup := drivers[driver.Scheme()]; dup {
		panic("metric: Register called twice for driver " + driver.Scheme())
	}
	drivers[driver.Scheme()] = driver
}

// Get returns the metric url opener registered with the given scheme.
//
// If no driver is register with the scheme, nil will be returned.
func Get(scheme string) ExporterURLOpener {
	driversMu.Lock()
	defer driversMu.Unlock()
	if b, ok := drivers[scheme]; ok {
		return b
	}
	return nil
}

// SetDefaultScheme sets the default scheme that will be used. The default
// scheme is "passthrough".
//
// NOTE: this function must only be called during initialization time (i.e. in
// an init() function), and is not thread-safe. The scheme set last overrides
// previously set values.
func SetDefaultScheme(scheme string) {
	defaultScheme = scheme
}

// GetScheme gets the default scheme that will be used.
func GetScheme() string {
	return defaultScheme
}

func unregisterAllDrivers() {
	driversMu.Lock()
	defer driversMu.Unlock()
	// For tests.
	drivers = make(map[string]ExporterURLOpener)
}

// Drivers returns a sorted list of the names of the registered drivers.
func Drivers() []string {
	driversMu.RLock()
	defer driversMu.RUnlock()
	list := make([]string, 0, len(drivers))
	for name := range drivers {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}
