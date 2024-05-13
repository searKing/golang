// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package context

import "net/textproto"

//go:generate go-option -type=mapTags
type mapTags struct {
	isMimeKey bool // represents a MIME-style key mapping

	values map[string]any
}

func (t *mapTags) Set(key string, value any) {
	if t.isMimeKey {
		key = textproto.CanonicalMIMEHeaderKey(key)
	}
	t.values[key] = value
}

func (t *mapTags) Get(key string) (any, bool) {
	if t.isMimeKey {
		key = textproto.CanonicalMIMEHeaderKey(key)
	}
	val, ok := t.values[key]
	return val, ok
}

// Del deletes the values associated with key.
func (t *mapTags) Del(key string) {
	if t.isMimeKey {
		key = textproto.CanonicalMIMEHeaderKey(key)
	}
	delete(t.values, key)
}

func (t *mapTags) Values() map[string]any {
	return t.values
}

// WithMapTagsMimeKey represents a MIME-style key mapping
func WithMapTagsMimeKey() MapTagsOption {
	return MapTagsOptionFunc(func(t *mapTags) {
		t.isMimeKey = true
	})
}
