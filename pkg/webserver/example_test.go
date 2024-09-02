// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webserver_test

import (
	"context"
	"fmt"
	"net/http"
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
		Name:        "MockWebServer",
		BindAddress: ":8080",
		Validator:   getValidator(t),
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
			t.Errorf("GET %q failed: %s", url, err)
			return
		}
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
