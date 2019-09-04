package supplier

import "github.com/searKing/golang/go/util/optional"

type OptionalSupplier interface {
	Get() optional.Optional
}

type OptionalSupplierFunc func() optional.Optional

func (supplier OptionalSupplierFunc) Get() optional.Optional {
	return supplier()
}
