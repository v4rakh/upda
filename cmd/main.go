package main

import (
	"context"
	"git.myservermanager.com/varakh/upda/internal/commons"
	"github.com/urfave/cli/v3"
	"log"
	"os"
)

const (
	description = "command-line application for upda"
)

func main() {
	cli.VersionFlag = &cli.BoolFlag{
		Name:  "version",
		Usage: "show version",
	}

	app := &cli.Command{
		Name:    commons.Name,
		Usage:   description,
		Version: commons.Version,
		Commands: []*cli.Command{
			{
				Name:  "server",
				Usage: "Options for server",
				Commands: []*cli.Command{
					{
						Name:   "serve",
						Usage:  "let the server serve",
						Action: serverServe,
					},
				},
			},
			{
				Name:  "webhook",
				Usage: "Options for webhook",
				Commands: []*cli.Command{
					{
						Name:  "create",
						Usage: "Creates a webhook",
						Flags: []cli.Flag{
							serverUrlFlag,
							userFlag,
							passwordFlag,
							rawFlag,
						},
						ArgsUsage: "<label> [<type (generic|diun, default: generic)>] [<ignore-host (true|false, default: false)>]",
						Action:    webhookCreate,
					},
					{
						Name:  "send",
						Usage: "Sends data to a webhook",
						Flags: []cli.Flag{
							serverUrlFlag,
							webhookIdFlag,
							webhookTokenFlag,
						},
						ArgsUsage: "<json payload>",
						Action:    webhookSend,
					},
				},
			},
			{
				Name:  "update",
				Usage: "Options for update",
				Commands: []*cli.Command{
					{
						Name:  "show",
						Usage: "Shows updates",
						Flags: []cli.Flag{
							serverUrlFlag,
							userFlag,
							passwordFlag,
							updatePageSizeFlag,
							rawFlag,
						},
						Action: updateShow,
					},
				},
			},
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
