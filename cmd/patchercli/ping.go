package main

import "github.com/urfave/cli"

func executePing(ctx *cli.Context) error {
	config, err := parseConfig()
	if err != nil {
		return err
	}

	_, err = initServiceClient(config)
	return err
}
