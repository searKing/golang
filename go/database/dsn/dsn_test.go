package dsn_test

import (
	"fmt"
	"maps"
	"net/url"
	"slices"
	"testing"

	"github.com/searKing/golang/go/database/dsn"
)

func TestSplit(t *testing.T) {
	for _, tc := range []struct {
		dsn     string
		scheme  string
		connect string
		query   url.Values
	}{
		{"mysql://foo:bar@tcp(baz:1234)/db?foo=bar", "mysql", "foo:bar@tcp(baz:1234)/db", url.Values{"foo": {"bar"}}},
		{"foo:bar@tcp(baz:1234)/db", "", "foo:bar@tcp(baz:1234)/db", nil},
		{"sqlite3://file.db?cache=shared", "sqlite3", "file.db", url.Values{"cache": {"shared"}}},
		{"sqlite3:///tmp/file.db?cache=shared", "sqlite3", "/tmp/file.db", url.Values{"cache": {"shared"}}},
		{"postgres://foo:bar@baz/db?sslrootcert=/a@b/ca.pem", "postgres", "foo:bar@baz/db", url.Values{"sslrootcert": {"/a@b/ca.pem"}}},
		{"mysql://foo:bar@tcp(baz:1234)?foo=bar", "mysql", "foo:bar@tcp(baz:1234)", url.Values{"foo": {"bar"}}},
		// '?' in password
		{"mysql://foo:b?r@tcp(baz:1234)/db", "mysql", "foo:b?r@tcp(baz:1234)/db", nil},
		{"mysql://foo:b?r@tcp(baz:1234)/db?foo=bar", "mysql", "foo:b?r@tcp(baz:1234)/db", url.Values{"foo": {"bar"}}},
		{"mysql://foo@email.com:b?r@baz/db?foo=bar", "mysql", "foo@email.com:b?r@baz/db", url.Values{"foo": {"bar"}}},
	} {
		scheme, connect, query := dsn.Split(tc.dsn)
		if scheme != tc.scheme || connect != tc.connect || !maps.EqualFunc(query, tc.query, slices.Equal) {
			t.Errorf("Split(%q) = %q, %q, %v; want %q, %q, %v", tc.dsn, scheme, connect, query, tc.scheme, tc.connect, tc.query)
		}
	}
}

func TestMasking(t *testing.T) {
	for k, tc := range []struct {
		dsn       string
		maskedDsn string
	}{
		{dsn: "mysql://foo:bar@tcp(baz:1234)/db?foo=bar", maskedDsn: "mysql://*:*@tcp(baz:1234)/db?foo=bar"},
		{dsn: "mysql://foo@email.com:bar@tcp(baz:1234)/db?foo=bar", maskedDsn: "mysql://*:*@tcp(baz:1234)/db?foo=bar"},
		{dsn: "postgres://foo:bar@baz:1234/db?foo=bar", maskedDsn: "postgres://*:*@baz:1234/db?foo=bar"},
		{dsn: "postgres://foo@email.com:bar@baz:1234/db?foo=bar", maskedDsn: "postgres://*:*@baz:1234/db?foo=bar"},
		{dsn: "mysql://foo:b?r@tcp(baz:1234)/db", maskedDsn: "mysql://*:*@tcp(baz:1234)/db"},
		{dsn: "mysql://foo:b?r@tcp(baz:1234)/db?foo=bar", maskedDsn: "mysql://*:*@tcp(baz:1234)/db?foo=bar"},
		{dsn: "foo:b?r@tcp(baz:1234)/db?foo=bar", maskedDsn: "*:*@tcp(baz:1234)/db?foo=bar"},
		{dsn: "mysql://foo@email.com:b?r@tcp(baz:1234)/db?foo=bar", maskedDsn: "mysql://*:*@tcp(baz:1234)/db?foo=bar"},
		{dsn: "postgres://foo:b?r@baz:1234/db?sslrootcert=/a@b/ca.pem", maskedDsn: "postgres://*:*@baz:1234/db?sslrootcert=%2Fa%40b%2Fca.pem"},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			maskedDsn := dsn.Masking(tc.dsn)

			if maskedDsn != tc.maskedDsn {
				t.Fatalf("%s, expected %q, got %q", tc.dsn, tc.maskedDsn, maskedDsn)
			}
		})
	}
}
