// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync_test

import (
	"bytes"
	"fmt"
	"time"

	sync_ "github.com/searKing/golang/go/exp/sync"
)

var bufFixedPool = sync_.NewFixedPool[*bytes.Buffer](func() *bytes.Buffer {
	// The Pool's New function should generally only return pointer
	// types, since a pointer can be put into the return interface
	// value without an allocation:
	fmt.Println("allocating new buffer")
	return new(bytes.Buffer)
}, 1)

var bufCachedPool = sync_.NewCachedPool[*bytes.Buffer](func() *bytes.Buffer {
	// The Pool's New function should generally only return pointer
	// types, since a pointer can be put into the return interface
	// value without an allocation:
	fmt.Println("allocating new buffer")
	return new(bytes.Buffer)
})
var bufTempPool = sync_.NewTempPool[*bytes.Buffer](nil)

// timeNow is a fake version of time.Now for tests.
func timeNow() time.Time {
	return time.Unix(1136214245, 0)
}

func Log(bufPool *sync_.FixedPool[*bytes.Buffer], key, val string) {
	be := bufPool.Get()
	{
		b := be.Value
		b.Reset()
		// Replace this with time.Now() in a real logger.
		b.WriteString(timeNow().UTC().Format(time.RFC3339))
		b.WriteByte(' ')
		b.WriteString(key)
		b.WriteByte('=')
		b.WriteString(val)
		fmt.Println(b.String())
	}
	bufPool.Put(be)
}
func ExampleNewFixedPool() {
	Log(bufFixedPool, "path", "/search?q=flowers")
	Log(bufFixedPool, "path", "/search?q=vegetables")
	// Output:
	// 2006-01-02T15:04:05Z path=/search?q=flowers
	// 2006-01-02T15:04:05Z path=/search?q=vegetables
}

func ExampleNewCachedPool() {
	Log(bufCachedPool, "path", "/search?q=flowers")
	Log(bufCachedPool, "path", "/search?q=vegetables")
	// Output:
	// allocating new buffer
	// 2006-01-02T15:04:05Z path=/search?q=flowers
	// 2006-01-02T15:04:05Z path=/search?q=vegetables
}

func ExampleNewTempPool() {
	bufTempPool.Emplace(new(bytes.Buffer))
	Log(bufTempPool, "path", "/search?q=flowers")
	Log(bufTempPool, "path", "/search?q=vegetables")
	// Output:
	// 2006-01-02T15:04:05Z path=/search?q=flowers
	// 2006-01-02T15:04:05Z path=/search?q=vegetables
}
