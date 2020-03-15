// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package modload

import (
	"fmt"
	"os"
)

func ModEnabled(dir string) bool {
	force, _ := mustUseModules()
	gomod := FindGoMod(dir)
	return gomod != "" || force
}

func mustUseModules() (mustUseModules bool, err error) {
	env := os.Getenv("GO111MODULE")
	switch env {
	default:
		return false, fmt.Errorf("go: unknown environment setting GO111MODULE=%s", env)
	case "auto", "":
		mustUseModules = false
	case "on":
		mustUseModules = true
	case "off":
		mustUseModules = false
		return
	}
	return mustUseModules, nil
}
