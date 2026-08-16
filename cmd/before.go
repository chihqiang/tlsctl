package cmd

import (
	"context"
	"path"

	"github.com/chihqiang/logx"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v3"
)

func Before(ctx context.Context, command *cli.Command) (context.Context, error) {
	if err := godotenv.Load(path.Join(command.String(flgPath), ".env")); err != nil {
		logx.Warn("Could not load .env file: %v", err)
	}
	return ctx, nil
}
