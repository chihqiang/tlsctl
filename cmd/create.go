package cmd

import (
	"context"

	"github.com/chihqiang/cli"
)

func createCommand() *cli.Command {
	return &cli.Command{
		Name:  "create",
		Usage: "Obtain and install a new SSL certificate",
		Flags: createFlags(),
		Action: func(ctx context.Context, in *cli.Input, _ *cli.Output) error {
			domains, err := getDomain(in)
			if err != nil {
				return err
			}
			_, err = buildLegoSSL(in, domains)
			return err
		},
	}
}
