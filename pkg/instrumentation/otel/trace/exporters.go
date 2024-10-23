// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package trace

import (
	"sort"
	"sync"

	"github.com/searKing/golang/pkg/instrumentation/otel/trace/driver"
)

var (
	driversMu sync.RWMutex
	drivers   = make(map[string]driver.ExporterURLOpener)
	// defaultScheme is the default scheme to use.
	defaultScheme = "passthrough"
)

// Register makes a driver available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(driver driver.ExporterURLOpener) {
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
func Get(scheme string) driver.ExporterURLOpener {
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
	drivers = make(map[string]driver.ExporterURLOpener)
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
