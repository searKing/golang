// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
)

//go:generate go-enum -type Transform  -trimprefix=Transform -transform=lower
type Transform int

const (
	TransformOne Transform = iota
	TransformTwo
	TransformThree
	TransformAnotherOne = TransformOne
)

func main() {
	ckRegistered(TransformOne, true)
	ckRegistered(TransformTwo, true)
	ckRegistered(TransformThree, true)
	ckRegistered(TransformAnotherOne, true)
	ckRegistered(Transform(127), false)

	ckString(TransformOne, "one")
	ckString(TransformTwo, "two")
	ckString(TransformThree, "three")
	ckString(TransformAnotherOne, "one")
	ckString(Transform(127), "Transform(127)")

	ckJson(TransformOne, `"one"`)
	ckJson(TransformTwo, `"two"`)
	ckJson(TransformThree, `"three"`)
	ckJson(TransformAnotherOne, `"one"`)
	ckJson(Transform(127), `"Transform(127)"`)

	//ckYamlMarshal(TransformOne, "One\n")
	//ckYamlMarshal(TransformTwo, "Two\n")
	//ckYamlMarshal(TransformThree, "Three\n")
	//ckYamlMarshal(TransformAnotherOne, "One\n")
	//ckYamlMarshal(Transform(127), "Transform(127)\n")

	ckTextMarshal(TransformOne, "one")
	ckTextMarshal(TransformTwo, "two")
	ckTextMarshal(TransformThree, "three")
	ckTextMarshal(TransformAnotherOne, "one")
	ckTextMarshal(Transform(127), "Transform(127)")

	ckSqlScan("one", TransformOne, false)
	ckSqlScan("two", TransformTwo, false)
	ckSqlScan("three", TransformThree, false)
	ckSqlScan("one", TransformAnotherOne, false)
	ckSqlScan("Transform(127)", Transform(127), true)

	ckSqlValue(TransformOne, "one")
	ckSqlValue(TransformTwo, "two")
	ckSqlValue(TransformThree, "three")
	ckSqlValue(TransformAnotherOne, "one")
	ckSqlValue(Transform(127), "Transform(127)")
}

func ckRegistered(prefix Transform, registered bool) {
	if prefix.Registered() == registered {
		return
	}
	panic(fmt.Sprintf("Transform.go: got %s, expect in %v", prefix.String(), prefix))
}

func ckString(prefix Transform, str string) {
	if prefix.String() == str {
		return
	}
	panic(fmt.Sprintf("Transform.go: got %s, expect %s", prefix.String(), str))
}

func ckJson(prefix Transform, str string) {
	bytes, err := json.Marshal(prefix)
	if err != nil {
		panic(fmt.Sprintf("Transform.go: json.Marshal failed: %s", err))
	}
	if string(bytes) == str {
		return
	}
	panic(fmt.Sprintf("Transform.go: got %s, expect %s", string(bytes), str))
}

//func ckYamlMarshal(prefix Transform, str string) {
//	bytes, err := yaml.Marshal(prefix)
//	if err != nil {
//		panic(fmt.Sprintf("Transform.go: yaml.Marshal failed: %s", err))
//	}
//	if string(bytes) == str {
//		return
//	}
//	panic(fmt.Sprintf("Transform.go: got %s, expect %s", string(bytes), str))
//}

func ckTextMarshal(prefix Transform, str string) {
	bytes, err := prefix.MarshalText()
	if err != nil {
		panic(fmt.Sprintf("Transform.go: MarshalText failed: %s", err))
	}
	if string(bytes) == str {
		return
	}
	panic(fmt.Sprintf("Transform.go: got %s, expect %s", string(bytes), str))
}

func ckSqlScan(str string, prefix Transform, errorExpected bool) {
	var gotPrefix Transform
	err := (&gotPrefix).Scan(str)
	if errorExpected {
		if err != nil {
			return
		}
		panic(fmt.Sprintf("Transform.go: sql.Scan expect err, but success"))
	}

	if err != nil && !errorExpected {
		panic(fmt.Sprintf("Transform.go: sql.Scan failed: %s", err))
	}
	if gotPrefix == prefix {
		return
	}
	panic(fmt.Sprintf("Transform.go: got %s, expect %s", gotPrefix, prefix))
}

func ckSqlValue(prefix Transform, str string) {
	val, err := prefix.Value()
	if err != nil {
		panic(fmt.Sprintf("Transform.go: driver.Value failed: %s", err))
	}

	if val.(string) == str {
		return
	}
	panic(fmt.Sprintf("Transform.go: got %s, expect %s", val.(string), str))
}
