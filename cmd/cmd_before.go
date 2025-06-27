package cmd

import (
	"wangzhiqiang/tlsctl/pkg/log"
	"context"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v3"
	"path"
)

func Before(ctx context.Context, command *cli.Command) (context.Context, error) {
	if err := godotenv.Load(path.Join(command.String(flgPath), ".env")); err != nil {
		log.Warn("Could not load .env file: %v", err)
	}
	return ctx, nil
}
