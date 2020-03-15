// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package io

import "io"

type Watcher interface {
	Watch(p []byte, n int, err error) (int, error)
}
type WatcherFunc func(p []byte, n int, err error) (int, error)

func (f WatcherFunc) Watch(p []byte, n int, err error) (int, error) {
	return f(p, n, err)
}

type watchReader struct {
	source  io.Reader
	watcher Watcher
}

func (r *watchReader) Read(p []byte) (int, error) {
	var dummy io.Reader
	if r == nil || r.source == nil {
		dummy = EOFReader()
	} else {
		dummy = r.source
	}
	n, err := dummy.Read(p)
	if r.watcher == nil {
		return n, err
	}
	return r.watcher.Watch(p, n, err)
}

// WatchReader returns a Reader that's watch the Read state of
// the provided input reader.
func WatchReader(r io.Reader, watcher Watcher) io.Reader {
	return &watchReader{source: r, watcher: watcher}
}
