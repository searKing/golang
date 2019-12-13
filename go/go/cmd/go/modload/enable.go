// Copyright (c) 2019 The searKing authors. All Rights Reserved.
//
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file in the root of the source
// tree. An additional intellectual property rights grant can be found
// in the file PATENTS.  All contributing project authors may
// be found in the AUTHORS file in the root of the source tree.

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
