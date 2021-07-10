package service

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"patcher/api/patcher"
	"patcher/pkg/common/infrastructure/git"
	commonreporter "patcher/pkg/common/infrastructure/reporter"
	"patcher/pkg/patchercli/app"
)

func NewProjectService(
	repoManager git.RepoManager,
	patcherClient patcher.PatcherServiceClient,
	userInfoProvider UserInfoProvider,
	reporter commonreporter.Reporter,
) app.ProjectService {
	return &projectService{
		repoManager:      repoManager,
		patcherClient:    patcherClient,
		userInfoProvider: userInfoProvider,
		reporter:         reporter,
	}
}

type UserInfoProvider interface {
	Username() (string, error)
	DeviceName() (string, error)
}

type projectService struct {
	repoManager      git.RepoManager
	patcherClient    patcher.PatcherServiceClient
	userInfoProvider UserInfoProvider
	reporter         commonreporter.Reporter
}

func (service *projectService) InitializeProject(configsDir string) (string, error) {
	panic("implement me")
}

func (service *projectService) ApplyPatch(param app.ApplyPatchParam) error {
	projectName, err := service.repoManager.RemoteProjectName()
	if err != nil {
		return err
	}

	if param.PatchID == nil {
		id, err2 := service.fetchOnePatchForProject(projectName)
		if err2 != nil {
			return err2
		}

		param.PatchID = &id
	}

	patchContentResp, err := service.patcherClient.GetPatchContent(context.Background(), &patcher.GetPatchContentRequest{
		PatchID: string(*param.PatchID),
	})
	if err != nil {
		return err
	}

	err = service.repoManager.ApplyPatch([]byte(patchContentResp.Content))
	if err != nil {
		return err
	}

	if param.WithApply {
		_, err = service.patcherClient.ApplyPatch(context.Background(), &patcher.ApplyPatchRequest{
			PatchID: string(*param.PatchID),
		})
		if err != nil {
			return err
		}
	}

	service.reporter.Info(fmt.Sprintf("Patch with id %s applied ✅", *param.PatchID))

	return nil
}

func (service *projectService) PushCurrentChanges() error {
	changedFiles, err := service.repoManager.ListChangedFiles()
	if err != nil {
		return err
	}

	service.reporter.Info(fmt.Sprintf("Collecting patch with %d files with changes", len(changedFiles)))

	if len(changedFiles) == 0 {
		service.reporter.Info("No changes or changes ignored")
		return nil
	}

	const maxChangedFilesToNotify = 20

	if len(changedFiles) > maxChangedFilesToNotify {
		service.reporter.Info(fmt.Sprintf("count of changed files more than %d, skipped showing files", maxChangedFilesToNotify))
	} else {
		service.reporter.Info("Changed files:")
		for _, file := range changedFiles {
			service.reporter.Info(fmt.Sprintf("  %s", file))
		}
	}

	changes, err := service.repoManager.GetCurrentChanges(true)
	if err != nil {
		return err
	}

	const maxSizeOfChanges = 10 * 1024 * 1024
	if len(changes) > maxSizeOfChanges {
		return errors.New("size of changes more than allowed")
	}

	username, err := service.userInfoProvider.Username()
	if err != nil {
		return err
	}

	deviceName, err := service.userInfoProvider.DeviceName()
	if err != nil {
		return err
	}

	projectName, err := service.repoManager.RemoteProjectName()
	if err != nil {
		return err
	}

	_, err = service.patcherClient.AddPatch(context.Background(), &patcher.AddPatchRequest{
		Project:      projectName,
		Author:       username,
		Device:       deviceName,
		PatchContent: string(changes),
	})

	if err != nil {
		return err
	}

	service.reporter.Info("Patch sent ✅")
	return nil
}

func (service projectService) fetchOnePatchForProject(projectName string) (app.PatchID, error) {
	resp, err := service.patcherClient.QueryPatches(context.Background(), &patcher.QueryPatchesRequest{
		Projects: []string{projectName},
	})
	if err != nil {
		return "", err
	}

	if len(resp.Patches) == 0 {
		return "", app.ErrNoPatchesForProject
	}

	if len(resp.Patches) != 1 {
		return "", app.ErrTooManyPatchesForProject
	}

	return app.PatchID(resp.Patches[0].Id), nil
}
