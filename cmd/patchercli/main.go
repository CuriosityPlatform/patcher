package main

import (
	stdlog "log"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"patcher/pkg/common/infrastructure/reporter"
)

var (
	appID   = "UNKNOWN"
	version = "UNKNOWN"
)

func main() {
	err := runApp(os.Args)
	if err != nil {
		stdlog.Fatal(err)
	}
}

func runApp(args []string) error {
	app := &cli.App{
		Name:                 appID,
		Version:              version,
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name: "quiet",
			},
		},
		Commands: []cli.Command{
			{
				Name:        "init",
				Description: "Initialize synchronization config for project",
				Action:      executeInit,
				Flags:       nil,
			},
			{
				Name:        "push",
				Description: "Push current patch to server",
				Action:      executePush,
			},
			{
				Name:        "ping",
				Description: "Ping connection to remote service",
				Action:      executePing,
			},
		},
	}

	return app.Run(args)
}

func initReporter(ctx *cli.Context) reporter.Reporter {
	impl := logrus.New()
	impl.SetFormatter(&logrus.TextFormatter{
		TimestampFormat:  time.RFC3339Nano,
		DisableTimestamp: true,
	})

	return reporter.New(
		ctx.Bool("quiet"),
		impl,
	)
}
