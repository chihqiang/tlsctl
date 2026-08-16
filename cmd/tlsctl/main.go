package main

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/chihqiang/tlsctl/cmd"
	"github.com/chihqiang/tlsctl/pkg/log"
	"github.com/urfave/cli/v3"
)

var (
	version = "main"
)

func main() {
	app := &cli.Command{}
	app.Name = "tlsctl"
	app.Usage = "A command line tool designed for developers and operators, it supports the application, renewal, and deployment of SSL/TLS certificates, helping you to easily manage the entire HTTPS process."
	app.Version = version
	cli.VersionPrinter = func(cmd *cli.Command) {
		fmt.Printf("tlsctl version %s %s/%s\n", cmd.Version, runtime.GOOS, runtime.GOARCH)
	}
	app.Flags = cmd.CreateFlags()
	app.Before = cmd.Before
	app.Commands = cmd.CreateCommands()
	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Error(err.Error())
	}
}
