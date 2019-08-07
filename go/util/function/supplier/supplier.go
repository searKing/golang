package supplier

import (
	"github.com/searKing/golang/go/error/exception"
	"github.com/searKing/golang/go/util/class"
)

/**
 * Represents a supplier of results.
 *
 * <p>There is no requirement that a new or distinct result be returned each
 * time the supplier is invoked.
 *
 * <p>This is a <a href="package-summary.html">functional interface</a>
 * whose functional method is {@link #get()}.
 *
 * @param <T> the type of results supplied by this supplier
 *
 * @since 1.8
 */
type Supplier interface {
	/**
	 * Gets a result.
	 *
	 * @return a result
	 */
	Get() interface{}
}
type SupplierFunc func() interface{}

func (supplier SupplierFunc) Get() interface{} {
	return supplier()
}

type AbstractSupplierClass struct {
	class.Class
}

func (sup *AbstractSupplierClass) Get() interface{} {
	panic(exception.NewIllegalStateException1("called wrong Get method"))
}
