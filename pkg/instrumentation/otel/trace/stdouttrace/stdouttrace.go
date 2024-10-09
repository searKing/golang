// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stdouttrace

import (
	"github.com/searKing/golang/pkg/instrumentation/otel/trace"
)

var _ trace.ExporterURLOpener = (*URLOpener)(nil)

func init() {
	trace.Register(&URLOpener{})
}
