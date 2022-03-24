// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Enumeration with an offset.
// Also includes a duplicate.

package main

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
)

//go:generate go-atomicvalue -type "Jsons<encoding/json.Token>"
type Jsons atomic.Value

const (
	_ = iota
	One
	Two
	Three
	AnotherOne = One // Duplicate; note that AnotherOne doesn't appear below.
)

func main() {
	var times Jsons
	times.Store(One)
	ck(times, One)
	times.Store(Two)
	ck(times, Two)
	times.Store(Three)
	ck(times, Three)
	times.Store(AnotherOne)
	ck(times, One)
	times.Store(127)
	ck(times, 127)
}

func ck(jsons Jsons, t json.Token) {
	val := jsons.Load()
	if val != t {
		panic(fmt.Sprintf("Jsons<encoding/json.Token>.go: %s", t))
	}
}
