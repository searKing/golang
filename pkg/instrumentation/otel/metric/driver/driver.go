// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package driver

import (
	"context"
	"net/url"

	"go.opentelemetry.io/otel/sdk/metric"
)

// ReaderURLOpener is the interface that must be implemented by a metric reader
// driver.
// ReaderURLOpener represents types that can open metric readers based on a URL.
// The opener must not modify the URL argument. OpenReaderURL must be safe to
// call from multiple goroutines.
//
// This interface is generally implemented by types in driver packages.
type ReaderURLOpener interface {
	// OpenReaderURL creates a new reader for the given target.
	OpenReaderURL(ctx context.Context, u *url.URL) (metric.Reader, error)

	// Scheme returns the scheme supported by this metric reader.
	// Scheme is defined at https://github.com/grpc/grpc/blob/master/doc/naming.md.
	Scheme() string
}
