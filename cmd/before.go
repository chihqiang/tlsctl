package cmd

import (
	"context"
	"path"

	"github.com/chihqiang/cli"
	"github.com/chihqiang/logx"
	"github.com/joho/godotenv"
)

func Before(ctx context.Context, in *cli.Input, _ *cli.Output) (context.Context, error) {
	// 根命令不持有业务标志；--path 定义在各子命令上，在全局 Before 执行时尚未解析，
	// 因此这里在未显式设置时回退到默认存储路径。
	dir := in.String(flgPath)
	if dir == "" {
		dir = defaultStorePath()
	}
	if err := godotenv.Load(path.Join(dir, ".env")); err != nil {
		logx.Warn("Could not load .env file: %v", err)
	}
	return ctx, nil
}
