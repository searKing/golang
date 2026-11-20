// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpc_test

import (
	"context"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	reflectionpb "google.golang.org/grpc/reflection/grpc_reflection_v1"

	http_ "github.com/searKing/golang/go/net/http"
	grpc_ "github.com/searKing/golang/third_party/github.com/grpc-ecosystem/grpc-gateway-v2/grpc"
)

func TestWithGrpcOrHttpHandlerDecorators(t *testing.T) {
	var mu sync.Mutex
	var paths []string
	decorator := http_.HandlerDecoratorFunc(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mu.Lock()
			paths = append(paths, r.URL.Path)
			mu.Unlock()
			next.ServeHTTP(w, r)
		})
	})
	seen := func(path string) bool {
		mu.Lock()
		defer mu.Unlock()
		for _, p := range paths {
			if p == path {
				return true
			}
		}
		return false
	}

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	gw := grpc_.NewGateway(l.Addr().String(),
		grpc_.WithGrpcReflectionService(true),
		grpc_.WithGrpcOrHttpHandlerDecorators(decorator))
	go func() { _ = gw.Serve(l) }()
	t.Cleanup(func() { _ = gw.Close() })

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// native gRPC request must pass the decorator and still reach the gRPC server.
	conn, err := grpc.NewClient(l.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	stream, err := reflectionpb.NewServerReflectionClient(conn).ServerReflectionInfo(ctx, grpc.WaitForReady(true))
	if err != nil {
		t.Fatalf("ServerReflectionInfo: %v", err)
	}
	if err := stream.Send(&reflectionpb.ServerReflectionRequest{
		MessageRequest: &reflectionpb.ServerReflectionRequest_ListServices{},
	}); err != nil {
		t.Fatalf("Send: %v", err)
	}
	resp, err := stream.Recv()
	if err != nil {
		t.Fatalf("Recv: %v", err)
	}
	if len(resp.GetListServicesResponse().GetService()) == 0 {
		t.Fatalf("got no services from gRPC reflection, gRPC request may not reach the gRPC server")
	}
	if p := "/grpc.reflection.v1.ServerReflection/ServerReflectionInfo"; !seen(p) {
		t.Errorf("gRPC request %q not seen by decorator", p)
	}

	// HTTP request must pass the decorator too.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+l.Addr().String()+"/not-found", nil)
	if err != nil {
		t.Fatal(err)
	}
	httpResp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("HTTP request: %v", err)
	}
	_ = httpResp.Body.Close()
	if !seen("/not-found") {
		t.Errorf("HTTP request %q not seen by decorator", "/not-found")
	}
}
