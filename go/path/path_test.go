package path_test

import (
	"strings"
	"testing"

	"github.com/searKing/golang/go/path"
)

type RelTests struct {
	root, path, want string
}

var reltests = []RelTests{
	{"a/b", "a/b", "."},
	{"a/b/.", "a/b", "."},
	{"a/b", "a/b/.", "."},
	{"./a/b", "a/b", "."},
	{"a/b", "./a/b", "."},
	{"ab/cd", "ab/cde", "../cde"},
	{"ab/cd", "ab/c", "../c"},
	{"a/b", "a/b/c/d", "c/d"},
	{"a/b", "a/b/../c", "../c"},
	{"a/b/../c", "a/b", "../b"},
	{"a/b/c", "a/c/d", "../../c/d"},
	{"a/b", "c/d", "../../c/d"},
	{"a/b/c/d", "a/b", "../.."},
	{"a/b/c/d", "a/b/", "../.."},
	{"a/b/c/d/", "a/b", "../.."},
	{"a/b/c/d/", "a/b/", "../.."},
	{"../../a/b", "../../a/b/c/d", "c/d"},
	{"/a/b", "/a/b", "."},
	{"/a/b/.", "/a/b", "."},
	{"/a/b", "/a/b/.", "."},
	{"/ab/cd", "/ab/cde", "../cde"},
	{"/ab/cd", "/ab/c", "../c"},
	{"/a/b", "/a/b/c/d", "c/d"},
	{"/a/b", "/a/b/../c", "../c"},
	{"/a/b/../c", "/a/b", "../b"},
	{"/a/b/c", "/a/c/d", "../../c/d"},
	{"/a/b", "/c/d", "../../c/d"},
	{"/a/b/c/d", "/a/b", "../.."},
	{"/a/b/c/d", "/a/b/", "../.."},
	{"/a/b/c/d/", "/a/b", "../.."},
	{"/a/b/c/d/", "/a/b/", "../.."},
	{"/../../a/b", "/../../a/b/c/d", "c/d"},
	{".", "a/b", "a/b"},
	{".", "..", ".."},

	// can't do purely lexically
	{"..", ".", "err"},
	{"..", "a", "err"},
	{"../..", "..", "err"},
	{"a", "/a", "err"},
	{"/a", "a", "err"},
}

func TestRel(t *testing.T) {
	tests := append([]RelTests{}, reltests...)
	for _, test := range tests {
		got, err := path.Rel(test.root, test.path)
		if test.want == "err" {
			if err == nil {
				t.Errorf("Rel(%q, %q)=%q, want error", test.root, test.path, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("Rel(%q, %q): want %q, got error: %s", test.root, test.path, test.want, err)
		}
		if got != test.want {
			t.Errorf("Rel(%q, %q)=%q, want %q", test.root, test.path, got, test.want)
		}
	}
}

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
		toPath := path.ResolveReference(test.FromPath, test.FromBase, test.ToBase)
		if !strings.EqualFold(toPath, test.ToPath) {
			t.Errorf("#%d. got %q, want %q", i, toPath, test.ToPath)
		}
	}
}
