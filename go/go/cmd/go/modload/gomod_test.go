// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package modload_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/searKing/golang/go/go/cmd/go/modload"
)

func TestFindModuleName(t *testing.T) {
	tests := []struct {
		gomod   string
		want    string
		wantErr bool
	}{
		{gomod: "module example.com/a\n", want: "example.com/a"},
		{gomod: "\n\nmodule example.com/a\n", want: "example.com/a"},
		{gomod: "// comment\n\nmodule example.com/a // trailing comment\n\ngo 1.27\n", want: "example.com/a"},
		{gomod: "module \"example.com/a\"\n", want: "example.com/a"},
		{gomod: "module `example.com/a`\n", want: "example.com/a"},
		{gomod: "go 1.27\n", want: ""},
		{gomod: "module\n", wantErr: true},
		{gomod: "module example.com/a example.com/b\n", wantErr: true},
		{gomod: "module \"example.com/a\n", wantErr: true},
	}
	for _, tt := range tests {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(tt.gomod), 0o644); err != nil {
			t.Fatal(err)
		}
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("FindModuleName(%q) panicked: %v", tt.gomod, r)
				}
			}()
			got, err := modload.FindModuleName(dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindModuleName(%q) error = %v, wantErr %v", tt.gomod, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("FindModuleName(%q) = %q, want %q", tt.gomod, got, tt.want)
			}
		}()
	}
}

func TestFindGoMod(t *testing.T) {
	root := t.TempDir()
	sub := filepath.Join(root, "a", "b")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	gomod := filepath.Join(root, "go.mod")
	if err := os.WriteFile(gomod, []byte("module example.com/a\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, dir := range []string{root, sub} {
		if got := modload.FindGoMod(dir); got != gomod {
			t.Errorf("FindGoMod(%q) = %q, want %q", dir, got, gomod)
		}
	}
}

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
	//cwd = "//workspace/go/src/github.com/searKing/golang/go/go/cmd/go/modload/init.go"
	srcDir, importPath, err := modload.ImportFile(cwd)
	if err != nil {
		panic(fmt.Errorf("find import path: %w", err))
	}
	_ = srcDir
	fmt.Println(importPath)
	// github.com/searKing/golang/go/go/cmd/go/modload
}

func ExampleImportPackage() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("getwd: %w", err))
	}
	//cwd = "/workspace/go/src/github.com/searKing/golang/go/go/cmd/go/modload/init.go"
	srcDir, modname, err := modload.ImportPackage("github.com/searKing/golang/go/go/cmd/go/modload", cwd)
	if err != nil {
		panic(fmt.Errorf("find import path: %w", err))
	}
	_ = srcDir
	fmt.Println(modname)
	// github.com/searKing/golang/go
}
