package main

import (
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"

	"patcher/api/patcher"
	"patcher/pkg/common/infrastructure/git"
	"patcher/pkg/patchercli/app"
	"patcher/pkg/patchercli/infrastructure/os"
	"patcher/pkg/patchercli/infrastructure/service"
)

func executePush(ctx *cli.Context) error {
	config, err := parseConfig()
	if err != nil {
		return err
	}

	executor, err := git.NewExecutor()
	if err != nil {
		return err
	}

	repoManager := git.NewRepoManager(".", executor)

	client, err := initServiceClient(config)
	if err != nil {
		return err
	}

	projectService := service.NewProjectService(
		repoManager,
		client,
		os.NewHostInfoProvider(),
		initReporter(ctx),
	)

	return projectService.PushCurrentChanges(app.PushCurrentChangesParam{
		NoReset: ctx.Bool("no-reset"),
	})
}

func initServiceClient(config *config) (patcher.PatcherServiceClient, error) {
	conn, err := grpc.Dial(config.PatcherServiceAddress, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return patcher.NewPatcherServiceClient(conn), nil
}
