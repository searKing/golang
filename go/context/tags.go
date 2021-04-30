// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package context

import (
	"context"
)

var (
	// NopTags is a trivial, minimum overhead implementation of Tags for which all operations are no-ops.
	NopTags = &nopTags{}
)

// Tags is the interface used for storing request tags between Context calls.
// The default implementation is *not* thread safe, and should be handled only in the context of the request.
type Tags interface {
	// Set sets the given key in the metadata tags.
	Set(key string, value interface{})
	// Get gets if the metadata tags got by the given key exists.
	Get(key string) (interface{}, bool)
	// Del deletes the values associated with key.
	Del(key string)
	// Values returns a map of key to values.
	// Do not modify the underlying map, please use Set instead.
	Values() map[string]interface{}
}

// ExtractTags returns a pre-existing Tags object in the Context.
// If the context wasn't set in a tag interceptor, a no-op Tag storage is returned that will *not* be propagated in context.
func ExtractTags(ctx context.Context, key interface{}) (tags Tags, has bool) {
	t, ok := ctx.Value(key).(Tags)
	if !ok {
		return NopTags, false
	}

	return t, true
}

// ExtractOrCreateTags extracts or create tags from context by key
func ExtractOrCreateTags(ctx context.Context, key interface{}, options ...MapTagsOption) (
	ctx_ context.Context, stags Tags) {
	tags, has := ExtractTags(ctx, key)
	if has {
		return ctx, tags
	}
	tags = NewMapTags(options...)
	return WithTags(ctx, key, tags), tags
}

func WithTags(ctx context.Context, key interface{}, tags Tags) context.Context {
	return context.WithValue(ctx, key, tags)
}

func NewMapTags(options ...MapTagsOption) Tags {
	t := &mapTags{values: make(map[string]interface{})}
	t.ApplyOptions(options...)
	return t
}
