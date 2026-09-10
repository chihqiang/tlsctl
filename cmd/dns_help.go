package cmd

import (
	"context"

	"github.com/chihqiang/cli"
	"github.com/chihqiang/tlsctl/dns"
	"github.com/chihqiang/tlsctl/pkg/structs"
)

func dnsHelpCommand() *cli.Command {
	return &cli.Command{
		Name:  "help:dns",
		Usage: `Show DNS providers and their corresponding environment/configuration fields`,
		Flags: []cli.Flag{},
		Action: func(ctx context.Context, _ *cli.Input, _ *cli.Output) error {
			return printTagTable("dns", structs.TagsMaps[dns.IDNSProvider](dns.All()))
		},
	}
}
