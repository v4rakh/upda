package main

import (
	"context"
	"git.myservermanager.com/varakh/upda/internal/meta"
	"git.myservermanager.com/varakh/upda/internal/server"
	"git.myservermanager.com/varakh/upda/internal/terminal"
	"github.com/urfave/cli/v3"
	golog "log"
	"os"
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
					server.ServeCmd,
				},
			},
			{
				Name: "webhook",
				Commands: []*cli.Command{
					terminal.WebhookCreateCmd,
					terminal.WebhookSendCmd,
				},
			},
			{
				Name: "update",
				Commands: []*cli.Command{
					terminal.UpdateShowCmd,
				},
			},
		},
	}

	if err := application.Run(context.Background(), os.Args); err != nil {
		golog.Fatal(err)
	}
}
