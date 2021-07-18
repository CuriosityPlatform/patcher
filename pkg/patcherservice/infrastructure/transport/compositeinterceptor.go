package transport

import (
	"context"
	"fmt"
	"time"

	log "github.com/CuriosityMusicStreaming/ComponentsPool/pkg/app/logger"
	"google.golang.org/grpc"
)

type GrpcRequestMiddleware func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo) (context.Context, error)

func NewCompositeInterceptor(logger log.Logger, middlewares ...GrpcRequestMiddleware) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		for _, middleware := range middlewares {
			ctx, err = middleware(ctx, req, info)
			if err != nil {
				break
			}
		}

		if err != nil {
			logger.WithFields(log.Fields{
				"args":   req,
				"method": getGRPCMethodName(info),
			}).Error(err, "call failed")
			return nil, err
		}

		start := time.Now()

		resp, err = handler(ctx, req)

		fields := log.Fields{
			"args":     req,
			"duration": fmt.Sprintf("%v", time.Since(start)),
			"method":   getGRPCMethodName(info),
		}

		entry := logger.WithFields(fields)
		if err != nil {
			entry.Error(err, "call failed")
		} else {
			entry.Info("call finished")
		}

		return resp, translateError(err)
	}
}
