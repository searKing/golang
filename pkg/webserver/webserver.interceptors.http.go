// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webserver

import (
	"github.com/rs/cors"

	http_ "github.com/searKing/golang/go/net/http"
	"github.com/searKing/golang/pkg/webserver/pkg/otel"
)

func (f *Factory) HttpHandlerDecorators(decorators ...http_.HandlerDecorator) []http_.HandlerDecorator {
	// otel
	if f.fc.EnableOpenTelemetry {
		decorators = append(decorators, otel.HttpHandlerDecorators()...)
	}

	// cors
	decorators = append(decorators, http_.HandlerDecoratorFunc(cors.New(f.fc.Cors).Handler))
	return decorators
}
