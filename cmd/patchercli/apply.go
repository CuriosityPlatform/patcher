package main

import (
	"github.com/urfave/cli/v2"

	"patcher/pkg/common/infrastructure/git"
	"patcher/pkg/patchercli/app"
	"patcher/pkg/patchercli/infrastructure/os"
	"patcher/pkg/patchercli/infrastructure/service"
)

func executeApply(ctx *cli.Context) error {
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

	var patchID *app.PatchID
	if patchIDString := ctx.String("patch"); patchIDString != "" {
		patchID = (*app.PatchID)(&patchIDString)
	}

	return projectService.ApplyPatch(app.ApplyPatchParam{
		PatchID:   patchID,
		WithApply: ctx.Bool("no-apply"),
	})
}
