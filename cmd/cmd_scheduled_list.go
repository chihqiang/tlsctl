package cmd

import (
	"context"
	"strings"

	"github.com/chihqiang/tlsctl/deploy"
	"github.com/chihqiang/tlsctl/pkg/stdout"
	"github.com/urfave/cli/v3"
)

func listScheduledCommand() *cli.Command {
	return &cli.Command{
		UseShortOptionHandling: true,
		Name:                   "scheduled:list",
		Usage:                  "List scheduled tasks",
		Flags:                  []cli.Flag{},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			domainDeploys, err := deploy.JSONFileLoad(getDeployJson(cmd))
			if err != nil {
				return err
			}
			table := stdout.NewTablePrinter()
			table.Add([]string{"Domain", "deploys"})
			for _, domainDeploy := range domainDeploys {
				table.Add([]string{domainDeploy.Domain, strings.Join(domainDeploy.Deploys, ", ")})
			}
			return table.Print()
		},
	}
}
