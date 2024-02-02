// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gopathload_test

import (
	"go/build"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/searKing/golang/go/go/cmd/go/gopathload"
)

func getFakeFS(t *testing.T, files ...string) string {
	tempDir, err := os.MkdirTemp("", "go_test")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}

	for _, f := range files {
		if err := os.MkdirAll(filepath.Join(tempDir, filepath.Dir(f)), 0770); err != nil {
			t.Fatalf("Failed to create directory structure for %v", f)
		}
		if err := os.WriteFile(filepath.Join(tempDir, f), nil, 0660); err != nil {
			t.Fatalf("Failed to create dummy file")
		}
	}
	return tempDir
}

func TestImportPackage(t *testing.T) {
	files := []string{
		"src/pkg1/sub/dummy.go",
		"src/pkg2/sub/dummy.go",
	}
	goPath1 := getFakeFS(t, files...)
	goPath2 := getFakeFS(t, files...)
	defer os.RemoveAll(goPath1)
	defer os.RemoveAll(goPath2)

	defer os.Setenv("GOPATH", os.Getenv("GOPATH"))
	os.Setenv("GOPATH", strings.Join([]string{goPath1, goPath2}, string(filepath.ListSeparator)))

	defer func(gopath string) {
		build.Default.GOPATH = gopath
	}(build.Default.GOPATH)
	build.Default.GOPATH = os.Getenv("GOPATH")

	tests := []struct {
		importPath string
		want       string
		wantErr    bool
	}{
		{
			importPath: "pkg1/sub/dummy.go",
			want:       filepath.Join(goPath1, "src"),
		},
		{
			importPath: "pkg2/sub/dummy.go",
			want:       filepath.Join(goPath1, "src"),
		},
		{
			importPath: "pkg3/sub/dummy.go",
			wantErr:    true,
		},
	}

	for n, tt := range tests {
		srcdir, err := gopathload.ImportPackage(tt.importPath)
		gotErr := err != nil
		if tt.wantErr != gotErr {
			t.Errorf("#%d: importPath %v; err %v; got %v; want %v", n, tt.importPath, err, gotErr, tt.wantErr)
			continue
		}
		if tt.want != srcdir {
			t.Errorf("#%d: importPath %v; got %v; want %v", n, tt.importPath, srcdir, tt.want)
		}
	}
}

func TestImportFile(t *testing.T) {
	files := []string{
		"src/pkg1/sub/dummy.go",
		"src/pkg2/sub/dummy.go",
	}
	goPath1 := getFakeFS(t, files...)
	goPath2 := getFakeFS(t, files...)
	defer os.RemoveAll(goPath1)
	defer os.RemoveAll(goPath2)

	defer os.Setenv("GOPATH", os.Getenv("GOPATH"))
	os.Setenv("GOPATH", strings.Join([]string{goPath1, goPath2}, string(filepath.ListSeparator)))

	defer func(gopath string) {
		build.Default.GOPATH = gopath
	}(build.Default.GOPATH)
	build.Default.GOPATH = os.Getenv("GOPATH")

	tests := []struct {
		filename       string
		wantSrcDir     string
		wantImportPath string
		wantErr        bool
	}{
		{
			filename:       filepath.Join(goPath1, "src", "pkg1/sub/dummy.go"),
			wantSrcDir:     filepath.Join(goPath1, "src"),
			wantImportPath: path.Dir("pkg1/sub/dummy.go"),
		},
		{
			filename:       filepath.Join(goPath1, "src", "pkg2/sub/dummy.go"),
			wantSrcDir:     filepath.Join(goPath1, "src"),
			wantImportPath: path.Dir("pkg2/sub/dummy.go"),
		},
		{
			filename: "pkg3/sub/dummy.go",
			wantErr:  true,
		},
	}

	for n, tt := range tests {
		srcdir, importPath, err := gopathload.ImportFile(tt.filename)
		gotErr := err != nil
		if tt.wantErr != gotErr {
			t.Errorf("#%d: filename %v; err %v; got %v; want %v", n, tt.filename, err, gotErr, tt.wantErr)
			continue
		}
		if tt.wantSrcDir != srcdir {
			t.Errorf("#%d: filename %v; got %v; want %v", n, tt.filename, srcdir, tt.wantSrcDir)
		}
		if tt.wantImportPath != importPath {
			t.Errorf("#%d: filename %v; got %v; want %v", n, tt.filename, importPath, tt.wantImportPath)
		}
	}
}
