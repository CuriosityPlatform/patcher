package transport

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"patcher/pkg/common/infrastructure/auth"
)

func NewAuthorizationMiddleware(manager auth.TokenManager) GrpcRequestMiddleware {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo) (context.Context, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return ctx, errors.New("failed to extract metadata")
		}

		authorizationData := md.Get("authorization")
		if len(authorizationData) != 1 {
			return ctx, status.Errorf(codes.PermissionDenied, "missing authorization header")
		}

		token := authorizationData[0]

		return ctx, manager.Verify(token)
	}
}
