// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webserver_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"testing"
	"time"

	"github.com/searKing/golang/pkg/webserver"
)

func TestNewWebServer(t *testing.T) {
	srv, err := webserver.NewWebServer(webserver.FactoryConfig{
		Name:        "MockWebServer",
		BindAddress: ":8080",
	})
	if err != nil {
		t.Fatalf("create web server failed: %s", err)
	}
	pws, err := srv.PrepareRun()
	if err != nil {
		t.Fatalf("prepare web server failed: %s", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	go func() {
		time.Sleep(time.Millisecond)
		url := "http://localhost:8080/healthz"
		resp, err := http.Get(url)
		if err != nil {
			t.Fatalf("GET %q failed: %s", url, err)
		}
		data, err := httputil.DumpResponse(resp, true)
		if err != nil {
			t.Fatalf("dump response failed: %s", err)
		}
		fmt.Printf("GET %s\n: %s\n", url, string(data))
	}()
	err = pws.Run(ctx)
	if err != nil {
		t.Fatalf("run web server failed: %s", err)
	}
}
