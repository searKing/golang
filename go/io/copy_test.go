// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package io_test

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	io_ "github.com/searKing/golang/go/io"
)

func TestCopy(t *testing.T) {
	copyWithFileRange := true
	copyWithFileClone := true
	doCopyTest(t, &copyWithFileRange, &copyWithFileClone)
}

func TestCopyWithoutRange(t *testing.T) {
	copyWithFileRange := false
	copyWithFileClone := false
	doCopyTest(t, &copyWithFileRange, &copyWithFileClone)
}

func TestCopyDir(t *testing.T) {
	srcDir, err := os.MkdirTemp("", "srcDir")
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
	populateSrcDir(t, srcDir, 3)

	dstDir, err := os.MkdirTemp("", "testdst")
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
	defer os.RemoveAll(dstDir)

	err = io_.CopyDir(srcDir, dstDir, io_.Content)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}

	err = filepath.Walk(srcDir, func(srcPath string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Rebase path
		relPath, err := filepath.Rel(srcDir, srcPath)
		if err != nil {
			t.Errorf("expect nil, got %v", err)
		}
		if relPath == "." {
			return nil
		}

		dstPath := filepath.Join(dstDir, relPath)
		if err != nil {
			t.Errorf("expect nil, got %v", err)
		}

		// If we add non-regular dirs and files to the test
		// then we need to add more checks here.
		_, err = os.Lstat(dstPath)
		if err != nil {
			t.Errorf("expect nil, got %v", err)
		}

		return nil
	})
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
}

func randomMode(baseMode int) os.FileMode {
	for i := 0; i < 7; i++ {
		baseMode = baseMode | (1&rand.Intn(2))<<uint(i)
	}
	return os.FileMode(baseMode)
}

func populateSrcDir(t *testing.T, srcDir string, remainingDepth int) {
	if remainingDepth == 0 {
		return
	}

	for i := 0; i < 10; i++ {
		dirName := filepath.Join(srcDir, fmt.Sprintf("srcdir-%d", i))
		// Owner all bits set
		err := os.Mkdir(dirName, randomMode(0700))
		if err != nil {
			t.Errorf("expect nil, got %v", err)
		}
		populateSrcDir(t, dirName, remainingDepth-1)
	}

	for i := 0; i < 10; i++ {
		fileName := filepath.Join(srcDir, fmt.Sprintf("srcfile-%d", i))
		// Owner read bit set
		err := os.WriteFile(fileName, []byte{}, randomMode(0400))
		if err != nil {
			t.Errorf("expect nil, got %v", err)
		}
	}
}

func doCopyTest(t *testing.T, copyWithFileRange, copyWithFileClone *bool) {
	dir, err := os.MkdirTemp("", "docker-copy-check")
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
	defer os.RemoveAll(dir)
	srcFilename := filepath.Join(dir, "srcFilename")
	dstFilename := filepath.Join(dir, "dstilename")

	r := rand.New(rand.NewSource(0))
	buf := make([]byte, 1024)
	_, err = r.Read(buf)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
	err = os.WriteFile(srcFilename, buf, 0777)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
	fileinfo, err := os.Stat(srcFilename)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}

	err = io_.CopyRegular(srcFilename, dstFilename, fileinfo)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
	readBuf, err := os.ReadFile(dstFilename)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}

	if !bytes.Equal(buf, readBuf) {
		t.Errorf("expect true, got %v", false)
	}
}

func TestCopyHardlink(t *testing.T) {
	var srcFile1FileInfo, srcFile2FileInfo, dstFile1FileInfo, dstFile2FileInfo os.FileInfo

	srcDir, err := os.MkdirTemp("", "srcDir")
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
	defer os.RemoveAll(srcDir)

	dstDir, err := os.MkdirTemp("", "dstDir")
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
	defer os.RemoveAll(dstDir)

	srcFile1 := filepath.Join(srcDir, "file1")
	srcFile2 := filepath.Join(srcDir, "file2")
	dstFile1 := filepath.Join(dstDir, "file1")
	dstFile2 := filepath.Join(dstDir, "file2")
	err = os.WriteFile(srcFile1, []byte{}, 0777)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
	err = os.Link(srcFile1, srcFile2)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}

	err = io_.CopyDir(srcDir, dstDir, io_.Content)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}

	srcFile1FileInfo, err = os.Stat(srcFile1)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
	dstFile1FileInfo, err = os.Stat(dstFile1)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
	srcFile2FileInfo, err = os.Stat(srcFile2)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
	dstFile2FileInfo, err = os.Stat(dstFile2)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}

	if srcFile1FileInfo.Name() != dstFile1FileInfo.Name() {
		t.Error("expect equal, got unequal")
	}
	if srcFile1FileInfo.Size() != dstFile1FileInfo.Size() {
		t.Error("expect equal, got unequal")
	}
	if srcFile1FileInfo.IsDir() != dstFile1FileInfo.IsDir() {
		t.Error("expect equal, got unequal")
	}
	if srcFile1FileInfo.Mode() != dstFile1FileInfo.Mode() {
		t.Error("expect equal, got unequal")
	}
	if srcFile2FileInfo.Name() != dstFile2FileInfo.Name() {
		t.Error("expect equal, got unequal")
	}
	if srcFile2FileInfo.Size() != dstFile2FileInfo.Size() {
		t.Error("expect equal, got unequal")
	}
	if srcFile2FileInfo.IsDir() != dstFile2FileInfo.IsDir() {
		t.Error("expect equal, got unequal")
	}
	if srcFile2FileInfo.Mode() != dstFile2FileInfo.Mode() {
		t.Error("expect equal, got unequal")
	}
}
