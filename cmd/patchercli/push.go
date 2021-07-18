package main

import (
	"github.com/urfave/cli/v2"

	"patcher/pkg/common/infrastructure/git"
	"patcher/pkg/patchercli/app"
	"patcher/pkg/patchercli/infrastructure"
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

	container := infrastructure.NewDependencyContainer(infrastructure.DependenciesConfig{
		PatcherServiceAddress: config.PatcherServiceAddress,
		SigningKey:            config.SigningKey,
	})

	client, err := container.PatcherClient()
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
		Message: ctx.String("message"),
	})
}
