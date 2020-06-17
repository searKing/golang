// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import "errors"

// ErrServerClosed is returned by the Server's Serve and ListenAndServe
// methods after a call to Shutdown or Close.
var ErrServerClosed = errors.New("websocket: Server closed")
var ErrNotFound = errors.New("websocket: Server not found")
var ErrClientClosed = errors.New("websocket: Client closed")
var ErrUnImplement = errors.New("UnImplement Method")
