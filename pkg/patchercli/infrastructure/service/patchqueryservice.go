package service

import (
	"context"
	"time"

	"patcher/api/patcher"
	"patcher/pkg/common/infrastructure/git"
	"patcher/pkg/patchercli/app"
)

func NewPatchQueryService(repoManager git.RepoManager, patcherClient patcher.PatcherServiceClient) app.PatchQueryService {
	return &patchQueryService{repoManager: repoManager, patcherClient: patcherClient}
}

type patchQueryService struct {
	repoManager   git.RepoManager
	patcherClient patcher.PatcherServiceClient
}

func (service *patchQueryService) Query(spec app.PatchSpec) ([]app.Patch, error) {
	resp, err := service.patcherClient.QueryPatches(context.Background(), &patcher.QueryPatchesRequest{
		Projects: spec.Projects,
		Authors:  spec.Authors,
		Devices:  spec.Devices,
	})
	if err != nil {
		return nil, err
	}

	result := make([]app.Patch, 0, len(resp.Patches))

	for _, patch := range resp.Patches {
		unixTime := time.Unix(patch.CreatedAt, 0)
		result = append(result, app.Patch{
			ID:        app.PatchID(patch.Id),
			Project:   patch.Project,
			Applied:   patch.Applied,
			Author:    patch.Author,
			Device:    patch.Device,
			CreatedAt: &unixTime,
		})
	}

	return result, nil
}
