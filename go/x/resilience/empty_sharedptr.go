package resilience

import (
	"context"
	"github.com/sirupsen/logrus"
)

// BackgroundPtr returns a non-nil, empty Ptr.
func BackgroundSharedPtr(ctx context.Context, l logrus.FieldLogger) *SharedPtr {
	return NewSharedPtr(ctx, func() (Ptr, error) {
		return BackgroundPtr(), nil
	}, l)
}

// TODO returns a non-nil, empty Ptr. Code should use context.TODO when
// it's unclear which Ptr to use or it is not yet available .
func TODOSharedPtr(ctx context.Context, l logrus.FieldLogger) *SharedPtr {
	return NewSharedPtr(ctx, func() (Ptr, error) {
		return TODOPtr(), nil
	}, l)
}
