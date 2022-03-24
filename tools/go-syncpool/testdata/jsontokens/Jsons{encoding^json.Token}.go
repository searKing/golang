// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Enumeration with an offset.
// Also includes a duplicate.

package main

import (
	"encoding/json"
	"fmt"
	"sync"
)

//go:generate go-syncpool -type "Jsons<encoding/json.Token>"
type Jsons sync.Pool

const (
	_ = iota
	One
	Two
	Three
	AnotherOne = One // Duplicate; note that AnotherOne doesn't appear below.
)

func main() {
	var times Jsons
	times.Put(One)
	ck(times, One)
	times.Put(Two)
	ck(times, Two)
	times.Put(Three)
	ck(times, Three)
	times.Put(AnotherOne)
	ck(times, One)
	times.Put(127)
	ck(times, 127)
}

func ck(jsons Jsons, t json.Token) {
	val := jsons.Get()
	if val != t {
		panic(fmt.Sprintf("Jsons<encoding/json.Token>.go: %s", t))
	}
}
