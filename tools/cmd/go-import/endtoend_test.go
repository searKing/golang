// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// go command is not available on android

// +build !android

package main

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	io_ "github.com/searKing/golang/go/io"
)

// This file contains a test that compiles and runs each program in testdata
// after generating the string method for its type. The rule is that for testdata/x.go
// we run stringer -type X and then compile and run the program. The resulting
// binary panics if the String method for X is not correct, including for error cases.

func TestEndToEnd(t *testing.T) {
	dir, gooption := buildOptions(t)
	defer os.RemoveAll(dir)
	// Read the testdata directory.
	walkDir(dir, gooption, "testdata", t)
}

func walkDir(dir, gooption, dirname string, t *testing.T) {
	// Generate, compile, and run the test programs.
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		t.Fatalf("read dir[%s] failed %s", dirname, err)
		return
	}
	for _, file := range files {
		if file.IsDir() {
			goimportCompileAndRun(t, dir, gooption, filepath.Join(dirname, file.Name()))
			continue
		}
	}
}

// buildOptions creates a temporary directory and installs go-import there.
func buildOptions(t *testing.T) (dir string, gooptions string) {
	t.Helper()
	dir, err := ioutil.TempDir("", "go-import")
	if err != nil {
		t.Fatal(err)
	}
	gooptions = filepath.Join(dir, "go-import.exe")
	err = run("go", "build", "-o", gooptions)
	if err != nil {
		t.Fatalf("building go-import: %s", err)
	}
	return dir, gooptions
}

// goimportCompileAndRun runs stringer for the named file and compiles and
// runs the target binary in directory dir. That binary will panic if the String method is incorrect.
func goimportCompileAndRun(t *testing.T, dir, goimport, fileName string) {
	t.Helper()
	t.Logf("run: %s\n", fileName)
	target := filepath.Join(dir, fileName)
	source := fileName
	err := os.MkdirAll(filepath.Dir(target), os.ModePerm)
	if err != nil {
		t.Fatalf("mkdir temporary directory: %s", err)
	}
	err = io_.CopyDir(source, target, io_.Content)
	if err != nil {
		t.Fatalf("copying file to temporary directory: %s", err)
	}

	// Run goimport in the temporary directory.
	err = run(goimport, target)
	if err != nil {
		t.Fatal(err)
	}
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
