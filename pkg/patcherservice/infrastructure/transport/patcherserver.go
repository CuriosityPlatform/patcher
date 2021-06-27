package transport

import (
	"github.com/google/uuid"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/emptypb"

	api "patcher/api/patcher"
	"patcher/pkg/patcherservice/app"
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

func (server *patcherServer) ApplyPatch(ctx context.Context, req *api.ApplyPatchRequest) (*emptypb.Empty, error) {
	patchID, err := uuid.Parse(req.PatchID)
	if err != nil {
		return nil, err
	}

	patchService := server.container.PatchService()

	err = patchService.ApplyPatch(app.PatchID(patchID))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (server *patcherServer) GetPatch(ctx context.Context, req *api.GetPatchRequest) (*api.GetPatchResponse, error) {
	patchID, err := uuid.Parse(req.PatchID)
	if err != nil {
		return nil, err
	}

	patchQueryService := server.container.PatchQueryService()

	patch, err := patchQueryService.GetPatch(app.PatchID(patchID))
	if err != nil {
		return nil, err
	}

	return &api.GetPatchResponse{
		Id:        uuid.UUID(patch.ID).String(),
		Applied:   patch.Applied,
		Author:    string(patch.Author),
		Device:    string(patch.Device),
		CreatedAt: patch.CreatedAt.Unix(),
	}, nil
}

func (server *patcherServer) GetPatchContent(ctx context.Context, req *api.GetPatchContentRequest) (*api.GetPatchContentResponse, error) {
	patchID, err := uuid.Parse(req.PatchID)
	if err != nil {
		return nil, err
	}

	patchQueryService := server.container.PatchQueryService()

	content, err := patchQueryService.GetPatchContent(app.PatchID(patchID))
	if err != nil {
		return nil, err
	}

	return &api.GetPatchContentResponse{
		Content: string(content),
	}, nil
}

func (server *patcherServer) QueryPatches(ctx context.Context, req *api.QueryPatchesRequest) (*api.QueryPatchesResponse, error) {
	showApplied := false

	patchIDs, err := stringsToPatchIDs(req.PatchIDs)
	if err != nil {
		return nil, err
	}

	spec := app.PatchSpecification{
		PatchIDS:    patchIDs,
		Authors:     stringsToPatchAuthors(req.Authors),
		Devices:     stringsToDevices(req.Devices),
		ShowApplied: &showApplied,
	}

	patchQueryService := server.container.PatchQueryService()

	patches, err := patchQueryService.GetPatches(spec)
	if err != nil {
		return nil, err
	}

	result := make([]*api.Patch, 0, len(patches))

	for _, patch := range patches {
		result = append(result, &api.Patch{
			Id:        uuid.UUID(patch.ID).String(),
			Applied:   patch.Applied,
			Author:    string(patch.Author),
			Device:    string(patch.Device),
			CreatedAt: patch.CreatedAt.Unix(),
		})
	}

	return &api.QueryPatchesResponse{Patches: result}, nil
}

func stringsToPatchIDs(ids []string) ([]app.PatchID, error) {
	result := make([]app.PatchID, 0, len(ids))
	for _, id := range ids {
		parsedUUID, err := uuid.Parse(id)
		if err != nil {
			return nil, err
		}

		result = append(result, app.PatchID(parsedUUID))
	}
	return result, nil
}

func stringsToPatchAuthors(authors []string) []app.PatchAuthor {
	result := make([]app.PatchAuthor, 0, len(authors))
	for _, author := range authors {
		result = append(result, app.PatchAuthor(author))
	}
	return result
}

func stringsToDevices(authors []string) []app.Device {
	result := make([]app.Device, 0, len(authors))
	for _, author := range authors {
		result = append(result, app.Device(author))
	}
	return result
}
