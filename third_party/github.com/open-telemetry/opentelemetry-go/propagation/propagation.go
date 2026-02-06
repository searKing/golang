// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package propagation

import (
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc/metadata"
)

// MetadataCarrier adapts metadata.MD to satisfy the TextMapCarrier and ValuesGetter interfaces.
type MetadataCarrier metadata.MD

// Compile time check that MetadataCarrier implements TextMapCarrier.
var _ propagation.TextMapCarrier = MetadataCarrier{}

// Compile time check that MetadataCarrier implements ValuesGetter.
var _ propagation.ValuesGetter = MetadataCarrier{}

// Get returns the first value associated with the passed key.
func (hc MetadataCarrier) Get(key string) string {
	v := metadata.MD(hc).Get(key)
	if len(v) == 0 {
		return ""
	}
	return v[0]
}

// Values returns all values associated with the passed key.
func (hc MetadataCarrier) Values(key string) []string {
	return metadata.MD(hc).Get(key)
}

// Set stores the key-value pair.
func (hc MetadataCarrier) Set(key, value string) {
	metadata.MD(hc).Set(key, value)
}

// Keys lists the keys stored in this carrier.
func (hc MetadataCarrier) Keys() []string {
	keys := make([]string, 0, len(hc))
	for k := range hc {
		keys = append(keys, k)
	}
	return keys
}
