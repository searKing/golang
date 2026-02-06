// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package propagation_test

import (
	"slices"
	"testing"

	"github.com/searKing/golang/third_party/github.com/open-telemetry/opentelemetry-go/propagation"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

func TestMetadataCarrierGet(t *testing.T) {
	carrier := propagation.MetadataCarrier(metadata.New(map[string]string{
		"foo": "bar",
		"baz": "qux",
	}))

	assert.Equal(t, "bar", carrier.Get("foo"))
	assert.Equal(t, "qux", carrier.Get("baz"))
}

func TestMetadataCarrierSet(t *testing.T) {
	carrier := make(propagation.MetadataCarrier)
	carrier.Set("foo", "bar")
	carrier.Set("baz", "qux")

	assert.Equal(t, "bar", carrier["foo"][0])
	assert.Equal(t, "qux", carrier["baz"][0])
}

func TestMetadataCarrierKeys(t *testing.T) {
	carrier := propagation.MetadataCarrier(metadata.New(map[string]string{
		"foo": "bar",
		"baz": "qux",
	}))

	keys := carrier.Keys()
	slices.Sort(keys)
	assert.Equal(t, []string{"baz", "foo"}, keys)
}
