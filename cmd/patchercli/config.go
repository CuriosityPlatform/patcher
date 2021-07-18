package main

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

func parseConfig() (*config, error) {
	c := &config{
		ConfigDir:             "${HOME}/.config",
		PatcherServiceAddress: "patcher.makerov.space:8002",
	}

	if err := envconfig.Process(appID, c); err != nil {
		return nil, errors.Wrap(err, "failed to parse env")
	}
	return c, nil
}

type config struct {
	ConfigDir             string `envconfig:"config_dir"`
	PatcherServiceAddress string `envconfig:"patcher_service_address"`
	SigningKey            string
}
