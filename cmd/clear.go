package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/chihqiang/cli"
)

func clearCommand() *cli.Command {
	return &cli.Command{
		Name:  "clear",
		Usage: "Clean up all SSL certificates, ACME account data, and related temporary files",
		Flags: pathFlags(),
		Action: func(ctx context.Context, in *cli.Input, _ *cli.Output) error {
			dir := in.String(flgPath)
			if err := validateRemovablePath(dir); err != nil {
				return err
			}
			return os.RemoveAll(dir)
		},
	}
}

// validateRemovablePath 防止 clear 误删危险路径（如根目录、用户主目录、常见系统目录）。
func validateRemovablePath(dir string) error {
	home, _ := os.UserHomeDir()
	abs, err := filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("resolving path %q: %w", dir, err)
	}
	if abs == "/" || abs == home {
		return fmt.Errorf("refusing to remove unsafe path: %q", dir)
	}
	// 常见系统目录前缀，防止误删
	for _, prefix := range []string{"/etc", "/usr", "/bin", "/sbin", "/var", "/opt", "/Library", "/System"} {
		if abs == prefix || strings.HasPrefix(abs, prefix+string(os.PathSeparator)) {
			return fmt.Errorf("refusing to remove system directory: %q", dir)
		}
	}
	return nil
}
