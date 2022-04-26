[![Build Status](https://travis-ci.org/searKing/travis-ci.svg?branch=mux)](https://travis-ci.org/searKing/travis-ci)
[![GoDoc](https://godoc.org/github.com/searKing/golang/tools/mux?status.svg)](https://godoc.org/github.com/searKing/golang/go/net/mux)
[![Report card](https://goreportcard.com/badge/github.com/searKing/golang/tools/mux)](https://goreportcard.com/report/github.com/searKing/golang/go/net/mux)
[![Sourcegraph](https://sourcegraph.com/github.com/searKing/golang/-/badge.svg)](https://sourcegraph.com/github.com/searKing/travis-ci@mux?badge)

# mux: Connection Mux

mux is a generic Go library to multiplex connections based on their payload. Using mux, you can serve gRPC, SSH, HTTPS,
HTTP, Go RPC, and pretty much any other protocol on the same TCP listener.

## How-To

Simply create your main listener, create a mux for that listener, and then match connections:

```go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"

	"github.com/searKing/golang/go/net/mux"
	"golang.org/x/net/http2/hpack"
)

func main() {
	srv := mux.NewServer()
	defer srv.Close()

	// We first match the connection against HTTP2 fields. If matched, the
	// connection will be sent through the "grpcl" listener.
	grpcl := mux.HandleListener(mux.HTTP2HeaderFieldValue(false, strings.EqualFold, hpack.HeaderField{
		Name:  "Content-Type",
		Value: "application/grpc",
	}))
	//Otherwise, we match it againts a websocket upgrade request.
	var wslH = http.Header{}
	wslH.Set("Upgrade", "websocket")
	wsl := mux.HandleListener(mux.HTTP1HeaderEqual(wslH))

	// Otherwise, we match it againts HTTP1 methods. If matched,
	// it is sent through the "httpl" listener.
	httpl := mux.HandleListener(mux.HTTP1Fast())
	// If not matched by HTTP, we assume it is an RPC connection.
	rpcl := mux.HandleListener(mux.Any())

	// Then we used the muxed listeners.
	// See safeServe in github.com/searKing/golang/go/net/mux/mux_helper_test.go
	go serveGRPC(grpcl)
	go serveWS(wsl)
	go serveHTTP(httpl)
	go serveRPC(rpcl)

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("mux server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()
	if err := srv.ListenAndServe("localhost:0"); err != mux.ErrServerClosed {
		// Error starting or closing listener:
		log.Printf("mux server ListenAndServe: %v", err)
	}
	<-idleConnsClosed
}
```

Take a look at [other examples in the GoDoc](http://godoc.org/github.com/searKing/golang/go/net/mux/#pkg-examples).

## Docs

* [GoDocs](https://godoc.org/github.com/searKing/golang/go/net/mux)

## Performance

There is room for improvment but, since we are only matching the very first bytes of a connection, the performance
overheads on long-lived connections (i.e., RPCs and pipelined HTTP streams)
is negligible.

## Limitations

* *TLS*: `net/http` uses a type assertion to identify TLS connections; since mux's lookahead-implementing connection
  wraps the underlying TLS connection, this type assertion fails. Because of that, you can serve HTTPS using mux
  but `http.Request.TLS`
  would not be set in your handlers.

* *Different Protocols on The Same Connection*: `mux` matches the connection when it's accepted. For example, one
  connection can be either gRPC or REST, but not both. That is, we assume that a client connection is either used for
  gRPC or REST.

* *Java gRPC Clients*: Java gRPC client blocks until it receives a SETTINGS frame from the server. If you are using the
  Java client to connect to a mux'ed gRPC server please match with writers:

```go
package main

import (
	"strings"

	"github.com/searKing/golang/go/net/mux"
	"golang.org/x/net/http2/hpack"
)

func main() {
	grpcl := mux.HandleListener(mux.HTTP2HeaderFieldValue(false, strings.EqualFold, hpack.HeaderField{
		Name:  "Content-Type",
		Value: "application/grpc",
	}))
	_ = grpcl
}
```

## Thanks to

+ [cmux](https://github.com/soheilhy/cmux.git).

# Copyright and License

Copyright 2019 The searKing Authors. All rights reserved.

Code is released under
[the MIT license](https://github.com/searKing/golang/go/net/mux/blob/master/LICENSE).
