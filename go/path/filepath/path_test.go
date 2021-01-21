package filepath_test

import (
	"path/filepath"
	"strings"
	"testing"

	filepath_ "github.com/searKing/golang/go/path/filepath"
)

func TestResolveReference(t *testing.T) {
	table := []struct {
		FromBase, ToBase, FromPath, ToPath string
	}{
		{
			FromBase: "/data/fruits",
			ToBase:   "/data/animals",
			FromPath: "apple",
			ToPath:   "/data/animals/apple",
		},
		{
			FromBase: "/data/fruits",
			ToBase:   "/data/animals",
			FromPath: "/data/fruits/apple",
			ToPath:   "/data/animals/apple",
		},
		{
			FromBase: "/data/fruits",
			ToBase:   "/data/animals",
			FromPath: "./apple",
			ToPath:   "/data/animals/apple",
		},
		{
			FromBase: "/data/fruits",
			ToBase:   "/data/animals",
			FromPath: "/data/stars/moon",
			ToPath:   "/data/stars/moon",
		},
		{
			FromBase: "./data/fruits",
			ToBase:   "/data/animals",
			FromPath: "/data/stars/moon",
			ToPath:   "/data/animals/data/stars/moon",
		},
	}
	for i, test := range table {
		toPath := filepath_.ResolveReference(filepath.FromSlash(test.FromPath), filepath.FromSlash(test.FromBase), filepath.FromSlash(test.ToBase))
		if !strings.EqualFold(filepath.ToSlash(toPath), test.ToPath) {
			t.Errorf("#%d. got %q, want %q", i, toPath, test.ToPath)
		}
	}
}
