package cmd

import (
	"context"
	"github.com/urfave/cli/v3"
	"os"
)

func clearCommand() *cli.Command {
	return &cli.Command{
		UseShortOptionHandling: true,
		Name:                   "clear",
		Usage:                  "Clean up all SSL certificates, ACME account data, and related temporary files",
		Flags:                  []cli.Flag{},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return os.RemoveAll(cmd.String(flgPath))
		},
	}
}
