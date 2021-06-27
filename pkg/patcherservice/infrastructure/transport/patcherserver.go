package transport

import (
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/emptypb"

	api "patcher/api/patcher"
	"patcher/pkg/patcherservice/infrastructure"
)

func NewPatcherServer(container infrastructure.DependencyContainer) api.PatcherServiceServer {
	return &patcherServer{container: container}
}

type patcherServer struct {
	container infrastructure.DependencyContainer
}

func (server *patcherServer) AddPatch(ctx context.Context, req *api.AddPatchRequest) (*emptypb.Empty, error) {
	patchService := server.container.PatchService()

	err := patchService.AddPatch(req.Author, req.Device, []byte(req.PatchContent))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
