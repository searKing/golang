// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webserver

import (
	"github.com/rs/cors"
	"github.com/searKing/golang/pkg/webserver/pkg/recovery"

	http_ "github.com/searKing/golang/go/net/http"
	"github.com/searKing/golang/pkg/webserver/pkg/otel"
)

func (f *Factory) HttpHandlerDecorators(decorators ...http_.HandlerDecorator) []http_.HandlerDecorator {
	var s []http_.HandlerDecorator
	// recovery
	s = append(s, recovery.HttpHandlerDecorator())
	// otel
	if f.fc.OtelHandling {
		s = append(s, otel.HttpHandlerDecorators(f.fc.OtelHttpOptions...)...)
	}

	// cors
	s = append(s, http_.HandlerDecoratorFunc(cors.New(f.fc.Cors).Handler))
	return append(s, decorators...)
}
