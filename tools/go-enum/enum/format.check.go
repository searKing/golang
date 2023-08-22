// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package enum

import (
	"fmt"
	"sync"
)

var checkImportPackages = []string{`fmt`}

// Arguments to format are:
//
//	[1]: type name
const stringNameToValueMethod = `
// Parse%[1]sString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func Parse%[1]sString(s string) (%[1]s, error) {
	if val, ok := _%[1]s_name_to_values[s]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%%s does not belong to %[1]s values", s)
}
`

// Arguments to format are:
//
//	[1]: type name
const stringValuesMethod = `
// %[1]sValues returns all values of the enum
func %[1]sValues() []%[1]s {
	return _%[1]s_values
}
`

// Arguments to format are:
//
//	[1]: type name
const stringBelongsMethodLoop = `
// IsA%[1]s returns "true" if the value is listed in the enum definition. "false" otherwise
func (i %[1]s) Registered() bool {
	for _, v := range _%[1]s_values {
		if i == v {
			return true
		}
	}
	return false
}
`

// Arguments to format are:
//
//	[1]: type name
const stringBelongsMethodSet = `
// IsA%[1]s returns "true" if the value is listed in the enum definition. "false" otherwise
func (i %[1]s) IsA%[1]s() bool {
	_, ok := _%[1]sMap[i] 
	return ok
}
`

var buildCheckOnce sync.Once

func (g *Generator) buildCheck(runs [][]Value, typeName string, runsThreshold int) {
	buildCheckOnce.Do(func() {
		// At this moment, either "g.declareIndexAndNameVars()" or "g.declareNameVars()" has been called

		// Print the slice of values
		g.Printf("\nvar _%[1]s_values = []%[1]s{", typeName)
		for _, values := range runs {
			for _, value := range values {
				g.Printf("\t%[1]s, ", value.valueInfo.str)
			}
		}
		g.Printf("}\n\n")

		// Print the map between name and value
		g.Printf("\nvar _%[1]s_name_to_values = map[string]%[1]s{\n", typeName)
		thereAreRuns := len(runs) > 1 && len(runs) <= runsThreshold
		var n int
		var runID string
		for i, values := range runs {
			if thereAreRuns {
				runID = "_" + fmt.Sprintf("%d", i)
				n = 0
			} else {
				runID = ""
			}

			for _, value := range values {
				g.Printf("\t_%s_name%s[%d:%d]: %s,\n", typeName, runID, n, n+len(value.nameInfo.trimmedName), &value)
				n += len(value.nameInfo.trimmedName)
			}
		}
		g.Printf("}\n\n")

		// Print the basic extra methods
		g.Printf(stringNameToValueMethod, typeName)
		g.Printf(stringValuesMethod, typeName)
		if len(runs) <= runsThreshold {
			g.Printf(stringBelongsMethodLoop, typeName)
		} else { // There is a map of values, the code is simpler then
			g.Printf(stringBelongsMethodSet, typeName)
		}
	})
}
