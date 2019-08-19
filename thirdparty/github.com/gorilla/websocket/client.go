package websocket

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/searKing/golang/go/util/object"
	"net/http"
	"time"
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
	httpRespHandler OnHTTPResponseHandler
}

func NewClientFunc(onHTTPRespHandler OnHTTPResponseHandler,
	onOpenHandler OnOpenHandler,
	onMsgReadHandler OnMsgReadHandler,
	onMsgHandleHandler OnMsgHandleHandler,
	onCloseHandler OnCloseHandler,
	onErrorHandler OnErrorHandler) *Client {
	return &Client{
		Server:          NewServerFunc(nil, onOpenHandler, onMsgReadHandler, onMsgHandleHandler, onCloseHandler, onErrorHandler),
		httpRespHandler: object.RequireNonNullElse(onHTTPRespHandler, NopOnHTTPResponseHandler).(OnHTTPResponseHandler),
	}
}
func NewClient(h ClientHandler) *Client {
	return NewClientFunc(h, h, h, h, h, h)
}

// Deprecated: use DialAndServe instead.
func (cli *Client) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return ErrUnImplement
}

// OnHandshake takes over the http handler
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
	err = OnHTTPResponse(resp)
	if cli.Server.CheckError(nil, err) != nil {
		return err
	}

	defer ws.Close()
	ctx := context.WithValue(context.Background(), ClientContextKey, cli)

	// takeover the connect
	c := cli.Server.newConn(ws)
	// Handle websocket On
	err = OnOpen(c.rwc)
	if err = cli.Server.CheckError(c.rwc, err); err != nil {
		c.close()
		return err
	}
	c.setState(c.rwc, StateNew) // before Serve can return
	c.serve(ctx)
	return nil
}
