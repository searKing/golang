// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package negroni

import "github.com/urfave/negroni"

func Clone(n *negroni.Negroni) *negroni.Negroni {
	// make a deepcopy
	var handlers = make([]negroni.Handler, len(n.Handlers()))
	copy(handlers, n.Handlers())
	return negroni.New(handlers...)
}
