package cmd

import (
	"context"

	"github.com/chihqiang/tlsctl/deploy"
	"github.com/urfave/cli/v3"
)

func removeScheduledCommand() *cli.Command {
	return &cli.Command{
		UseShortOptionHandling: true,
		Name:                   "scheduled:remove",
		Usage:                  "Delete scheduled tasks based on domain name",
		Flags:                  []cli.Flag{},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			domains, err := getDomain(cmd)
			if err != nil {
				return err
			}
			return deploy.JSONFileRemove(getDeployJson(cmd), domains)
		},
	}
}
