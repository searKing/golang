// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package modload_test

import (
	"fmt"
	"os"

	"github.com/searKing/golang/go/go/cmd/go/modload"
)

func ExampleFindModuleName() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("getwd: %w", err))
	}
	importPath, err := modload.FindModuleName(modload.FindModuleRoot(cwd))
	if err != nil {
		panic(fmt.Errorf("find mod name: %w", err))
	}
	fmt.Print(importPath)
	// Output:
	// github.com/searKing/golang/go
}

func ExampleImportFile() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("getwd: %w", err))
	}
	//cwd = "/Users/chenhaixin/workspace/go/src/github.com/searKing/golang/go/go/cmd/go/modload/init.go"
	srcDir, importPath, err := modload.ImportFile(cwd)
	if err != nil {
		panic(fmt.Errorf("find import path: %w", err))
	}
	_ = srcDir
	fmt.Println(importPath)
	// Output:
	// github.com/searKing/golang/go/go/cmd/go/modload
}

func ExampleImportPackage() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("getwd: %w", err))
	}
	//cwd = "/Users/chenhaixin/workspace/go/src/github.com/searKing/golang/go/go/cmd/go/modload/init.go"
	srcDir, modname, err := modload.ImportPackage("github.com/searKing/golang/go/go/cmd/go/modload", cwd)
	if err != nil {
		panic(fmt.Errorf("find import path: %w", err))
	}
	_ = srcDir
	fmt.Println(modname)
	// Output:
	// github.com/searKing/golang/go
}
