package main

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/chihqiang/cli"
	"github.com/chihqiang/logx"
	"github.com/chihqiang/tlsctl/cmd"
)

var (
	version = "main"
)

func main() {
	app := &cli.Command{}
	app.Name = "tlsctl"
	app.Usage = "A command line tool designed for developers and operators, it supports the application, renewal, and deployment of SSL/TLS certificates, helping you to easily manage the entire HTTPS process."
	app.Version = version
	app.VersionPrinter = func(_ context.Context, _ *cli.Input, _ *cli.Output) {
		fmt.Printf("tlsctl version %s %s/%s\n", app.Version, runtime.GOOS, runtime.GOARCH)
	}
	app.Before = cmd.Before
	app.Subcommands = cmd.CreateCommands()
	if err := app.Run(context.Background(), os.Args); err != nil {
		logx.Error("tlsctl: %v", err)
		os.Exit(1)
	}
}
