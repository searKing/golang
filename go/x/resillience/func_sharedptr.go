package resillience

import (
	"context"
	"github.com/sirupsen/logrus"
)

func NewSharedPtrFunc(ctx context.Context,
	new func() (interface{}, error),
	ready func(x interface{}) error,
	close func(x interface{}), l logrus.FieldLogger) *SharedPtr {
	return NewSharedPtr(ctx, WithFuncNewer(new, ready, close), l)
}
