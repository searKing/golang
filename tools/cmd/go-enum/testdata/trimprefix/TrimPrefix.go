// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
)

//go:generate go-enum -type TrimPrefix -trimprefix=TrimPrefix
type TrimPrefix int

const (
	TrimPrefixOne TrimPrefix = iota
	TrimPrefixTwo
	TrimPrefixThree
	TrimPrefixAnotherOne = TrimPrefixOne
)

func main() {
	ckRegistered(TrimPrefixOne, true)
	ckRegistered(TrimPrefixTwo, true)
	ckRegistered(TrimPrefixThree, true)
	ckRegistered(TrimPrefixAnotherOne, true)
	ckRegistered(TrimPrefix(127), false)

	ckString(TrimPrefixOne, "One")
	ckString(TrimPrefixTwo, "Two")
	ckString(TrimPrefixThree, "Three")
	ckString(TrimPrefixAnotherOne, "One")
	ckString(TrimPrefix(127), "TrimPrefix(127)")

	ckJson(TrimPrefixOne, `"One"`)
	ckJson(TrimPrefixTwo, `"Two"`)
	ckJson(TrimPrefixThree, `"Three"`)
	ckJson(TrimPrefixAnotherOne, `"One"`)
	ckJson(TrimPrefix(127), `"TrimPrefix(127)"`)

	//ckYamlMarshal(TrimPrefixOne, "One\n")
	//ckYamlMarshal(TrimPrefixTwo, "Two\n")
	//ckYamlMarshal(TrimPrefixThree, "Three\n")
	//ckYamlMarshal(TrimPrefixAnotherOne, "One\n")
	//ckYamlMarshal(TrimPrefix(127), "TrimPrefix(127)\n")

	ckTextMarshal(TrimPrefixOne, "One")
	ckTextMarshal(TrimPrefixTwo, "Two")
	ckTextMarshal(TrimPrefixThree, "Three")
	ckTextMarshal(TrimPrefixAnotherOne, "One")
	ckTextMarshal(TrimPrefix(127), "TrimPrefix(127)")

	ckSqlScan("One", TrimPrefixOne, false)
	ckSqlScan("Two", TrimPrefixTwo, false)
	ckSqlScan("Three", TrimPrefixThree, false)
	ckSqlScan("One", TrimPrefixAnotherOne, false)
	ckSqlScan("TrimPrefix(127)", TrimPrefix(127), true)

	ckSqlValue(TrimPrefixOne, "One")
	ckSqlValue(TrimPrefixTwo, "Two")
	ckSqlValue(TrimPrefixThree, "Three")
	ckSqlValue(TrimPrefixAnotherOne, "One")
	ckSqlValue(TrimPrefix(127), "TrimPrefix(127)")
}

func ckRegistered(prefix TrimPrefix, registered bool) {
	if prefix.Registered() == registered {
		return
	}
	panic(fmt.Sprintf("TrimPrefix.go: got %s, expect in %v", prefix.String(), prefix))
}

func ckString(prefix TrimPrefix, str string) {
	if prefix.String() == str {
		return
	}
	panic(fmt.Sprintf("TrimPrefix.go: got %s, expect %s", prefix.String(), str))
}

func ckJson(prefix TrimPrefix, str string) {
	bytes, err := json.Marshal(prefix)
	if err != nil {
		panic(fmt.Sprintf("TrimPrefix.go: json.Marshal failed: %s", err))
	}
	if string(bytes) == str {
		return
	}
	panic(fmt.Sprintf("TrimPrefix.go: got %s, expect %s", string(bytes), str))
}

//func ckYamlMarshal(prefix TrimPrefix, str string) {
//	bytes, err := yaml.Marshal(prefix)
//	if err != nil {
//		panic(fmt.Sprintf("TrimPrefix.go: yaml.Marshal failed: %s", err))
//	}
//	if string(bytes) == str {
//		return
//	}
//	panic(fmt.Sprintf("TrimPrefix.go: got %s, expect %s", string(bytes), str))
//}

func ckTextMarshal(prefix TrimPrefix, str string) {
	bytes, err := prefix.MarshalText()
	if err != nil {
		panic(fmt.Sprintf("TrimPrefix.go: MarshalText failed: %s", err))
	}
	if string(bytes) == str {
		return
	}
	panic(fmt.Sprintf("TrimPrefix.go: got %s, expect %s", string(bytes), str))
}

func ckSqlScan(str string, prefix TrimPrefix, errorExpected bool) {
	var gotPrefix TrimPrefix
	err := (&gotPrefix).Scan(str)
	if errorExpected {
		if err != nil {
			return
		}
		panic(fmt.Sprintf("TrimPrefix.go: sql.Scan expect err, but success"))
	}

	if err != nil && !errorExpected {
		panic(fmt.Sprintf("TrimPrefix.go: sql.Scan failed: %s", err))
	}
	if gotPrefix == prefix {
		return
	}
	panic(fmt.Sprintf("TrimPrefix.go: got %s, expect %s", gotPrefix, prefix))
}

func ckSqlValue(prefix TrimPrefix, str string) {
	val, err := prefix.Value()
	if err != nil {
		panic(fmt.Sprintf("TrimPrefix.go: driver.Value failed: %s", err))
	}

	if val.(string) == str {
		return
	}
	panic(fmt.Sprintf("TrimPrefix.go: got %s, expect %s", val.(string), str))
}
