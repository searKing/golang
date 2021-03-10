package os_test

import (
	"io/ioutil"
	"os"
	"testing"

	os_ "github.com/searKing/golang/go/os"
)

// tmpDir creates a temporary directory and returns its name.
func tmpFile(t *testing.T) string {
	tmp, err := ioutil.TempFile("", "")
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
	defer f.Close()
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
	defer f.Close()
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
	defer f.Close()
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
