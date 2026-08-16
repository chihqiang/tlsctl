package cmd

import (
	"context"

	"github.com/urfave/cli/v3"
)

func createCommand() *cli.Command {
	return &cli.Command{
		UseShortOptionHandling: true,
		Name:                   "create",
		Usage:                  "Obtain and install a new SSL certificate",
		Flags:                  []cli.Flag{},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			domains, err := getDomain(cmd)
			if err != nil {
				return err
			}
			_, err = buildLegoSSL(cmd, domains)
			return err
		},
	}
}
