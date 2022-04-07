// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package constraints

import "golang.org/x/exp/constraints"

// Number is a constraint that permits any number type: any type
// that supports the operators < <= >= > - + * /.
// If future releases of Go add new number types,
// this constraint will be modified to include them.
type Number interface {
	constraints.Integer | constraints.Float
}
