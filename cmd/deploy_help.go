package cmd

import (
	"context"

	"github.com/chihqiang/cli"
	"github.com/chihqiang/tlsctl/deploy"
	"github.com/chihqiang/tlsctl/pkg/structs"
)

func deployHelpCommand() *cli.Command {
	return &cli.Command{
		Name:  "help:deploy",
		Usage: `Display fields and environment configs required by each deploy type`,
		Flags: []cli.Flag{},
		Action: func(ctx context.Context, _ *cli.Input, _ *cli.Output) error {
			return printTagTable("deploy", structs.TagsMaps[deploy.IDeploy](deploy.All()))
		},
	}
}
