package net_test

import (
	"testing"

	"github.com/searKing/golang/go/net"
)

func TestSplitHostPort(t *testing.T) {
	table := []struct {
		Hostport string
		Ip       string
		Port     string
	}{
		{
			Hostport: `localhost`,
			Ip:       "localhost",
			Port:     "",
		},
		{
			Hostport: `localhost:`,
			Ip:       "localhost",
			Port:     "",
		},
		{
			Hostport: `localhost:80`,
			Ip:       "localhost",
			Port:     "80",
		},
		{
			Hostport: `:80`,
			Ip:       "",
			Port:     "80",
		},
		{
			Hostport: `:`,
			Ip:       "",
			Port:     "",
		},
		{
			Hostport: ``,
			Ip:       "",
			Port:     "",
		},
		{
			Hostport: `[::1]:80`,
			Ip:       "::1",
			Port:     "80",
		},
		{
			Hostport: `[::1]`,
			Ip:       "::1",
			Port:     "",
		},
		{
			Hostport: `[::1%lo0]:80`,
			Ip:       "::1%lo0",
			Port:     "80",
		},
		{
			Hostport: `[::1%lo0]:`,
			Ip:       "::1%lo0",
			Port:     "",
		},
	}

	for i, test := range table {
		qIp, qPort, err := net.SplitHostPort(test.Hostport)
		if err != nil {
			t.Errorf("#%d. got err %s, want err nil", i, err)
		}
		if qIp != test.Ip || qPort != test.Port {
			t.Errorf("#%d. got %q:%q, want %q:%q", i, qIp, qPort, test.Ip, test.Port)
		}
	}
}

func TestHostportOrDefault(t *testing.T) {
	table := []struct {
		Hostport, DefHostport, R string
	}{
		{
			Hostport:    ``,
			DefHostport: `127.0.0.1:443`,
			R:           `127.0.0.1:443`,
		},
		{
			Hostport:    `localhost`,
			DefHostport: `127.0.0.1:443`,
			R:           `localhost:443`,
		},
		{
			Hostport:    `localhost:`,
			DefHostport: `127.0.0.1:443`,
			R:           `localhost:443`,
		},
		{
			Hostport:    `localhost:80`,
			DefHostport: `127.0.0.1:443`,
			R:           `localhost:80`,
		},
		{
			Hostport:    `:80`,
			DefHostport: `127.0.0.1:443`,
			R:           `127.0.0.1:80`,
		},
	}

	for i, test := range table {
		qr := net.HostportOrDefault(test.Hostport, test.DefHostport)
		if qr != test.R {
			t.Errorf("%d. got %q, want %q", i, qr, test.R)
		}
	}
}
