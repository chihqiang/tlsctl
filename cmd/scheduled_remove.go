package cmd

import (
	"context"

	"github.com/chihqiang/cli"
	"github.com/chihqiang/tlsctl/deploy"
)

func removeScheduledCommand() *cli.Command {
	return &cli.Command{
		Name:  "scheduled:remove",
		Usage: "Delete scheduled tasks based on domain name",
		Flags: []cli.Flag{pathFlag(), domainFlag()},
		Action: func(ctx context.Context, in *cli.Input, _ *cli.Output) error {
			domains, err := getDomain(in)
			if err != nil {
				return err
			}
			return deploy.JSONFileRemove(getDeployJson(in), domains)
		},
	}
}
