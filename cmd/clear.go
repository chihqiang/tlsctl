package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

func clearCommand() *cli.Command {
	return &cli.Command{
		UseShortOptionHandling: true,
		Name:                   "clear",
		Usage:                  "Clean up all SSL certificates, ACME account data, and related temporary files",
		Flags:                  []cli.Flag{},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			dir := cmd.String(flgPath)
			if err := validateRemovablePath(dir); err != nil {
				return err
			}
			return os.RemoveAll(dir)
		},
	}
}

// validateRemovablePath 防止 clear 误删危险路径（如根目录、用户主目录）。
func validateRemovablePath(dir string) error {
	home, _ := os.UserHomeDir()
	if dir == "" || dir == "/" || dir == home {
		return fmt.Errorf("refusing to remove unsafe path: %q", dir)
	}
	return nil
}
