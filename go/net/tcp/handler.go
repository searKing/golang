// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcp

import (
	"io"
	"net"
)

type OnOpenHandler interface {
	OnOpen(conn net.Conn) error
}
type OnOpenHandlerFunc func(conn net.Conn) error

func (f OnOpenHandlerFunc) OnOpen(conn net.Conn) error { return f(conn) }

type OnMsgReadHandler interface {
	OnMsgRead(b io.Reader) (msg interface{}, err error)
}
type OnMsgReadHandlerFunc func(r io.Reader) (msg interface{}, err error)

func (f OnMsgReadHandlerFunc) OnMsgRead(r io.Reader) (msg interface{}, err error) { return f(r) }

type OnMsgHandleHandler interface {
	OnMsgHandle(b io.Writer, msg interface{}) error
}
type OnMsgHandleHandlerFunc func(w io.Writer, msg interface{}) error

func (f OnMsgHandleHandlerFunc) OnMsgHandle(w io.Writer, msg interface{}) error { return f(w, msg) }

type OnCloseHandler interface {
	OnClose(w io.Writer, r io.Reader) error
}
type OnCloseHandlerFunc func(w io.Writer, r io.Reader) error

func (f OnCloseHandlerFunc) OnClose(w io.Writer, r io.Reader) error { return f(w, r) }

type OnErrorHandler interface {
	OnError(w io.Writer, r io.Reader, err error) error
}
type OnErrorHandlerFunc func(w io.Writer, r io.Reader, err error) error

func (f OnErrorHandlerFunc) OnError(w io.Writer, r io.Reader, err error) error { return f(w, r, err) }

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

func (srv *nopServerClient) OnOpen(conn net.Conn) error { return nil }

func (srv *nopServerClient) OnMsgRead(r io.Reader) (msg interface{}, err error) {
	return nil, nil
}

func (srv *nopServerClient) OnMsgHandle(w io.Writer, msg interface{}) error {
	return nil
}

func (srv *nopServerClient) OnClose(w io.Writer, r io.Reader) error { return nil }

func (srv *nopServerClient) OnError(w io.Writer, r io.Reader, err error) error { return err }
