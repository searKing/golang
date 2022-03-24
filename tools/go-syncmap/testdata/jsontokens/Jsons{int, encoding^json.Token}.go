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

//go:generate go-syncmap -type "Jsons<int, encoding/json.Token>"
type Jsons sync.Map

const (
	_ = iota
	One
	Two
	Three
	AnotherOne = One // Duplicate; note that AnotherOne doesn't appear below.
)

func main() {
	var times Jsons
	times.Store(One, One)
	times.Store(Two, Two)
	times.Store(Three, Three)
	times.Store(AnotherOne, One)
	ck(times, One, One)
	ck(times, Two, Two)
	ck(times, Three, Three)
	ck(times, AnotherOne, One)
	ck(times, 127, 127)
}

func ck(jsons Jsons, num int, t json.Token) {
	val, loaded := jsons.Load(num)
	if num < One || num > Three {
		if loaded {
			panic(fmt.Sprintf("Jsons<int, encoding/json.Token>.go: %s", t))
		}
		return
	}
	if !loaded || val != t {
		panic(fmt.Sprintf("Jsons<int, encoding/json.Token>.go: %s", t))
	}
}
