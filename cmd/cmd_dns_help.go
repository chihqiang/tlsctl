package cmd

import (
	"context"
	"fmt"

	"github.com/chihqiang/tlsctl/dns"
	"github.com/chihqiang/tlsctl/pkg/stdout"
	"github.com/chihqiang/tlsctl/pkg/structs"
	"github.com/urfave/cli/v3"
)

func dnsHelpCommand() *cli.Command {
	return &cli.Command{
		UseShortOptionHandling: true,
		Name:                   "help:dns",
		Usage:                  `Show DNS providers and their corresponding environment/configuration fields`,
		Flags:                  []cli.Flag{},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			maps := structs.TagsMaps[dns.IDNSProvider](dns.All())
			tp := stdout.NewTablePrinter()
			for _, key := range maps.Keys {
				tp.SetTitle(fmt.Sprintf("`%s` dns Environment variables and other tags", key))
				tp.Add([]string{"Field Name", "Type", "Environment", "JSON", "Yaml"})
				for _, tag := range maps.Maps[key] {
					tp.Add([]string{tag.Field, tag.Type, tag.Env, tag.Json, tag.Yaml})
				}
				tp.Print()
				tp.Reset()
			}
			return nil
		},
	}
}
