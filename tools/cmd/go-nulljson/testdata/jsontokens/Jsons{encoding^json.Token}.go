// Copyright 2019 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Enumeration with an offset.
// Also includes a duplicate.

package main

import (
	"encoding/json"
	"fmt"
)

//go:generate go-nulljson -type "Jsons<encoding/json.Token>"

const (
	_ = iota
	One
	Two
	Three
	AnotherOne = One // Duplicate; note that AnotherOne doesn't appear below.
)

func main() {
	var tokens Jsons
	err := tokens.Scan(One)
	ckError(err)
	ck(tokens, One)

	err = tokens.Scan(Two)
	ckError(err)
	ck(tokens, Two)

	err = tokens.Scan(Three)
	ckError(err)
	ck(tokens, Three)

	err = tokens.Scan(AnotherOne)
	ckError(err)
	ck(tokens, One)

	err = tokens.Scan(127)
	ckError(err)
	ck(tokens, 127)
}

func ck(jsons Jsons, t json.Token) {
	val := jsons.Data
	if int(val.(float64)) != t {
		panic(fmt.Sprintf("Jsons<encoding/json.Token>.go:got %s expected %s", val, t))
	}
}

func ckError(err error) {
	if err != nil {
		panic(fmt.Sprintf("Jsons<encoding/json.Token>.go: error happened %s", err))
	}
}
