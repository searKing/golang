// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package supplier

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
