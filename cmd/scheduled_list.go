package cmd

import (
	"context"
	"strings"

	"github.com/chihqiang/cli"
	"github.com/chihqiang/tlsctl/deploy"
	"github.com/chihqiang/tlsctl/pkg/stdout"
)

func listScheduledCommand() *cli.Command {
	return &cli.Command{
		Name:  "scheduled:list",
		Usage: "List scheduled tasks",
		Flags: pathFlags(),
		Action: func(ctx context.Context, in *cli.Input, _ *cli.Output) error {
			domainDeploys, err := deploy.JSONFileLoad(getDeployJson(in))
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
