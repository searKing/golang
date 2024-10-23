// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package otlpmetricgrpc

import (
	"github.com/searKing/golang/pkg/instrumentation/otel/metric"
	"github.com/searKing/golang/pkg/instrumentation/otel/metric/driver"
)

var _ driver.ReaderURLOpener = (*URLOpener)(nil)

func init() {
	metric.Register(&URLOpener{})
}
