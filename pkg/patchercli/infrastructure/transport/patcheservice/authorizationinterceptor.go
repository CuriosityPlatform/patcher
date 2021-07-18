package patcheservice

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"patcher/pkg/common/infrastructure/auth"
)

func NewAuthorizationInterceptor(manager auth.TokenManager) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		token, err := manager.NewToken()
		if err != nil {
			return err
		}

		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", token)

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
