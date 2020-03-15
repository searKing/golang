// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package supplier

import "github.com/searKing/golang/go/util/optional"

type OptionalSupplier interface {
	Get() optional.Optional
}

type OptionalSupplierFunc func() optional.Optional

func (supplier OptionalSupplierFunc) Get() optional.Optional {
	return supplier()
}
