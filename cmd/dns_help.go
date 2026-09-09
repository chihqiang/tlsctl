package cmd

import (
	"context"
	"fmt"

	"github.com/chihqiang/cli"
	"github.com/chihqiang/tlsctl/dns"
	"github.com/chihqiang/tlsctl/pkg/stdout"
	"github.com/chihqiang/tlsctl/pkg/structs"
)

func dnsHelpCommand() *cli.Command {
	return &cli.Command{
		Name:  "help:dns",
		Usage: `Show DNS providers and their corresponding environment/configuration fields`,
		Flags: []cli.Flag{},
		Action: func(ctx context.Context, _ *cli.Input, _ *cli.Output) error {
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
