// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// go command is not available on android

//go:build !android

package main

import (
	"go/build"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// This file contains a test that compiles and runs each program in testdata
// after generating the string method for its type. The rule is that for testdata/x.go
// we run stringer -type X and then compile and run the program. The resulting
// binary panics if the String method for X is not correct, including for error cases.

func TestEndToEnd(t *testing.T) {
	dir, gounions := buildUnions(t)
	defer os.RemoveAll(dir)
	// Read the testdata directory.
	walkDir(dir, gounions, "testdata", t)
}

func walkDir(dir, gounions, dirname string, t *testing.T) {
	// Generate, compile, and run the test programs.
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		t.Fatalf("read dir[%s] failed %s", dirname, err)
		return
	}
	for _, file := range files {
		name := file.Name()
		if file.IsDir() {
			walkDir(dir, gounions, filepath.Join(dirname, name), t)
			continue
		}
		if file.Mode().IsRegular() {
			if !strings.HasSuffix(name, ".go") {
				t.Errorf("%s is not a Go file", name)
				continue
			}
			if strings.HasPrefix(name, ".") {
				// This file is ignored by the build tool since it's name starts with '.'.
				continue
			}
			if strings.HasPrefix(name, "tag_") || strings.HasPrefix(name, "vary_") {
				// This file is used for tag processing in TestTags or TestConstValueChange, below.
				continue
			}
			if strings.HasSuffix(name, "_unions.go") {
				// This file is ignored by the build tool since it's name ends with '_unions.go'.
				continue
			}
			if name == "cgo.go" && !build.Default.CgoEnabled {
				t.Logf("cgo is not enabled for %s", name)
				continue
			}
			// Names are known to be ASCII and long enough.
			typeName := castFileNameToTypeName(name[:len(name)-len(".go")])
			gounionsCompileAndRun(t, dir, gounions, typeName, filepath.Join(dirname, name))
		}
	}
}

// buildUnions creates a temporary directory and installs go-union there.
func buildUnions(t *testing.T) (dir string, gounions string) {
	t.Helper()
	dir, err := ioutil.TempDir("", "go-union")
	if err != nil {
		t.Fatal(err)
	}
	gounions = filepath.Join(dir, "go-union.exe")
	err = run("go", "build", "-o", gounions)
	if err != nil {
		t.Fatalf("building go-union: %s", err)
	}
	return dir, gounions
}

// gounionsCompileAndRun runs stringer for the named file and compiles and
// runs the target binary in directory dir. That binary will panic if the String method is incorrect.
func gounionsCompileAndRun(t *testing.T, dir, gounions, typeName, fileName string) {
	t.Helper()
	t.Logf("run: %s %s\n", fileName, typeName)
	source := filepath.Join(dir, fileName)
	target := fileName
	err := os.MkdirAll(filepath.Dir(source), os.ModePerm)
	if err != nil {
		t.Fatalf("mkdir temporary directory: %s", err)
	}
	err = copy(source, target)
	if err != nil {
		t.Fatalf("copying file to temporary directory: %s", err)
	}

	unionsSource := filepath.Join(filepath.Dir(source), castTypeNameToFileName(typeName+"_unions.go"))
	// Run gounions in temporary directory.
	err = run(gounions, "-type", typeName, "-output", unionsSource, source)
	if err != nil {
		t.Fatal(err)
	}
	// Run the binary in the temporary directory.
	err = run("go", "run", unionsSource, source)
	if err != nil {
		t.Fatal(err)
	}
}

// castFileNameToTypeName replace "{" "}" "^" "@" with "<" ">" "/" "*"
// to fulfill windows os's constraint
// https://docs.microsoft.com/zh-cn/windows/win32/fileio/naming-a-file
func castFileNameToTypeName(name string) string {
	name = strings.ReplaceAll(name, "{", "<")
	name = strings.ReplaceAll(name, "}", ">")
	name = strings.ReplaceAll(name, "^", "/")
	name = strings.ReplaceAll(name, "@", "*")
	return name
}

// castFileNameToTypeName replace "<" ">" "/" "*" with "{" "}" "^" "@"
// to fulfill windows os's constraint
// https://docs.microsoft.com/zh-cn/windows/win32/fileio/naming-a-file
func castTypeNameToFileName(name string) string {
	name = strings.ReplaceAll(name, "<", "{")
	name = strings.ReplaceAll(name, ">", "}")
	name = strings.ReplaceAll(name, "/", "^")
	name = strings.ReplaceAll(name, "*", "@")

	return name
}

// copy copies the from file to the to file.
func copy(to, from string) error {
	toFd, err := os.Create(to)
	if err != nil {
		return err
	}
	defer toFd.Close()
	fromFd, err := os.Open(from)
	if err != nil {
		return err
	}
	defer fromFd.Close()
	_, err = io.Copy(toFd, fromFd)
	return err
}

// run runs a single command and returns an error if it does not succeed.
// os/exec should have this function, to be honest.
func run(name string, arg ...string) error {
	return runInDir(".", name, arg...)
}

// runInDir runs a single command in directory dir and returns an error if
// it does not succeed.
func runInDir(dir, name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), "GO111MODULE=auto")
	return cmd.Run()
}
