package os_test

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	os_ "github.com/searKing/golang/go/os"
)

// tmpDir creates a temporary directory and returns its name.
func tmpFile(t *testing.T) string {
	tmp, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("temp file creation failed: %v", err)
	}
	defer os.Remove(tmp.Name())
	defer tmp.Close()
	return tmp.Name()
}

func TestCreateAll(t *testing.T) {
	tmp := tmpFile(t)
	f, err := os_.CreateAll(tmp)
	if err != nil {
		t.Fatalf("temp file CreateAll failed: %v", err)
	}
	err = f.Close()
	if err != nil {
		t.Fatalf("temp file Close failed: %v", err)
	}
	if err := os.Remove(tmp); err != nil {
		t.Fatalf("temp file Remove failed: %v", err)
	}
}

func TestTouchAll(t *testing.T) {
	tmp := tmpFile(t)
	f, err := os_.TouchAll(tmp)
	if err != nil {
		t.Fatalf("temp file TouchAll failed: %v", err)
	}
	err = f.Close()
	if err != nil {
		t.Fatalf("temp file Close failed: %v", err)
	}
	if err := os.Remove(tmp); err != nil {
		t.Fatalf("temp file Remove failed: %v", err)
	}
}

func TestCreateAllIfNotExist(t *testing.T) {
	tmp := tmpFile(t)
	f, err := os_.CreateAllIfNotExist(tmp)
	if err != nil {
		t.Fatalf("temp file CreateAllIfNotExist failed: %v", err)
	}
	err = f.Close()
	if err != nil {
		t.Fatalf("temp file Close failed: %v", err)
	}
	if err := os.Remove(tmp); err != nil {
		t.Fatalf("temp file CreateAllIfNotExist failed: %v", err)
	}
}

func TestRelink(t *testing.T) {
	tmpOld := tmpFile(t)
	tmpNew := tmpFile(t)
	func() {
		f, err := os_.CreateAllIfNotExist(tmpOld)
		if err != nil {
			t.Fatalf("temp file CreateAllIfNotExist failed: %v", err)
		}
		defer f.Close()
	}()
	err := os_.ReLink(tmpOld, tmpNew)
	if err != nil {
		t.Fatalf("temp file ReSymlink failed: %v", err)
	}
	if err := os.Remove(tmpOld); err != nil {
		t.Fatalf("temp file[%s] Remove failed: %v", tmpOld, err)
	}
	if err := os.Remove(tmpNew); err != nil {
		t.Fatalf("temp file[%s] Remove failed: %v", tmpNew, err)
	}
}

func TestReSymlink(t *testing.T) {
	tmpOld := tmpFile(t)
	tmpNew := tmpFile(t)
	func() {
		f, err := os_.CreateAllIfNotExist(tmpOld)
		if err != nil {
			t.Fatalf("temp file CreateAllIfNotExist failed: %v", err)
		}
		defer f.Close()
	}()
	err := os_.ReSymlink(tmpOld, tmpNew)
	if err != nil {
		t.Fatalf("temp file ReSymlink failed: %v", err)
	}
	if err := os.Remove(tmpOld); err != nil {
		t.Fatalf("temp file[%s] Remove failed: %v", tmpOld, err)
	}
	if err := os.Remove(tmpNew); err != nil {
		t.Fatalf("temp file[%s] Remove failed: %v", tmpNew, err)
	}
}

func TestOpenFileAll_Dir(t *testing.T) {
	sep := string(filepath.Separator)
	tests := []struct {
		name    string
		path    func(dir string) string
		open    func(path string) (*os.File, error)
		wantErr bool
		wantDir bool
	}{
		{name: "OpenAll not exist", path: func(dir string) string { return filepath.Join(dir, "a", "b") + sep }, open: os_.OpenAll, wantErr: true},
		{name: "OpenAll existing dir", path: func(dir string) string { return dir }, open: os_.OpenAll, wantDir: true},
		{name: "CreateAll trailing separator", path: func(dir string) string { return filepath.Join(dir, "a", "b") + sep }, open: os_.CreateAll, wantErr: true, wantDir: true},
		{name: "CreateAllIfNotExist existing dir", path: func(dir string) string { return dir }, open: os_.CreateAllIfNotExist, wantErr: true, wantDir: true},
		{name: "LockAll trailing separator", path: func(dir string) string { return filepath.Join(dir, "a", "b") + sep }, open: os_.LockAll, wantErr: true, wantDir: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.path(t.TempDir())
			f, err := tt.open(path)
			if f != nil {
				defer f.Close()
			}
			if (err != nil) != tt.wantErr {
				t.Fatalf("open %q error = %v, wantErr %v", path, err, tt.wantErr)
			}
			if err == nil && f == nil {
				t.Fatalf("open %q = nil, nil, want non-nil file", path)
			}
			fi, statErr := os.Stat(path)
			if gotDir := statErr == nil && fi.IsDir(); gotDir != tt.wantDir {
				t.Errorf("%q is a directory = %v, want %v", path, gotDir, tt.wantDir)
			}
		})
	}
}

func TestOpenAll_NotExistNoMkdir(t *testing.T) {
	dir := t.TempDir()
	if _, err := os_.OpenAll(filepath.Join(dir, "a", "f")); !os.IsNotExist(err) {
		t.Errorf("OpenAll() error = %v, want not exist", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "a")); !os.IsNotExist(err) {
		t.Errorf("parent dir is created by OpenAll: %v", err)
	}
}

func TestTouchAll_Dir(t *testing.T) {
	sep := string(filepath.Separator)
	tests := []struct {
		name string
		path func(dir string) string
	}{
		{name: "trailing separator", path: func(dir string) string { return filepath.Join(dir, "a", "b") + sep }},
		{name: "existing dir", path: func(dir string) string { return dir }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.path(t.TempDir())
			old := time.Now().Add(-time.Hour)
			if err := os.MkdirAll(path, 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.Chtimes(path, old, old); err != nil {
				t.Fatal(err)
			}

			f, err := os_.TouchAll(path)
			if err != nil {
				t.Fatalf("TouchAll(%q) error = %v", path, err)
			}
			defer f.Close()
			fi, err := f.Stat()
			if err != nil {
				t.Fatal(err)
			}
			if !fi.IsDir() {
				t.Errorf("TouchAll(%q) opened a non directory", path)
			}
			if !fi.ModTime().After(old) {
				t.Errorf("ModTime = %v, want after %v", fi.ModTime(), old)
			}
		})
	}
}

// setupLinkReplace creates oldname with mode 0600, and newname by create, in a temp dir.
func setupLinkReplace(t *testing.T, create func(dir, newname string)) (dir, oldname, newname string) {
	dir = t.TempDir()
	oldname = filepath.Join(dir, "old")
	newname = filepath.Join(dir, "new")
	if err := os.WriteFile(oldname, []byte("old"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(oldname, 0600); err != nil {
		t.Fatal(err)
	}
	create(dir, newname)
	return dir, oldname, newname
}

// checkLinkReplaced checks that oldname keeps mode 0600, and no temp file is left in dir.
func checkLinkReplaced(t *testing.T, dir, oldname string, wantEntries int) {
	t.Helper()
	fi, err := os.Stat(oldname)
	if err != nil {
		t.Fatal(err)
	}
	if got := fi.Mode().Perm(); got != 0600 {
		t.Errorf("mode of oldname = %v, want %v", got, os.FileMode(0600))
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != wantEntries {
		var names []string
		for _, e := range entries {
			names = append(names, e.Name())
		}
		t.Errorf("entries in dir = %v, want %d entries", names, wantEntries)
	}
}

func TestReLink_ReplaceExisting(t *testing.T) {
	dir, oldname, newname := setupLinkReplace(t, func(dir, newname string) {
		if err := os.WriteFile(newname, []byte("new"), 0644); err != nil {
			t.Fatal(err)
		}
	})
	if err := os_.ReLink(oldname, newname); err != nil {
		t.Fatalf("ReLink() error = %v", err)
	}
	if !os_.SameFile(oldname, newname) {
		t.Errorf("newname is not a hard link to oldname")
	}
	checkLinkReplaced(t, dir, oldname, 2)
}

func TestReSymlink_ReplaceExisting(t *testing.T) {
	tests := []struct {
		name        string
		create      func(dir, newname string)
		wantEntries int
	}{
		{
			name: "regular file",
			create: func(dir, newname string) {
				if err := os.WriteFile(newname, []byte("new"), 0644); err != nil {
					t.Fatal(err)
				}
			},
			wantEntries: 2,
		},
		{
			name: "symlink",
			create: func(dir, newname string) {
				other := filepath.Join(dir, "other")
				if err := os.WriteFile(other, []byte("other"), 0644); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(other, newname); err != nil {
					t.Fatal(err)
				}
			},
			wantEntries: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, oldname, newname := setupLinkReplace(t, tt.create)
			if err := os_.ReSymlink(oldname, newname); err != nil {
				t.Fatalf("ReSymlink() error = %v", err)
			}
			got, err := os.Readlink(newname)
			if err != nil {
				t.Fatal(err)
			}
			if got != oldname {
				t.Errorf("Readlink(newname) = %q, want %q", got, oldname)
			}
			checkLinkReplaced(t, dir, oldname, tt.wantEntries)
		})
	}
}

func TestCopyRenameTruncateAll(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.log")
	dst := filepath.Join(dir, "sub", "dst.log")
	f, err := os.OpenFile(src, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if _, err := f.WriteString("old"); err != nil {
		t.Fatal(err)
	}

	if err := os_.CopyRenameTruncateAll(dst, src); err != nil {
		t.Fatalf("CopyRenameTruncateAll() error = %v", err)
	}
	if _, err := f.WriteString("new"); err != nil {
		t.Fatal(err)
	}

	if got, err := os.ReadFile(src); err != nil || string(got) != "old" {
		t.Errorf("ReadFile(src) = %q, %v, want %q", got, err, "old")
	}
	if got, err := os.ReadFile(dst); err != nil || string(got) != "new" {
		t.Errorf("ReadFile(dst) = %q, %v, want %q", got, err, "new")
	}
}

func TestWriteRenameAll_Perm(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("permission bits are not supported on windows")
	}
	dir := t.TempDir()
	existing := filepath.Join(dir, "existing")
	if err := os.WriteFile(existing, []byte("old"), 0640); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(existing, 0640); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		filename string
		wantPerm os.FileMode
	}{
		{name: "new file", filename: filepath.Join(dir, "sub", "new"), wantPerm: 0600},
		{name: "existing file", filename: existing, wantPerm: 0640},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := os_.WriteRenameAll(tt.filename, []byte("new")); err != nil {
				t.Fatalf("WriteRenameAll() error = %v", err)
			}
			if got, err := os.ReadFile(tt.filename); err != nil || string(got) != "new" {
				t.Errorf("ReadFile() = %q, %v, want %q", got, err, "new")
			}
			fi, err := os.Stat(tt.filename)
			if err != nil {
				t.Fatal(err)
			}
			if got := fi.Mode().Perm(); got != tt.wantPerm {
				t.Errorf("perm = %v, want %v", got, tt.wantPerm)
			}
			entries, err := os.ReadDir(filepath.Dir(tt.filename))
			if err != nil {
				t.Fatal(err)
			}
			for _, e := range entries {
				if strings.HasSuffix(e.Name(), ".rename") {
					t.Errorf("temp file %q left", e.Name())
				}
			}
		})
	}
}

func TestNextFile(t *testing.T) {
	t.Parallel()

	dir, err := os.MkdirTemp("", "TestNextFileBadDir")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	tests := []struct{ pattern, prefix, suffix string }{
		{filepath.Join(dir, "tempfile_test"), "tempfile_test", ""},
		{filepath.Join(dir, "tempfile_test*"), "tempfile_test", ""},
		{filepath.Join(dir, "tempfile_test*xyz"), "tempfile_test", "xyz"},
	}

	for _, test := range tests {
		f, seq, err := os_.NextFile(test.pattern, 0)
		if err != nil {
			t.Errorf("CreateTemp(..., %q) error: %v", test.pattern, err)
			continue
		}
		defer os.Remove(f.Name())
		base := filepath.Base(f.Name())
		f.Close()
		_ = seq
		if !(strings.HasPrefix(base, test.prefix) && strings.HasSuffix(base, test.suffix)) {
			t.Errorf("NextFile pattern %q created bad name %q; want prefix %q suffix %q",
				test.pattern, base, test.prefix, test.suffix)
		}
	}
}
