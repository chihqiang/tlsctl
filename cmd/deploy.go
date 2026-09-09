package cmd

import (
	"context"

	"github.com/chihqiang/cli"
	"github.com/chihqiang/tlsctl/deploy"
)

func deployCommand() *cli.Command {
	return &cli.Command{
		Name:  "deploy",
		Usage: `Publish the generated certificate and add it to the scheduled monitoring deployment`,
		Flags: deployFlags(),
		Action: func(ctx context.Context, in *cli.Input, _ *cli.Output) error {
			rCache, err := setupResourceCache(in)
			if err != nil {
				return err
			}
			domains, err := getDomain(in)
			if err != nil {
				return err
			}
			// 部署以第一个域名为准（证书主域名，用于定位已保存的资源）
			res, err := rCache.ReadResource(domains[0])
			if err != nil {
				return err
			}
			return deploy.RunWithJSONFile(getDeployJson(in), in.String(flgDeploy), res)
		},
	}
}
