package exec

import (
	"context"
	"io"
)

func CommandWithCancel(handle func(reader io.Reader), name string, arg ...string) (err error) {
	cs, err := newCommandServerWithCancel(handle, name, arg...)
	err = cs.wait()
	if err != nil {
		cs.Stop()
		return err
	}
	return nil
}
func newCommandServerWithCancel(handle func(reader io.Reader), name string, arg ...string) (*commandServer, error) {
	ctx, stop := context.WithCancel(context.Background())
	return newCommandServer(ctx, stop, handle, name, arg...)
}
