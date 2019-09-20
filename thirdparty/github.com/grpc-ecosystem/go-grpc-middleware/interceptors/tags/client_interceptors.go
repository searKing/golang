package tags

import (
	"context"
	"google.golang.org/grpc"
)

// UnaryClientInterceptor returns a new unary client interceptor with tags in context.
func UnaryClientInterceptor(key interface{}, values map[string]interface{}) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		newCtx := newContextTagsForCall(ctx, KindClient, method, key, values)
		return invoker(newCtx, method, req, reply, cc, opts...)
	}
}

// StreamServerInterceptor returns a new streaming client interceptor with tags in context.
func StreamClientInterceptor(key interface{}, values map[string]interface{}) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		newCtx := newContextTagsForCall(ctx, KindClient, method, key, values)
		return streamer(newCtx, desc, cc, method, opts...)
	}
}
