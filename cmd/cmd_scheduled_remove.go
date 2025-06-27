package cmd

import (
	"wangzhiqiang/tlsctl/deploy"
	"context"
	"github.com/urfave/cli/v3"
)

func removeScheduledCommand() *cli.Command {
	return &cli.Command{
		UseShortOptionHandling: true,
		Name:                   "scheduled:remove",
		Usage:                  "Delete scheduled tasks based on domain name",
		Flags:                  []cli.Flag{},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return deploy.JSONFileRemove(getDeployJson(cmd), []string{getDomain(cmd)})
		},
	}
}
