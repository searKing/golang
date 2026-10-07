// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exec

import (
	"context"
	"io"
)

// CommandWithCancel runs the named program, feeding its stdout to handle, and waits for it to exit.
func CommandWithCancel(handle func(reader io.Reader), name string, arg ...string) (err error) {
	return commandContext(context.Background(), handle, name, arg...)
}
