package patcheservice

import (
	"google.golang.org/grpc"

	patcherapi "patcher/api/patcher"
)

type ClientConfig struct {
	PatcherServiceAddress string
}

func NewClient(config ClientConfig, interceptor grpc.UnaryClientInterceptor) (patcherapi.PatcherServiceClient, error) {
	opts := []grpc.DialOption{grpc.WithInsecure()}

	if interceptor != nil {
		opts = append(opts, grpc.WithUnaryInterceptor(interceptor))
	}

	conn, err := grpc.Dial(config.PatcherServiceAddress, opts...)
	if err != nil {
		return nil, err
	}

	return patcherapi.NewPatcherServiceClient(conn), nil
}
