package filepath

import (
	"os"
	"path/filepath"
)

// Pathify Expand, Abs and Clean the path
func Pathify(path string) string {
	p := os.ExpandEnv(path)

	if filepath.IsAbs(p) {
		return filepath.Clean(p)
	}

	p, err := filepath.Abs(p)
	if err == nil {
		return filepath.Clean(p)
	}
	return ""
}
