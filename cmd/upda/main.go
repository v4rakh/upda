package main

import (
	"context"
	golog "log"
	"os"

	"git.myservermanager.com/varakh/upda/internal/meta"
	"git.myservermanager.com/varakh/upda/internal/server"
	"git.myservermanager.com/varakh/upda/internal/terminal"
	"github.com/urfave/cli/v3"
)

func main() {
	cli.VersionFlag = &cli.BoolFlag{
		Name: "version",
	}

	application := &cli.Command{
		Name:    meta.Name,
		Usage:   "command-line interface for upda",
		Version: meta.Version,
		Commands: []*cli.Command{
			{
				Name: "server",
				Commands: []*cli.Command{
					serverServeCmd,
				},
			},
			{
				Name: "webhook",
				Commands: []*cli.Command{
					terminal.WebhookSendCmd,
				},
			},
		},
	}

	if err := application.Run(context.Background(), os.Args); err != nil {
		golog.Fatal(err)
	}
}

var serverServeCmd = &cli.Command{
	Name:  "serve",
	Usage: "Starts the server and keeps it running",
	Action: func(ctx context.Context, _ *cli.Command) error {
		server := server.New(&ctx)
		server.Start()
		return nil
	},
}
