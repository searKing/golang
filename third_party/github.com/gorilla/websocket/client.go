// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	slices_ "github.com/searKing/golang/go/exp/slices"
)

type ClientHandler interface {
	OnHTTPResponseHandler
	OnOpenHandler
	OnMsgReadHandler
	OnMsgHandleHandler
	OnCloseHandler
	OnErrorHandler
}
type Client struct {
	*Server
	onHttpResponseHandler OnHTTPResponseHandler
}

func NewClientFunc(onHTTPRespHandler OnHTTPResponseHandler,
	onOpenHandler OnOpenHandler,
	onMsgReadHandler OnMsgReadHandler,
	onMsgHandleHandler OnMsgHandleHandler,
	onCloseHandler OnCloseHandler,
	onErrorHandler OnErrorHandler) *Client {
	return &Client{
		Server:                NewServerFunc(nil, onOpenHandler, onMsgReadHandler, onMsgHandleHandler, onCloseHandler, onErrorHandler),
		onHttpResponseHandler: slices_.FirstOrZero[OnHTTPResponseHandler](onHTTPRespHandler, NopOnHTTPResponseHandler),
	}
}
func NewClient(h ClientHandler) *Client {
	return NewClientFunc(h, h, h, h, h, h)
}

// Deprecated: use DialAndServe instead.
func (cli *Client) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return ErrUnImplement
}

// DialAndServe takes over the http handler
func (cli *Client) DialAndServe(urlStr string, requestHeader http.Header) error {
	if cli.shuttingDown() {
		return ErrClientClosed
	}
	// transfer http to websocket
	dialer := *websocket.DefaultDialer
	dialer.HandshakeTimeout = time.Second
	ws, resp, err := dialer.Dial(urlStr, requestHeader)
	if cli.Server.CheckError(nil, err) != nil {
		return err
	}
	// Handle HTTP Response
	err = cli.onHttpResponseHandler.OnHTTPResponse(resp)
	if cli.Server.CheckError(nil, err) != nil {
		return err
	}

	defer ws.Close()
	ctx := context.WithValue(context.Background(), ClientContextKey, cli)

	// takeover this connect
	c := cli.Server.newConn(ws)
	// Handle websocket On
	err = cli.onOpenHandler.OnOpen(c.rwc)
	if err = cli.Server.CheckError(c.rwc, err); err != nil {
		c.close()
		return err
	}
	c.setState(c.rwc, StateNew) // before Serve can return
	c.serve(ctx)
	return nil
}
