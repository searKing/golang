// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated for package main by go-bindata DO NOT EDIT. (@generated)
// sources:
// tmpl/import.tmpl
package goimport

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() any {
	return nil
}

var _tmplImportTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x6c\x92\x41\x6f\xa3\x30\x10\x85\xef\xfc\x8a\x27\x8e\x59\xad\xf9\x01\x7b\xda\xdd\x54\x55\x0e\x4d\x7a\xe0\x8e\x9c\x78\x02\x56\xc1\x8e\x6c\x47\x55\x35\x9a\xff\x5e\x41\x9c\x34\x09\x1c\x61\xbe\x79\xef\xcd\x03\xe6\xdf\xa8\x56\x45\xdd\xd9\x08\x1b\x91\x3a\x82\x1d\x4e\x3e\x24\x24\x1a\x4e\xbd\x4e\xa4\x8a\x4d\x82\xa1\x44\x61\xb0\x8e\x22\x3a\xff\x39\x61\x47\x1f\x06\x9d\x92\x75\x6d\xde\x88\xd0\x81\x10\xc8\x19\x0a\x64\x54\xb1\xaa\x44\x0a\x66\x18\x3a\x5a\x47\x28\x3b\xd2\x86\x42\x13\xe9\x90\xac\x77\x25\x44\xaa\x0a\xff\xbd\x21\xb4\xe4\x28\xe8\x44\x06\xfb\x2f\x94\xcc\xea\xd5\x6f\x26\xc9\xda\xfb\x7e\xab\x07\x12\x61\x0e\xda\xb5\x84\x87\xd1\xdf\xd0\x46\x88\x80\x59\x8d\x04\x39\x23\x52\xfe\xc1\x7a\x87\xed\xae\xc6\xcb\x7a\x53\x2b\x66\x90\x33\x10\x29\xee\xa3\x9c\xf4\xe1\x43\xb7\x74\x9f\xa5\x00\x80\xfc\x7e\x14\x7c\xf3\xe6\xdc\xd3\xc5\xbc\x58\x56\xb9\x9c\xdd\xf4\x36\xa6\x9b\x02\x33\x72\xd0\x4b\xcc\x77\x9d\xba\x78\x1d\x66\xc0\x1e\xa1\x20\x92\x7b\x06\xd0\x4c\x47\x8b\x94\x3f\x3e\x19\xcd\x4f\xcb\xfe\xfb\xb3\xed\x4d\x93\x74\x3b\xbb\x23\x9b\xfc\x1b\x81\x5a\xb7\xf7\xfe\x55\xf5\x6b\xda\x1b\x6f\xbc\xce\x6f\x5b\x53\x83\x4f\x6e\xd7\xff\x60\xfe\x01\xd5\x1c\x59\x88\xb4\x40\x3d\xd7\xaf\x30\x67\x1e\xca\x9d\x34\xbe\x03\x00\x00\xff\xff\x62\x32\xc5\x48\xab\x02\x00\x00")

func tmplImportTmplBytes() ([]byte, error) {
	return bindataRead(
		_tmplImportTmpl,
		"tmpl/import.tmpl",
	)
}

func tmplImportTmpl() (*asset, error) {
	bytes, err := tmplImportTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "tmpl/import.tmpl", size: 683, mode: os.FileMode(420), modTime: time.Unix(1576512469, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"tmpl/import.tmpl": tmplImportTmpl,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//
//	data/
//	  foo.txt
//	  img/
//	    a.png
//	    b.png
//
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"tmpl": &bintree{nil, map[string]*bintree{
		"import.tmpl": &bintree{tmplImportTmpl, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
