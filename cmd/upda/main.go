package main

import (
	"context"
	"git.myservermanager.com/varakh/upda/internal/app"
	"git.myservermanager.com/varakh/upda/internal/server"
	"git.myservermanager.com/varakh/upda/internal/terminal"
	"github.com/urfave/cli/v3"
	"log"
	"os"
)

func main() {
	cli.VersionFlag = &cli.BoolFlag{
		Name: "version",
	}

	application := &cli.Command{
		Name:    app.Name,
		Usage:   "command-line interface for upda",
		Version: app.Version,
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
		log.Fatal(err)
	}
}
