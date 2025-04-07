package main

import (
	"context"
	"git.myservermanager.com/varakh/upda/server"
	"github.com/urfave/cli/v3"
)

func serverServe(ctx context.Context, _ *cli.Command) error {
	server.Start(ctx)
	return nil
}
