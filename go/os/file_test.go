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
	if err := os_.CreateAll(tmp, 0666); err != nil {
		t.Fatalf("temp file CreateAll failed: %v", err)
	}
	defer os.Remove(tmp)
}

func TestTouchAll(t *testing.T) {
	tmp := tmpFile(t)
	if err := os_.TouchAll(tmp, 0666); err != nil {
		t.Fatalf("temp file TouchAll failed: %v", err)
	}
	defer os.Remove(tmp)
}

func TestCreateAllIfNotExist(t *testing.T) {
	tmp := tmpFile(t)
	if err := os_.CreateAllIfNotExist(tmp, 0666); err != nil {
		t.Fatalf("temp file TouchAll failed: %v", err)
	}
	defer os.Remove(tmp)
}
