package dsn_test

import (
	"fmt"
	"testing"

	"github.com/searKing/golang/go/database/dsn"
)

func TestMasking(t *testing.T) {
	for k, tc := range []struct {
		dsn       string
		maskedDsn string
	}{
		{dsn: "mysql://foo:bar@tcp(baz:1234)/db?foo=bar", maskedDsn: "mysql://*:*@tcp(baz:1234)/db?foo=bar"},
		{dsn: "mysql://foo@email.com:bar@tcp(baz:1234)/db?foo=bar", maskedDsn: "mysql://*:*@tcp(baz:1234)/db?foo=bar"},
		{dsn: "postgres://foo:bar@baz:1234/db?foo=bar", maskedDsn: "postgres://*:*@baz:1234/db?foo=bar"},
		{dsn: "postgres://foo@email.com:bar@baz:1234/db?foo=bar", maskedDsn: "postgres://*:*@baz:1234/db?foo=bar"},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			maskedDsn := dsn.Masking(tc.dsn)

			if maskedDsn != tc.maskedDsn {
				t.Fatalf("%s, expected %q, got %q", tc.dsn, tc.maskedDsn, maskedDsn)
			}
		})
	}
}
