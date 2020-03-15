// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package context

type nopTags struct{}

func (t *nopTags) Set(key string, value interface{}) {
	return
}

func (t *nopTags) Get(key string) (interface{}, bool) {
	return nil, false
}

// Del deletes the values associated with key.
func (t *nopTags) Del(key string) {
	return
}
func (t *nopTags) Values() map[string]interface{} {
	return nil
}
