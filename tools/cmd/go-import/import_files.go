// Code generated for package main by go-bindata DO NOT EDIT. (@generated)
// sources:
// tmpl/import.tmpl
package main

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
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _tmplImportTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x6c\x92\xc1\x6e\xb3\x30\x10\x84\xef\x3c\xc5\x88\x63\x7e\xfd\xe6\x01\x7a\x6a\x9b\xaa\xca\xa1\x49\x0f\xdc\x91\x13\x6f\x8c\x55\xb0\x23\xdb\x51\x55\x59\xfb\xee\x15\xe0\xa4\x21\xe4\x06\xec\xe7\x99\xf1\x2c\x29\xfd\x47\xb5\x2a\xea\xd6\x04\x98\x80\xd8\x12\x4c\x7f\x72\x3e\x22\x52\x7f\xea\x64\x24\x51\x6c\x22\x14\x45\xf2\xbd\xb1\x14\xd0\xba\xef\x11\x3b\x3a\xdf\xcb\x18\x8d\xd5\xf9\x44\x80\xf4\x04\x4f\x56\x91\x27\x25\x8a\x55\xc5\x5c\xa4\x04\x45\x47\x63\x09\x65\x4b\x52\x91\x6f\x02\x1d\xa2\x71\xb6\x04\x73\x55\xe1\xd5\x29\x82\x26\x4b\x5e\x46\x52\xd8\xff\xa0\x4c\x49\xbc\xbb\xcd\x28\x59\x3b\xd7\x6d\x65\x4f\xcc\x29\x79\x69\x35\x61\x36\x7a\xf6\x3a\x80\x19\x29\x89\x81\x20\xab\x98\xcb\x27\xac\x77\xd8\xee\x6a\xbc\xad\x37\xb5\x48\x09\x64\x15\x98\x8b\xdb\x28\x27\x79\xf8\x92\x9a\x6e\xb3\x14\x00\x90\xbf\x0f\x82\x1f\x4e\x9d\x3b\x9a\xcc\x8b\xc7\x2a\xd3\xb5\x9b\xce\x84\x78\x55\x48\x09\x39\xe8\x14\xf3\x53\xc6\x36\x5c\x86\x19\x30\x47\x08\x30\xe7\x9e\x01\x34\xe3\xa5\x99\xcb\x3f\x9f\x8c\xe6\xb7\xc7\xfe\xfb\xb3\xe9\x54\x13\xa5\x5e\x44\x18\x1c\x5e\x86\x69\x2d\xf5\xcc\xbd\xaa\xfe\x8d\xa7\xae\x75\x5e\xa9\xbb\x1a\xb3\xd2\xf4\x3c\xb7\xbf\xfc\x18\xcb\x8d\x8a\x25\x72\x9f\xf1\x01\x72\xbf\x0c\x81\x25\x33\xab\x7a\xd4\xf8\x0d\x00\x00\xff\xff\x82\x57\x08\xc5\xb9\x02\x00\x00")

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

	info := bindataFileInfo{name: "tmpl/import.tmpl", size: 697, mode: os.FileMode(420), modTime: time.Unix(1576315168, 0)}
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
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
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
