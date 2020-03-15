// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcp

import (
	"context"
	"net"
)

type Client struct {
	*Server
}

func NewClientFunc(onOpenHandler OnOpenHandler,
	onMsgReadHandler OnMsgReadHandler,
	onMsgHandleHandler OnMsgHandleHandler,
	onCloseHandler OnCloseHandler,
	onErrorHandler OnErrorHandler) *Client {
	return &Client{
		Server: NewServerFunc(onOpenHandler, onMsgReadHandler, onMsgHandleHandler, onCloseHandler, onErrorHandler),
	}
}
func NewClient(h Handler) *Client {
	return NewClientFunc(h, h, h, h, h)
}

// Deprecated: use DialAndServe instead.
func (cli *Client) ListenAndServe() error {
	return ErrUnImplement
}

func (cli *Client) DialAndServe(network, address string) error {
	if cli.shuttingDown() {
		return ErrClientClosed
	}
	// transfer http to websocket
	conn, err := net.Dial(network, address)
	if cli.Server.CheckError(nil, nil, err) != nil {
		return err
	}

	defer conn.Close()
	ctx := context.WithValue(context.Background(), ClientContextKey, cli)

	// takeover the connect
	c := cli.Server.newConn(conn)
	// Handle websocket On
	err = cli.Server.onOpenHandler.OnOpen(c.rwc)
	if err = cli.Server.CheckError(c.rwc, c.rwc, err); err != nil {
		c.close()
		return err
	}
	c.setState(c.rwc, StateNew) // before Serve can return

	c.serve(ctx)
	return nil
}
