// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmux_test

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/searKing/golang/go/net/cmux"
)

type anotherHTTPHandler struct{}

func (h *anotherHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "example http response")
}

func serveHTTP1(l net.Listener) {
	s := &http.Server{
		Handler: &anotherHTTPHandler{},
	}
	if err := s.Serve(l); err != cmux.ErrListenerClosed {
		panic(err)
	}
}

func serveHTTPS(l net.Listener) {
	// Load certificates.
	certificate, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		log.Panic(err)
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{certificate},
		Rand:         rand.Reader,
	}

	// Create TLS listener.
	tlsl := tls.NewListener(l, config)

	// Serve HTTP over TLS.
	serveHTTP1(tlsl)
}

// This is an example for serving HTTP and HTTPS on the same port.
func ExampleListenAndServe_bothHTTPAndHTTPS() {

	// Create a mux.
	m := cmux.New(context.Background())

	// We first match on HTTP 1.1 methods.
	httpl := m.Match(cmux.HTTP1Fast())

	// If not matched, we assume that its TLS.
	//
	// Note that you can take this listener, do TLS handshake and
	// create another mux to multiplex the connections over TLS.
	tlsl := m.Match(cmux.Any())

	go serveHTTP1(httpl)
	go serveHTTPS(tlsl)

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.
		if err := m.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("cmux server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()
	if err := m.ListenAndServe("localhost:0"); err != cmux.ErrServerClosed {
		// Error starting or closing listener:
		log.Printf("cmux server ListenAndServe: %v", err)
	}
	<-idleConnsClosed

}
