// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package context

import "net/textproto"

//go:generate go-option -type=mapTags
type mapTags struct {
	isMimeKey bool // represents a MIME-style key mapping

	values map[string]interface{}
}

func (t *mapTags) Set(key string, value interface{}) {
	if t.isMimeKey {
		key = textproto.CanonicalMIMEHeaderKey(key)
	}
	t.values[key] = value
}

func (t *mapTags) Get(key string) (interface{}, bool) {
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

func (t *mapTags) Values() map[string]interface{} {
	return t.values
}

// represents a MIME-style key mapping
func WithTagsMimeKey() MapTagsOption {
	return MapTagsOptionFunc(func(t *mapTags) {
		t.isMimeKey = true
	})
}
