// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webserver_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"reflect"
	"strings"
	"testing"
	"time"
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/searKing/golang/pkg/webserver"
)

func TestNewWebServer(t *testing.T) {
	srv, err := webserver.NewWebServer(webserver.FactoryConfig{
		Name:             "MockWebServer",
		BindAddress:      ":8080",
		Validator:        getValidator(t),
		HTTPTraceLogging: true,
	})
	if err != nil {
		t.Fatalf("create web server failed: %s", err)
	}
	pws, err := srv.PrepareRun()
	if err != nil {
		t.Fatalf("prepare web server failed: %s", err)
	}
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	go func() {
		defer cancel()
		time.Sleep(time.Millisecond)
		url := "http://localhost:8080/healthz"
		resp, err := http.Get(url)
		if err != nil {
			t.Errorf("GET %q failed: %s", url, err)
			return
		}
		defer resp.Body.Close()
		data, err := httputil.DumpResponse(resp, true)
		if err != nil {
			t.Errorf("dump response failed: %s", err)
			return
		}
		fmt.Printf("GET %s\n: %s\n", url, string(data))
	}()
	err = pws.Run(ctx)
	if err != nil {
		t.Fatalf("run web server failed: %s", err)
	}
}

func ExampleWebServer_HttpRoundTripProxyFunc() {
	srv, err := webserver.NewWebServer(webserver.FactoryConfig{
		BindAddress:             "127.0.0.1:0",
		HTTPDynamicHostAndProxy: true,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	// base stands in for the caller's *http.Transport.
	// Production code clones http.DefaultTransport or its own transport the same way.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()
	base := ts.Client().Transport.(*http.Transport).Clone()
	if proxy := srv.HttpRoundTripProxyFunc(); proxy != nil {
		base.Proxy = proxy
	}
	client := &http.Client{
		Transport: srv.HttpRoundTripDecorators().WrapRoundTrip(base),
	}

	resp, err := client.Get(ts.URL)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	fmt.Println(resp.StatusCode)
	// Output:
	// 204
}

func isStrNotContainSpace(fl validator.FieldLevel) bool {
	field := fl.Field()
	switch field.Kind() {
	case reflect.String:
		return strings.IndexFunc(field.String(), unicode.IsSpace) < 0
	default:
		panic(fmt.Sprintf("Bad field type %T", field.Interface()))
	}
}

func getValidator(t *testing.T) *validator.Validate {
	v := validator.New()
	err := v.RegisterValidation("str_not_contain_space", isStrNotContainSpace)
	if err != nil {
		t.Fatalf("register validation failed: %s", err)
	}
	return v
}
