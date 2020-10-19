package sql

import "github.com/searKing/golang/go/strings"

// Placeholders behaves like strings.Join([]string{"?",...,"?"}, ",")
func Placeholders(n int) string {
	return strings.JoinRepeat("?", ",", n)
}
