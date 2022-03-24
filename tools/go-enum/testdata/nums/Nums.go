// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Enumeration with an offset.
// Also includes a duplicate.

package main

import (
	"encoding/json"
	"fmt"
)

//go:generate go-enum -type "Nums"
type Nums int

const (
	_ Nums = iota
	One
	Two
	Three
	AnotherOne = One // Duplicate; note that AnotherOne doesn't appear below.
)

func main() {
	ckRegistered(One, true)
	ckRegistered(Two, true)
	ckRegistered(Three, true)
	ckRegistered(AnotherOne, true)
	ckRegistered(Nums(127), false)

	ckString(One, "One")
	ckString(Two, "Two")
	ckString(Three, "Three")
	ckString(AnotherOne, "One")
	ckString(Nums(127), "Nums(127)")

	ckJson(One, `"One"`)
	ckJson(Two, `"Two"`)
	ckJson(Three, `"Three"`)
	ckJson(AnotherOne, `"One"`)
	ckJson(Nums(127), `"Nums(127)"`)

	//ckYamlMarshal(One, "One\n")
	//ckYamlMarshal(Two, "Two\n")
	//ckYamlMarshal(Three, "Three\n")
	//ckYamlMarshal(AnotherOne, "One\n")
	//ckYamlMarshal(Nums(127), "Nums(127)\n")

	ckTextMarshal(One, "One")
	ckTextMarshal(Two, "Two")
	ckTextMarshal(Three, "Three")
	ckTextMarshal(AnotherOne, "One")
	ckTextMarshal(Nums(127), "Nums(127)")

	ckSqlScan("One", One, false)
	ckSqlScan("Two", Two, false)
	ckSqlScan("Three", Three, false)
	ckSqlScan("One", AnotherOne, false)
	ckSqlScan("Nums(127)", Nums(127), true)

	ckSqlValue(One, "One")
	ckSqlValue(Two, "Two")
	ckSqlValue(Three, "Three")
	ckSqlValue(AnotherOne, "One")
	ckSqlValue(Nums(127), "Nums(127)")
}

func ckRegistered(nums Nums, registered bool) {
	if nums.Registered() == registered {
		return
	}
	panic(fmt.Sprintf("Nums.go: got %s, expect in %v", NumsValues(), nums))
}

func ckString(nums Nums, str string) {
	if nums.String() == str {
		return
	}
	panic(fmt.Sprintf("Nums.go: got %s, expect %s", nums.String(), str))
}

func ckJson(nums Nums, str string) {
	bytes, err := json.Marshal(nums)
	if err != nil {
		panic(fmt.Sprintf("Nums.go: json.Marshal failed: %s", err))
	}
	if string(bytes) == str {
		return
	}
	panic(fmt.Sprintf("Nums.go: got %s, expect %s", string(bytes), str))
}

//func ckYamlMarshal(nums Nums, str string) {
//	bytes, err := yaml.Marshal(nums)
//	if err != nil {
//		panic(fmt.Sprintf("Nums.go: yaml.Marshal failed: %s", err))
//	}
//	if string(bytes) == str {
//		return
//	}
//	panic(fmt.Sprintf("Nums.go: got %s, expect %s", string(bytes), str))
//}

func ckTextMarshal(nums Nums, str string) {
	bytes, err := nums.MarshalText()
	if err != nil {
		panic(fmt.Sprintf("Nums.go: MarshalText failed: %s", err))
	}
	if string(bytes) == str {
		return
	}
	panic(fmt.Sprintf("Nums.go: got %s, expect %s", string(bytes), str))
}

func ckSqlScan(str string, nums Nums, errorExpected bool) {
	var gotNum Nums
	err := (&gotNum).Scan(str)
	if errorExpected {
		if err != nil {
			return
		}
		panic(fmt.Sprintf("Nums.go: sql.Scan expect err, but success"))
	}

	if err != nil && !errorExpected {
		panic(fmt.Sprintf("Nums.go: sql.Scan failed: %s", err))
	}
	if gotNum == nums {
		return
	}
	panic(fmt.Sprintf("Nums.go: got %s, expect %s", gotNum, nums))
}

func ckSqlValue(nums Nums, str string) {
	val, err := nums.Value()
	if err != nil {
		panic(fmt.Sprintf("Nums.go: driver.Value failed: %s", err))
	}

	if val.(string) == str {
		return
	}
	panic(fmt.Sprintf("Nums.go: got %s, expect %s", val.(string), str))
}
