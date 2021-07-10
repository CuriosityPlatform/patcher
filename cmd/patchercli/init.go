package main

import "github.com/urfave/cli/v2"

func executeInit(ctx *cli.Context) error {
	config, err := parseConfig()
	if err != nil {
		return err
	}

	_ = config

	return nil
}
