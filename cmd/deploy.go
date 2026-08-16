package cmd

import (
	"context"

	"github.com/chihqiang/tlsctl/deploy"
	"github.com/urfave/cli/v3"
)

func deployCommand() *cli.Command {
	return &cli.Command{
		UseShortOptionHandling: true,
		Name:                   "deploy",
		Usage:                  `Publish the generated certificate and add it to the scheduled monitoring deployment`,
		Flags:                  []cli.Flag{},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			rCache, err := setupResourceCache(cmd)
			if err != nil {
				return err
			}
			domain, err := getDomain(cmd)
			if err != nil {
				return err
			}
			res, err := rCache.ReadResource(domain)
			if err != nil {
				return err
			}
			return deploy.RunWithJSONFile(getDeployJson(cmd), cmd.String(flgDeploy), res)
		},
	}
}
