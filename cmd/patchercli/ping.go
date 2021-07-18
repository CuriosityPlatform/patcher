package main

import (
	"github.com/urfave/cli/v2"

	"patcher/pkg/patchercli/infrastructure"
)

func executePing(ctx *cli.Context) error {
	config, err := parseConfig()
	if err != nil {
		return err
	}

	container := infrastructure.NewDependencyContainer(infrastructure.DependenciesConfig{
		PatcherServiceAddress: config.PatcherServiceAddress,
		SigningKey:            config.SigningKey,
	})

	_, err = container.PatcherClient()
	return err
}
