package main

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"patcher/pkg/common/infrastructure/git"
	"patcher/pkg/patchercli/app"
	"patcher/pkg/patchercli/infrastructure/service"
)

func executeQuery(ctx *cli.Context) error {
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

	reporter := initReporter(ctx)

	queryService := service.NewPatchQueryService(repoManager, client)

	patches, err := queryService.Query(app.PatchSpec{
		Projects: ctx.StringSlice("project"),
		Authors:  ctx.StringSlice("author"),
		Devices:  ctx.StringSlice("device"),
	})
	if err != nil {
		return err
	}

	if len(patches) == 0 {
		reporter.Info("No patches for query found")
		return nil
	}

	reporter.Info("Founded patches (without content):")

	for i, patch := range patches {
		fmt.Println(fmt.Sprintf(
			"ID : %s\n"+
				"Project: %s\n"+
				"Message: %s\n"+
				"IsApplied: %v\n"+
				"Author : %s\n"+
				"Device : %s\n"+
				"Created at : %s",
			patch.ID,
			patch.Project,
			patch.Message,
			patch.Applied,
			patch.Author,
			patch.Device,
			patch.CreatedAt.Format("2006-01-02 15:04:05"),
		))
		if i != len(patches)-1 {
			fmt.Println()
		}
	}

	return nil
}
