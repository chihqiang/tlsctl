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
			domains, err := getDomain(cmd)
			if err != nil {
				return err
			}
			// 部署以第一个域名为准（证书主域名，用于定位已保存的资源）
			res, err := rCache.ReadResource(domains[0])
			if err != nil {
				return err
			}
			return deploy.RunWithJSONFile(getDeployJson(cmd), cmd.String(flgDeploy), res)
		},
	}
}
