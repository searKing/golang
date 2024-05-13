// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import "net/http"

type OnHandshakeHandler interface {
	OnHandshake(http.ResponseWriter, *http.Request) error
}
type OnHandshakeHandlerFunc func(http.ResponseWriter, *http.Request) error

func (f OnHandshakeHandlerFunc) OnHandshake(w http.ResponseWriter, r *http.Request) error {
	return f(w, r)
}

type OnOpenHandler interface {
	OnOpen(conn WebSocketReadWriteCloser) error
}
type OnOpenHandlerFunc func(conn WebSocketReadWriteCloser) error

func (f OnOpenHandlerFunc) OnOpen(conn WebSocketReadWriteCloser) error {
	return f(conn)
}

type OnMsgReadHandler interface {
	OnMsgRead(conn WebSocketReadWriteCloser) (msg any, err error)
}
type OnMsgReadHandlerFunc func(conn WebSocketReadWriteCloser) (msg any, err error)

func (f OnMsgReadHandlerFunc) OnMsgRead(conn WebSocketReadWriteCloser) (msg any, err error) {
	return f(conn)
}

type OnMsgHandleHandler interface {
	OnMsgHandle(conn WebSocketReadWriteCloser, msg any) error
}
type OnMsgHandleHandlerFunc func(conn WebSocketReadWriteCloser, msg any) error

func (f OnMsgHandleHandlerFunc) OnMsgHandle(conn WebSocketReadWriteCloser, msg any) error {
	return f(conn, msg)
}

type OnCloseHandler interface {
	OnClose(conn WebSocketReadWriteCloser) error
}
type OnCloseHandlerFunc func(conn WebSocketReadWriteCloser) error

func (f OnCloseHandlerFunc) OnClose(conn WebSocketReadWriteCloser) error { return f(conn) }

type OnErrorHandler interface {
	OnError(conn WebSocketReadWriteCloser, err error) error
}
type OnErrorHandlerFunc func(conn WebSocketReadWriteCloser, err error) error

func (f OnErrorHandlerFunc) OnError(conn WebSocketReadWriteCloser, err error) error {
	return f(conn, err)
}

type OnHTTPResponseHandler interface {
	OnHTTPResponse(resp *http.Response) error
}
type OnHTTPResponseHandlerFunc func(resp *http.Response) error

func (f OnHTTPResponseHandlerFunc) OnHTTPResponse(resp *http.Response) error {
	return f(resp)
}

var NopOnHandshakeHandler = nopSC
var NopOnOpenHandler = nopSC
var NopOnMsgReadHandler = nopSC
var NopOnMsgHandleHandler = nopSC
var NopOnCloseHandler = nopSC
var NopOnErrorHandler = nopSC
var NopOnHTTPResponseHandler = nopSC
var nopSC = &nopServerClient{}

type NopServer struct{ nopServerClient }
type NopClient struct{ nopServerClient }
type nopServerClient struct {
}

func (srv *nopServerClient) OnHandshake(w http.ResponseWriter, r *http.Request) error { return nil }

func (srv *nopServerClient) OnOpen(conn WebSocketReadWriteCloser) error { return nil }

func (srv *nopServerClient) OnMsgRead(conn WebSocketReadWriteCloser) (msg any, err error) {
	return nil, nil
}

func (srv *nopServerClient) OnMsgHandle(conn WebSocketReadWriteCloser, msg any) error {
	return nil
}

func (srv *nopServerClient) OnClose(conn WebSocketReadWriteCloser) error { return nil }

func (srv *nopServerClient) OnError(conn WebSocketReadWriteCloser, err error) error { return err }

func (srv *nopServerClient) OnHTTPResponse(resp *http.Response) error { return nil }
